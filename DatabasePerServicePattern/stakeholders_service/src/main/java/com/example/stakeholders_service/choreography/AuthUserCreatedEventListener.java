package com.example.stakeholders_service.choreography;

import com.example.stakeholders_service.Config.ChoreographyRabbitMQConfig;
import com.example.stakeholders_service.Model.Role;
import com.example.stakeholders_service.Model.User;
import com.example.stakeholders_service.Repository.UserRepository;
import com.example.stakeholders_service.choreography.event.AuthUserCreatedEvent;
import com.example.stakeholders_service.choreography.event.AuthUserRollbackEvent;
import com.example.stakeholders_service.choreography.event.StakeholderProfileCreatedEvent;
import lombok.extern.slf4j.Slf4j;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.stereotype.Component;

/**
 * Stakeholder servisov učesnik u KOREOGRAFISANOM SAGA pattern-u.
 *
 * <p>Razlika od {@code RegistrationCommandListener} (orkestrisani pattern):
 *
 * <pre>
 * ORKESTRISANI (RegistrationCommandListener):
 *   - Prima KOMANDU: "Kreiraj ovog korisnika" (RegisterUserCommand)
 *   - Šalje REPLY: "Evo rezultata" (UserRegisteredReply)
 *   - Orchestrator (auth servis) koordiniše ceo tok
 *   - Stakeholder ne zna za celokupan tok sage
 *
 * KOREOGRAFISANI (AuthUserCreatedEventListener):
 *   - Prima DOGAĐAJ: "Auth servis je upisao korisnika" (AuthUserCreatedEvent)
 *   - Objavljuje novi DOGAĐAJ: "Profil je kreiran" ili "Treba rolbek"
 *   - Nema orchestratora - svaki servis zna svoju ulogu
 *   - Stakeholder zna da treba da kreira profil kad čuje AuthUserCreatedEvent
 * </pre>
 *
 * <p>Kompenzaciona logika:
 * <ul>
 *   <li>Uspeh → objavljuje {@code StakeholderProfileCreatedEvent} → auth označava COMPLETED</li>
 *   <li>Pad   → objavljuje {@code AuthUserRollbackEvent}          → auth briše iz auth baze</li>
 * </ul>
 */
@Slf4j
@Component
public class AuthUserCreatedEventListener {

    private final UserRepository userRepository;
    private final RabbitTemplate rabbitTemplate;

    public AuthUserCreatedEventListener(UserRepository userRepository,
                                        RabbitTemplate rabbitTemplate) {
        this.userRepository = userRepository;
        this.rabbitTemplate = rabbitTemplate;
    }

    /**
     * Obrada AuthUserCreatedEvent:
     * 1. Provjeri da li korisnik već postoji (idempotentnost)
     * 2. Kreiraj User entitet i sačuvaj u stakeholder bazu
     * 3. Objavi StakeholderProfileCreatedEvent (uspeh) ili AuthUserRollbackEvent (pad)
     */
    @RabbitListener(queues = ChoreographyRabbitMQConfig.AUTH_USER_CREATED_QUEUE)
    public void handleAuthUserCreatedEvent(AuthUserCreatedEvent event) {
        log.info("[CHOREOGRAPHY][STAKEHOLDER] Received AuthUserCreatedEvent sagaId={} user={}",
                event.getSagaId(), event.getUsername());

        try {
            // Idempotentnost: ako korisnik već postoji (duplikata poruka), ne prijavljujemo grešku
            if (userRepository.findByUserId(event.getAuthId()).isPresent()) {
                log.warn("[CHOREOGRAPHY][STAKEHOLDER] sagaId={} authId={} already exists, skipping",
                        event.getSagaId(), event.getAuthId());
                publishSuccess(event.getSagaId());
                return;
            }

            // Mapiraj ulogu na Role enum
            Role role;
            try {
                role = Role.valueOf(event.getRole().toUpperCase());
            } catch (IllegalArgumentException e) {
                String reason = "Unknown role: " + event.getRole();
                log.error("[CHOREOGRAPHY][STAKEHOLDER] sagaId={} {}", event.getSagaId(), reason);
                publishRollback(event.getSagaId(), reason);
                return;
            }

            // Kreiraj i sačuvaj korisnika u stakeholder bazi
            User user = new User();
            user.setUserId(event.getAuthId());      // registry_id = auth DB id
            user.setUsername(event.getUsername());
            user.setEmail(event.getEmail());
            user.setRole(role);
            user.setFirstName(event.getFirstname());
            user.setLastName(event.getLastname());
            user.setBlocked(false);

            userRepository.save(user);

            log.info("[CHOREOGRAPHY][STAKEHOLDER] sagaId={} user {} saved in stakeholder DB",
                    event.getSagaId(), event.getUsername());

            // Objavi događaj da je profil kreiran → auth servis označava sagu COMPLETED
            publishSuccess(event.getSagaId());

        } catch (Exception e) {
            log.error("[CHOREOGRAPHY][STAKEHOLDER] sagaId={} ERROR saving user: {}",
                    event.getSagaId(), e.getMessage(), e);
            // Objavi rolbek događaj → auth servis će obrisati korisnika iz auth baze
            publishRollback(event.getSagaId(), e.getMessage());
        }
    }

    // ── Pomoćne metode za objavljivanje događaja ──────────────────────────────

    /**
     * Objavljuje StakeholderProfileCreatedEvent.
     * Auth servis čuje ovaj događaj i označava sagu kao COMPLETED.
     */
    private void publishSuccess(String sagaId) {
        StakeholderProfileCreatedEvent event =
                new StakeholderProfileCreatedEvent(sagaId, true);
        try {
            rabbitTemplate.convertAndSend(
                    ChoreographyRabbitMQConfig.CHOREOGRAPHY_EXCHANGE,
                    ChoreographyRabbitMQConfig.STAKEHOLDER_PROFILE_CREATED_KEY,
                    event);
            log.info("[CHOREOGRAPHY][STAKEHOLDER] sagaId={} StakeholderProfileCreatedEvent published",
                    sagaId);
        } catch (Exception e) {
            log.error("[CHOREOGRAPHY][STAKEHOLDER] sagaId={} ERROR publishing success event: {}",
                    sagaId, e.getMessage(), e);
        }
    }

    /**
     * Objavljuje AuthUserRollbackEvent.
     * Auth servis čuje ovaj događaj i pokreće kompenzaciju (brisanje iz auth baze).
     */
    private void publishRollback(String sagaId, String errorMessage) {
        AuthUserRollbackEvent event = new AuthUserRollbackEvent(sagaId, errorMessage);
        try {
            rabbitTemplate.convertAndSend(
                    ChoreographyRabbitMQConfig.CHOREOGRAPHY_EXCHANGE,
                    ChoreographyRabbitMQConfig.AUTH_USER_ROLLBACK_KEY,
                    event);
            log.info("[CHOREOGRAPHY][STAKEHOLDER] sagaId={} AuthUserRollbackEvent published (compensation triggered)",
                    sagaId);
        } catch (Exception e) {
            log.error("[CHOREOGRAPHY][STAKEHOLDER] sagaId={} ERROR publishing rollback event: {}",
                    sagaId, e.getMessage(), e);
        }
    }
}
