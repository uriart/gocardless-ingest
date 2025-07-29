package clients

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gocardless-ingest/internal/models"
	"github.com/streadway/amqp"
)

type RabbitMQClient struct {
	connStr string
}

// NewRabbitMQClient crea una nueva instancia con la cadena de conexión AMQP
func NewRabbitMQClient(connStr string) *RabbitMQClient {
	return &RabbitMQClient{connStr: connStr}
}

// SendTransaction envía una transacción a la cola de RabbitMQ
func (r *RabbitMQClient) SendTransaction(tx models.Transaction) error {
	conn, err := amqp.Dial(r.connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"transactions-queue", // nombre
		true,                 // durable
		false,                // auto-delete
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	body, err := json.Marshal(tx)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}

	err = ch.Publish(
		"",                   // exchange (default)
		"transactions-queue", // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Println("✅ Sent transaction to queue:", tx.ID)
	return nil
}
