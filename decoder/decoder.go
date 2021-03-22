package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	// EXCHANGENAME is the name of the RabbitMQ exchange
	EXCHANGENAME = "go-test-exchange"
	// ROUTINGKEY is the name of the RabbitMQ routing key for routing messages
	ROUTINGKEY = "go-test-key"
	// QUEUENAME is the name of the RabbitMQ queue
	QUEUENAME = "go-test-queue"
	// CONSUMERNAME is the name of the RabbitMQ consumer
	CONSUMERNAME = "go-amqp-decoder"
)

var (
	amqpURI = flag.String("amqp", "amqp://guest:guest@localhost:5672/", "AMQP URI")
)

func init() {
	flag.Parse()
	initLog()
	initAmqp()
}

func initLog() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel) // TODO: configurable
}

func amqpFailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

var conn *amqp.Connection
var ch *amqp.Channel
var replies <-chan amqp.Delivery

func initAmqp() {
	var err error
	var q amqp.Queue

	conn, err = amqp.Dial(*amqpURI)
	amqpFailOnError(err, "Failed to connect to RabbitMQ")

	log.Infof("Got Connection, getting Channel...")

	ch, err = conn.Channel()
	amqpFailOnError(err, "Failed to open a channel")

	log.Infof("Got Channel, declaring Exchange (%s)", EXCHANGENAME)

	err = ch.ExchangeDeclare(
		EXCHANGENAME, // name of the exchange
		"direct",     // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	)
	amqpFailOnError(err, "Failed to declare the Exchange")

	log.Infof("Declared Exchange, declaring Queue (%s)", QUEUENAME)

	q, err = ch.QueueDeclare(
		QUEUENAME, // name, leave empty to generate a unique name
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	amqpFailOnError(err, "Error declaring the Queue")

	log.Infof("Declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		q.Name, q.Messages, q.Consumers, ROUTINGKEY)

	err = ch.QueueBind(
		q.Name,       // name of the queue
		ROUTINGKEY,   // bindingKey
		EXCHANGENAME, // sourceExchange
		false,        // noWait
		nil,          // arguments
	)
	amqpFailOnError(err, "Error binding to the Queue")

	log.Infof("Queue bound to Exchange, starting Consume (consumer tag %q)", CONSUMERNAME)

	replies, err = ch.Consume(
		q.Name,       // queue
		CONSUMERNAME, // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	amqpFailOnError(err, "Error consuming the Queue")
}

func main() {
	log.Infoln("Start decoder...")
	var count int = 1
	for r := range replies {
		log.Debugf("Consuming message number %d", count)
		log.Infof("[x] Received a message: %s", r.Body)
		count++
	}

	// Close Channel
	defer ch.Close()

	// Close Connection
	defer conn.Close()
}
