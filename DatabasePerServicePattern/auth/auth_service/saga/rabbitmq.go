package saga

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ queue i exchange nazivi za registraciju korisnika
const (
	OrchestrationExchange = "saga.registration.exchange"

	// Komande (auth -> stakeholder)
	RegisterUserCmdQueue = "register.user.command.queue"
	RegisterUserCmdKey   = "register.user.command"

	// Kompenzacione komande (auth -> auth, rollback)
	DeleteUserCmdQueue = "delete.user.command.queue"
	DeleteUserCmdKey   = "delete.user.command"

	// Reply queues (stakeholder -> auth)
	UserRegisteredReplyQueue = "user.registered.reply.queue"
	UserRegisteredReplyKey   = "user.registered.reply"

	UserDeletedReplyQueue = "user.deleted.reply.queue"
	UserDeletedReplyKey   = "user.deleted.reply"
)

// RabbitMQConnection drži konekciju i kanal prema RabbitMQ brokeru
type RabbitMQConnection struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// NewRabbitMQConnection otvara konekciju, kreira kanal i deklariše sve exchange-ove i queue-ove
func NewRabbitMQConnection(url string) (*RabbitMQConnection, error) {
	const (
		maxAttempts = 30
		retryDelay   = 2 * time.Second
	)

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		conn, err := amqp.Dial(url)
		if err != nil {
			lastErr = fmt.Errorf("rabbitmq dial failed: %w", err)
			log.Printf("[RabbitMQ] connect attempt %d/%d failed: %v", attempt, maxAttempts, lastErr)
		} else {
			ch, err := conn.Channel()
			if err != nil {
				conn.Close()
				lastErr = fmt.Errorf("rabbitmq channel open failed: %w", err)
				log.Printf("[RabbitMQ] connect attempt %d/%d failed: %v", attempt, maxAttempts, lastErr)
			} else {
				r := &RabbitMQConnection{Conn: conn, Channel: ch}
				if err := r.declareTopology(); err != nil {
					r.Close()
					lastErr = err
					log.Printf("[RabbitMQ] connect attempt %d/%d failed: %v", attempt, maxAttempts, lastErr)
				} else {
					log.Println("[RabbitMQ] Connected and topology declared")
					return r, nil
				}
			}
		}

		if attempt < maxAttempts {
			time.Sleep(retryDelay)
		}
	}

	return nil, fmt.Errorf("rabbitmq connection failed after %d attempts: %w", maxAttempts, lastErr)
}

// declareTopology deklariše exchange, queue-ove i binding-e
func (r *RabbitMQConnection) declareTopology() error {
	// Direct exchange za orchestration
	if err := r.Channel.ExchangeDeclare(
		OrchestrationExchange, "direct", true, false, false, false, nil,
	); err != nil {
		return fmt.Errorf("exchange declare failed: %w", err)
	}

	queues := []struct{ name, key string }{
		{RegisterUserCmdQueue, RegisterUserCmdKey},
		{DeleteUserCmdQueue, DeleteUserCmdKey},
		{UserRegisteredReplyQueue, UserRegisteredReplyKey},
		{UserDeletedReplyQueue, UserDeletedReplyKey},
	}

	for _, q := range queues {
		if _, err := r.Channel.QueueDeclare(q.name, true, false, false, false, nil); err != nil {
			return fmt.Errorf("queue declare %s failed: %w", q.name, err)
		}
		if err := r.Channel.QueueBind(q.name, q.key, OrchestrationExchange, false, nil); err != nil {
			return fmt.Errorf("queue bind %s failed: %w", q.name, err)
		}
	}
	return nil
}

// Publish šalje JSON poruku na zadati exchange sa zadatim routing key-em
func (r *RabbitMQConnection) Publish(routingKey string, body []byte) error {
	return r.Channel.Publish(
		OrchestrationExchange,
		routingKey,
		false, false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
}

// Consume vraća kanal poruka sa zadatog queue-a
func (r *RabbitMQConnection) Consume(queueName string) (<-chan amqp.Delivery, error) {
	return r.Channel.Consume(queueName, "", false, false, false, false, nil)
}

// Close zatvara kanal i konekciju
func (r *RabbitMQConnection) Close() {
	if r.Channel != nil {
		r.Channel.Close()
	}
	if r.Conn != nil {
		r.Conn.Close()
	}
}
