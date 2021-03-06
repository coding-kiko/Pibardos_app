package user

import (
	// Internal
	"encoding/json"
	"os"

	// Internal
	"github.com/coding-kiko/user_service/pkg/errors"
	"github.com/coding-kiko/user_service/pkg/log"

	// Third party
	"github.com/streadway/amqp"
)

var (
	UsersQueue  = os.Getenv("USERS_QUEUE")
	AvatarQueue = os.Getenv("AVATAR_QUEUE")
)

// Describes any consumer for a rabbit queue

type rabbitConsumer struct {
	conn    *amqp.Connection
	service Service
	logger  log.Logger
}

type RabbitConsumer interface {
	UsersQueue()
}

func NewRabbitConsumer(service Service, conn *amqp.Connection, logger log.Logger) RabbitConsumer {
	return &rabbitConsumer{
		conn:    conn,
		service: service,
		logger:  logger,
	}
}

// Receives new users or user update from the authentication service
func (r *rabbitConsumer) UsersQueue() {
	// create rabbit connection channel
	ch, err := r.conn.Channel()
	if err != nil {
		r.logger.Error("main.go", "main", err.Error())
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		UsersQueue, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		r.logger.Error("rabbitmq.go", "NewUsersQueueConsumer", err.Error())
		panic(err.Error())
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		r.logger.Error("rabbitmq.go", "NewUsersQueueConsumer", err.Error())
		panic(err.Error())
	}
	r.logger.Info("rabbitmq.go", "UsersQueue", "Started listening on usersQueue")
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			r.logger.Info("rabbitmq.go", "UsersQueue", "Received a message: "+string(d.Body))

			req := UpsertUserRequest{}
			err = json.Unmarshal(d.Body, &req)
			if err != nil {
				r.logger.Error("rabbitmq.go", "NewUsersQueueConsumer", err.Error())
			}
			_, err = r.service.UpsertUser(req)
			if err != nil {
				r.logger.Error("rabbitmq.go", "NewUsersQueueConsumer", err.Error())
			}
		}
	}()
	<-forever
}

// Describes any producer for a rabbit queue
type rabbitProducer struct {
	conn   *amqp.Connection
	logger log.Logger
}

type RabbitProducer interface {
	AvatarQueue(file File) error
}

func NewRabbitProducer(conn *amqp.Connection, logger log.Logger) RabbitProducer {
	return &rabbitProducer{
		conn:   conn,
		logger: logger,
	}
}

// send new avatar to be stored in the avatar service
func (r *rabbitProducer) AvatarQueue(file File) error {
	// create rabbit connection channel
	ch, err := r.conn.Channel()
	if err != nil {
		return errors.NewRabbitError("failed to create channel")
	}
	defer ch.Close()

	// declare queue
	q, err := ch.QueueDeclare(AvatarQueue, true, false, false, false, nil)
	if err != nil {
		return errors.NewRabbitError("failed to declare a queue")
	}

	headers := make(map[string]interface{})
	headers["filename"] = file.Name
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			Headers:      headers,
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         file.Data,
		})
	if err != nil {
		return errors.NewRabbitError("failed to publish message")
	}
	r.logger.Info("rabbitmq.go", "AvatarQueue", "avatar published")
	return nil
}
