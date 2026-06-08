package choreography

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ─────────────────────────────────────────────────────────────────────────────
// Koreografija koristi POSEBAN exchange i queue-ove od orkestrisanog pristupa.
// Tako oba pristupa mogu raditi ISTOVREMENO bez međusobnog ometanja.
//
// Exchange:    saga.choreography.exchange  (DirectExchange)
//
// Queues:
//   auth.user.created.queue           ← prima stakeholder (od auth-a)
//   stakeholder.profile.created.queue ← prima auth (od stakeholder-a) → uspeh
//   auth.user.rollback.queue          ← prima auth (od stakeholder-a) → pad
// ─────────────────────────────────────────────────────────────────────────────

const (
	ChoreographyExchange = "saga.choreography.exchange"

	// Auth objavljuje, stakeholder prima
	AuthUserCreatedQueue = "auth.user.created.queue"
	AuthUserCreatedKey   = "auth.user.created"

	// Stakeholder objavljuje na uspeh, auth prima
	StakeholderProfileCreatedQueue = "stakeholder.profile.created.queue"
	StakeholderProfileCreatedKey   = "stakeholder.profile.created"

	// Stakeholder objavljuje na pad, auth prima (kompenzacija)
	AuthUserRollbackQueue = "auth.user.rollback.queue"
	AuthUserRollbackKey   = "auth.user.rollback"
)

// RabbitMQConnection drži konekciju i kanal prema RabbitMQ brokeru
// (isti broker kao za orchestration, ali drugi exchange i queue-ovi)
type RabbitMQConnection struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

// NewRabbitMQConnection otvara konekciju, kreira kanal i deklariše topologiju
// specifičnu za koreografisani SAGA pattern.
func NewRabbitMQConnection(url string) (*RabbitMQConnection, error) {
	const (
		maxAttempts = 30
		retryDelay  = 2 * time.Second
	)

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		conn, err := amqp.Dial(url)
		if err != nil {
			lastErr = fmt.Errorf("rabbitmq dial failed: %w", err)
			log.Printf("[CHOREOGRAPHY][RabbitMQ] attempt %d/%d: %v", attempt, maxAttempts, lastErr)
		} else {
			ch, err := conn.Channel()
			if err != nil {
				conn.Close()
				lastErr = fmt.Errorf("channel open failed: %w", err)
				log.Printf("[CHOREOGRAPHY][RabbitMQ] attempt %d/%d: %v", attempt, maxAttempts, lastErr)
			} else {
				r := &RabbitMQConnection{Conn: conn, Channel: ch}
				if err := r.declareTopology(); err != nil {
					r.Close()
					lastErr = err
					log.Printf("[CHOREOGRAPHY][RabbitMQ] attempt %d/%d: %v", attempt, maxAttempts, lastErr)
				} else {
					log.Println("[CHOREOGRAPHY][RabbitMQ] Connected and topology declared")
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

// declareTopology deklariše exchange, queue-ove i binding-e za koreografiju
func (r *RabbitMQConnection) declareTopology() error {
	// DirectExchange za koreografisani pattern
	if err := r.Channel.ExchangeDeclare(
		ChoreographyExchange, "direct", true, false, false, false, nil,
	); err != nil {
		return fmt.Errorf("exchange declare failed: %w", err)
	}

	queues := []struct{ name, key string }{
		{AuthUserCreatedQueue, AuthUserCreatedKey},
		{StakeholderProfileCreatedQueue, StakeholderProfileCreatedKey},
		{AuthUserRollbackQueue, AuthUserRollbackKey},
	}

	for _, q := range queues {
		if _, err := r.Channel.QueueDeclare(q.name, true, false, false, false, nil); err != nil {
			return fmt.Errorf("queue declare %s failed: %w", q.name, err)
		}
		if err := r.Channel.QueueBind(q.name, q.key, ChoreographyExchange, false, nil); err != nil {
			return fmt.Errorf("queue bind %s failed: %w", q.name, err)
		}
	}
	return nil
}

// Publish šalje JSON poruku na choreography exchange sa zadatim routing key-em
func (r *RabbitMQConnection) Publish(routingKey string, body []byte) error {
	return r.Channel.Publish(
		ChoreographyExchange,
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
