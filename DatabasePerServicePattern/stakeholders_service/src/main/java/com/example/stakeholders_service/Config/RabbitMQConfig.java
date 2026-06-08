package com.example.stakeholders_service.Config;

import org.springframework.amqp.core.*;
import org.springframework.amqp.rabbit.config.SimpleRabbitListenerContainerFactory;
import org.springframework.amqp.rabbit.connection.ConnectionFactory;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.amqp.support.converter.Jackson2JsonMessageConverter;
import org.springframework.amqp.support.converter.MessageConverter;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * RabbitMQ konfiguracija za stakeholder servis.
 *
 * Exchange: saga.registration.exchange (DirectExchange)
 *
 * Queues koje PRIMA stakeholder servis:
 *   register.user.command.queue  ← prima RegisterUserCommand od auth servisa
 *
 * Queues na koje ŠALJE stakeholder servis (reply):
 *   user.registered.reply.queue  → šalje UserRegisteredReply auth servisu
 */
@Configuration
public class RabbitMQConfig {

    // ── Exchange ──────────────────────────────────────────────────────────────
    public static final String ORCHESTRATION_EXCHANGE = "saga.registration.exchange";

    // ── Komandni queue (prima stakeholder) ───────────────────────────────────
    public static final String REGISTER_USER_CMD_QUEUE = "register.user.command.queue";
    public static final String REGISTER_USER_CMD_KEY   = "register.user.command";

    // ── Reply queue (šalje stakeholder → prima auth) ─────────────────────────
    public static final String USER_REGISTERED_REPLY_QUEUE = "user.registered.reply.queue";
    public static final String USER_REGISTERED_REPLY_KEY   = "user.registered.reply";

    // ── Message converter i template ─────────────────────────────────────────

    @Bean
    public MessageConverter jsonMessageConverter() {
        return new Jackson2JsonMessageConverter();
    }

    @Bean
    public RabbitTemplate rabbitTemplate(ConnectionFactory connectionFactory) {
        RabbitTemplate template = new RabbitTemplate(connectionFactory);
        template.setMessageConverter(jsonMessageConverter());
        return template;
    }

    @Bean
    public SimpleRabbitListenerContainerFactory rabbitListenerContainerFactory(
            ConnectionFactory connectionFactory) {
        SimpleRabbitListenerContainerFactory factory = new SimpleRabbitListenerContainerFactory();
        factory.setConnectionFactory(connectionFactory);
        factory.setMessageConverter(jsonMessageConverter());
        return factory;
    }

    // ── Exchange ──────────────────────────────────────────────────────────────

    @Bean
    public DirectExchange orchestrationExchange() {
        return new DirectExchange(ORCHESTRATION_EXCHANGE);
    }

    // ── Queues ────────────────────────────────────────────────────────────────

    @Bean
    public Queue registerUserCmdQueue() {
        return QueueBuilder.durable(REGISTER_USER_CMD_QUEUE).build();
    }

    @Bean
    public Queue userRegisteredReplyQueue() {
        return QueueBuilder.durable(USER_REGISTERED_REPLY_QUEUE).build();
    }

    // ── Bindings ──────────────────────────────────────────────────────────────

    @Bean
    public Binding registerUserCmdBinding() {
        return BindingBuilder
                .bind(registerUserCmdQueue())
                .to(orchestrationExchange())
                .with(REGISTER_USER_CMD_KEY);
    }

    @Bean
    public Binding userRegisteredReplyBinding() {
        return BindingBuilder
                .bind(userRegisteredReplyQueue())
                .to(orchestrationExchange())
                .with(USER_REGISTERED_REPLY_KEY);
    }
}
