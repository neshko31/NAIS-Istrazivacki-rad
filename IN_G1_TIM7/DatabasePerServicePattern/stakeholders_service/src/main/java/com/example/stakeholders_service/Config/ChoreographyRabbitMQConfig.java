package com.example.stakeholders_service.Config;

import org.springframework.amqp.core.*;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * RabbitMQ konfiguracija za KOREOGRAFISANI SAGA pattern (stakeholder servis).
 *
 * <p>Razlika od orkestrisanog (RabbitMQConfig.java):
 * <ul>
 *   <li>ORKESTRISANI: stakeholder prima KOMANDE i šalje REPLY-eve orchestratoru</li>
 *   <li>KOREOGRAFISANI: stakeholder prima DOGAĐAJE i objavljuje nove DOGAĐAJE</li>
 * </ul>
 *
 * <p>Exchange: saga.choreography.exchange (DirectExchange)
 *
 * <p>Queues koje PRIMA stakeholder servis:
 * <ul>
 *   <li>auth.user.created.queue ← prima AuthUserCreatedEvent od auth servisa</li>
 * </ul>
 *
 * <p>Queues na koje ŠALJE stakeholder servis:
 * <ul>
 *   <li>stakeholder.profile.created.queue → šalje StakeholderProfileCreatedEvent (uspeh)</li>
 *   <li>auth.user.rollback.queue          → šalje AuthUserRollbackEvent (pad)</li>
 * </ul>
 *
 * <p>Koristi ISTI RabbitMQ broker kao orkestrisani pattern, ali RAZLIČIT exchange
 * i queue-ove, pa oba pristupa mogu raditi istovremeno bez konflikata.
 */
@Configuration
public class ChoreographyRabbitMQConfig {

    // ── Exchange ──────────────────────────────────────────────────────────────
    public static final String CHOREOGRAPHY_EXCHANGE = "saga.choreography.exchange";

    // ── Queue koji stakeholder PRIMA (od auth servisa) ────────────────────────
    public static final String AUTH_USER_CREATED_QUEUE = "auth.user.created.queue";
    public static final String AUTH_USER_CREATED_KEY   = "auth.user.created";

    // ── Queue na koji stakeholder ŠALJE (uspeh) ───────────────────────────────
    public static final String STAKEHOLDER_PROFILE_CREATED_QUEUE = "stakeholder.profile.created.queue";
    public static final String STAKEHOLDER_PROFILE_CREATED_KEY   = "stakeholder.profile.created";

    // ── Queue na koji stakeholder ŠALJE (pad → rolbek) ───────────────────────
    public static final String AUTH_USER_ROLLBACK_QUEUE = "auth.user.rollback.queue";
    public static final String AUTH_USER_ROLLBACK_KEY   = "auth.user.rollback";

    // ── Exchange bean ─────────────────────────────────────────────────────────

    @Bean
    public DirectExchange choreographyExchange() {
        return new DirectExchange(CHOREOGRAPHY_EXCHANGE);
    }

    // ── Queue beans ───────────────────────────────────────────────────────────

    @Bean
    public Queue authUserCreatedQueue() {
        return QueueBuilder.durable(AUTH_USER_CREATED_QUEUE).build();
    }

    @Bean
    public Queue stakeholderProfileCreatedQueue() {
        return QueueBuilder.durable(STAKEHOLDER_PROFILE_CREATED_QUEUE).build();
    }

    @Bean
    public Queue authUserRollbackQueue() {
        return QueueBuilder.durable(AUTH_USER_ROLLBACK_QUEUE).build();
    }

    // ── Binding beans ─────────────────────────────────────────────────────────

    @Bean
    public Binding authUserCreatedBinding() {
        return BindingBuilder
                .bind(authUserCreatedQueue())
                .to(choreographyExchange())
                .with(AUTH_USER_CREATED_KEY);
    }

    @Bean
    public Binding stakeholderProfileCreatedBinding() {
        return BindingBuilder
                .bind(stakeholderProfileCreatedQueue())
                .to(choreographyExchange())
                .with(STAKEHOLDER_PROFILE_CREATED_KEY);
    }

    @Bean
    public Binding authUserRollbackBinding() {
        return BindingBuilder
                .bind(authUserRollbackQueue())
                .to(choreographyExchange())
                .with(AUTH_USER_ROLLBACK_KEY);
    }
}
