package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	// EXCHANGENAME is the name of the RabbitMQ deduplication exchange
	EXCHANGENAME = "deduplication-exchange"
	// ROUTINGKEY is the name of the RabbitMQ routing key for routing messages
	ROUTINGKEY = "deduplication-routing-key"
	// QUEUENAME is the name of the RabbitMQ queue
	QUEUENAME = "deduplication-queue"
	// CONSUMERNAME is the name of the RabbitMQ consumer
	CONSUMERNAME = "decoder"
	// DECODEREXCHANGENAME is the name of the RabbitMQ decoded exchange
	DECODEREXCHANGENAME = "decoder-exchange"
	// DECODERROUTINGKEY is the name of the RabbitMQ routing key for routing messages
	DECODERROUTINGKEY = "1"
)

var (
	amqpURI = flag.String("amqp", "amqp://guest:guest@localhost:5672/", "AMQP URI")
)

func init() {
	initLog()
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

var (
	conn     *amqp.Connection
	ch       *amqp.Channel
	confirms chan amqp.Confirmation
	replies  <-chan amqp.Delivery
)

func configureAmqp() {
	var err error
	var q amqp.Queue

	conn, err = amqp.Dial(*amqpURI)
	amqpFailOnError(err, "Failed to connect to RabbitMQ")

	log.Infof("Got Connection, getting Channel...")

	ch, err = conn.Channel()
	amqpFailOnError(err, "Failed to open a channel")

	log.Infof("Got Channel, declaring Exchange (%s)", EXCHANGENAME)

	var exchangeArgs = make(amqp.Table)
	exchangeArgs["x-cache-persistence"] = "disk"
	exchangeArgs["x-cache-size"] = "10000"
	exchangeArgs["x-cache-ttl"] = "300"

	err = ch.ExchangeDeclare(
		EXCHANGENAME,              // name of the exchange
		"x-message-deduplication", // type
		true,                      // durable
		false,                     // delete when complete
		false,                     // internal
		false,                     // noWait
		exchangeArgs,              // arguments
	)
	amqpFailOnError(err, "Failed to declare the Exchange")

	log.Infof("Declared Exchange, declaring Queue (%s)", QUEUENAME)

	var qArgs = make(amqp.Table)
	qArgs["x-queue-type"] = "quorum"
	qArgs["x-quorum-initial-group-size"] = 3
	qArgs["x-single-active-consumer"] = true

	q, err = ch.QueueDeclare(
		QUEUENAME, // name, leave empty to generate a unique name
		true,      // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		qArgs,     // arguments
	)
	amqpFailOnError(err, "Error declaring the Queue")

	err = ch.Qos(
		30, // prefetch count
		0,  // prefetch size
		false,
	)
	amqpFailOnError(err, "Error declaring QoS")

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
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	amqpFailOnError(err, "Error consuming the Queue")

	err = ch.ExchangeDeclare(
		DECODEREXCHANGENAME, // name of the exchange
		"x-consistent-hash", // type
		true,                // durable
		false,               // delete when complete
		false,               // internal
		false,               // noWait
		nil,                 // arguments
	)
	amqpFailOnError(err, "Failed to declare the Exchange")

	// TODO: notifypublish confirm/nacks
	// Buffer of 1 for our single outstanding publishing
	// confirms = ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	// if err := ch.Confirm(false); err != nil {
	// 	log.Fatalf("confirm.select destination: %s", err)
	// }
}

func publishMessages(messages []string) {
	for i := 0; i < len(messages); i++ {
		err := ch.Publish(
			DECODEREXCHANGENAME, // exchange
			DECODERROUTINGKEY,   // routing key
			false,               // mandatory
			false,               // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(messages[i]),
				Timestamp:    time.Now(),
			})
		amqpFailOnError(err, "Failed to Publish on RabbitMQ") // TODO: nacks notifypublish

		log.Infof("[x] Sent a message: %s", messages[i])
	}
}

var totalCount int = 1

func processDeliveries() {
	var count int = 1
	var deliveries []amqp.Delivery

	for r := range replies {
		log.Debugf("Consuming message number batch %d, total %d", count, totalCount)
		log.Infof("[x] Received a message: %s", r.Body)
		deliveries = append(deliveries, r)
		count++
		totalCount++

		if len(deliveries) == 30 {
			// Decode message batch of 30
			decodeMessageBatch(deliveries)
			for _, d := range deliveries {
				// Publish messages to decoder exchanger
				publishMessages([]string{string(d.Body)}) // TODO: Replace with "decoded" messages array
			}
			// Acknolwedge all messages up to batch and reset
			ch.Ack(deliveries[len(deliveries)-1].DeliveryTag, true)
			deliveries = nil
		}
	}
}

func decodeMessageBatch(d []amqp.Delivery) {
	// TODO: Decoding business logic
}

func main() {
	flag.Parse()
	log.Infoln("Start decoder...")

	configureAmqp()

	// Process delivered messages in batches of 30
	processDeliveries()

	// Close Channel
	defer ch.Close()

	// Close Connection
	defer conn.Close()
}
