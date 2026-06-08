package com.example.stakeholders_service.Saga;

import com.example.stakeholders_service.Model.Role;
import com.example.stakeholders_service.Model.User;
import com.example.stakeholders_service.Repository.UserRepository;
import com.example.stakeholders_service.Saga.command.RegisterUserCommand;
import com.example.stakeholders_service.Saga.reply.UserRegisteredReply;
import com.example.stakeholders_service.Config.RabbitMQConfig;
import lombok.extern.slf4j.Slf4j;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.stereotype.Component;

/**
 * CommandListener za orchestrated SAGA na strani stakeholder servisa.
 *
 * Sluša na "register.user.command.queue" za RegisterUserCommand poruke koje
 * šalje SagaOrchestrator (auth servis). Kreira korisnika u stakeholder bazi,
 * zatim šalje UserRegisteredReply nazad orchestratoru.
 *
 * Kompenzacija NIJE potrebna na ovoj strani jer:
 * - Ako korak uspe → saga je završena (COMPLETED)
 * - Ako korak padne → Orchestrator (auth servis) sam briše korisnika iz auth baze
 */
@Slf4j
@Component
public class RegistrationCommandListener {

    private final UserRepository  userRepository;
    private final RabbitTemplate  rabbitTemplate;

    public RegistrationCommandListener(UserRepository userRepository,
                                       RabbitTemplate rabbitTemplate) {
        this.userRepository = userRepository;
        this.rabbitTemplate = rabbitTemplate;
    }

    /**
     * Obrada RegisterUserCommand:
     * 1. Proveri da li korisnik već postoji (idempotentnost)
     * 2. Kreiraj User entitet i sačuvaj u stakeholder bazu
     * 3. Pošalji UserRegisteredReply orchestratoru
     */
    @RabbitListener(queues = RabbitMQConfig.REGISTER_USER_CMD_QUEUE)
    public void handleRegisterUserCommand(RegisterUserCommand cmd) {
        log.info("[ORCHESTRATION][STAKEHOLDER] Received RegisterUserCommand sagaId={} user={}",
                cmd.getSagaId(), cmd.getUsername());

        UserRegisteredReply reply;

        try {
            // Idempotentnost: ako korisnik već postoji (duplikata poruka), ne prijavljujemo grešku
            if (userRepository.findByUserId(cmd.getAuthId()).isPresent()) {
                log.warn("[ORCHESTRATION][STAKEHOLDER] sagaId={} user authId={} already exists, skipping insert",
                        cmd.getSagaId(), cmd.getAuthId());
                reply = new UserRegisteredReply(cmd.getSagaId(), true, null);
            } else {
                // Mapiraj ulogu na Role enum
                Role role;
                try {
                    role = Role.valueOf(cmd.getRole().toUpperCase());
                } catch (IllegalArgumentException e) {
                    String reason = "Unknown role: " + cmd.getRole();
                    log.error("[ORCHESTRATION][STAKEHOLDER] sagaId={} {}", cmd.getSagaId(), reason);
                    reply = new UserRegisteredReply(cmd.getSagaId(), false, reason);
                    sendReply(reply, cmd.getSagaId());
                    return;
                }

                // Kreiraj i sačuvaj korisnika
                User user = new User();
                user.setUserId(cmd.getAuthId());       // registry_id = auth DB id
                user.setUsername(cmd.getUsername());
                user.setEmail(cmd.getEmail());
                user.setRole(role);
                user.setFirstName(cmd.getFirstname());
                user.setLastName(cmd.getLastname());
                user.setBlocked(false);

                userRepository.save(user);

                log.info("[ORCHESTRATION][STAKEHOLDER] sagaId={} user {} saved in stakeholder DB",
                        cmd.getSagaId(), cmd.getUsername());
                reply = new UserRegisteredReply(cmd.getSagaId(), true, null);
            }

        } catch (Exception e) {
            log.error("[ORCHESTRATION][STAKEHOLDER] sagaId={} ERROR saving user: {}",
                    cmd.getSagaId(), e.getMessage(), e);
            reply = new UserRegisteredReply(cmd.getSagaId(), false, e.getMessage());
        }

        sendReply(reply, cmd.getSagaId());
    }

    // ── Pomoćna metoda za slanje reply-a ─────────────────────────────────────

    private void sendReply(UserRegisteredReply reply, String sagaId) {
        try {
            rabbitTemplate.convertAndSend(
                    RabbitMQConfig.ORCHESTRATION_EXCHANGE,
                    RabbitMQConfig.USER_REGISTERED_REPLY_KEY,
                    reply);
            log.info("[ORCHESTRATION][STAKEHOLDER] sagaId={} UserRegisteredReply sent (success={})",
                    sagaId, reply.isSuccess());
        } catch (Exception e) {
            log.error("[ORCHESTRATION][STAKEHOLDER] sagaId={} ERROR sending UserRegisteredReply: {}",
                    sagaId, e.getMessage(), e);
        }
    }
}
