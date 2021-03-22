package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/dmichael/go-multicast/multicast"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	// EXCHANGENAME is the name of the RabbitMQ exchange
	EXCHANGENAME = "go-test-exchange"
	// ROUTINGKEY is the name of the RabbitMQ routing key for routing messages
	ROUTINGKEY = "go-test-key"
	// CONSUMERNAME is the name of the RabbitMQ consumer
	CONSUMERNAME = "go-amqp-decoder"
)

var (
	amqpURI      = flag.String("amqp", "amqp://guest:guest@localhost:5672/", "AMQP URI")
	multicastURI = flag.String("multicast", "239.0.0.0:9002", "UDP Multicast URI")
)

func init() {
	flag.Parse()
	initLog()
	initAmqp()
	initUDPListener()
}

func initLog() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel) // TODO: configurable
}

var conn *amqp.Connection
var ch *amqp.Channel

func initAmqp() {
	var err error

	conn, err = amqp.Dial(*amqpURI)
	amqpFailOnError(err, "Failed to connect to RabbitMQ")

	ch, err = conn.Channel()
	amqpFailOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		EXCHANGENAME, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // noWait
		nil,          // arguments
	)
	amqpFailOnError(err, "Failed to declare the Exchange")
}

func amqpFailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func publishMessages(messages []string) {
	for i := 0; i < len(messages); i++ {
		err := ch.Publish(
			EXCHANGENAME, // exchange
			ROUTINGKEY,   // routing key
			false,        // mandatory
			false,        // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(messages[i]),
				Timestamp:    time.Now(),
			})
		amqpFailOnError(err, "Failed to Publish on RabbitMQ")
		log.Infof("[x] Sent a message: %s", messages[i])
	}
}

func msgHandler(src *net.UDPAddr, n int, b []byte) {
	log.Debugln(n, "bytes read from", src)
	msgs := []string{string(b[:n])}
	publishMessages(msgs)
}

func initUDPListener() {
	log.Infof("Listening on %s", *multicastURI)
	multicast.Listen(*multicastURI, msgHandler)
}

func main() {
	log.Println("Starting listener...")

	// Close Channel
	defer ch.Close()

	// Close Connection
	defer conn.Close()
}

// confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))
// if err := ch.Confirm(false); err != nil {
// 	log.Fatalf("confirm.select destination: %s", err)
// }

// if confirmed := <-confirms; confirmed.Ack {
// 	msg.Ack(false)
// } else {
// 	msg.Nack(false, false)
// }

// func confirm(ack, nack chan int64, bodies, resends chan []byte, limit int) {
//     var (
//         sequence int64
//         pending  = make(map[int64][]body)
//         flow     = bodies
//     )
//     for {
//         if len(pending) > limit {
//             flow = nil
//         } else {
//             flow = bodies
//         }
//         select {
//         case body := <-flow:
//             sequence++
//             pending[sequence] = body
//         case ack := <-acks:
//             delete(pending, ack)
//         case nack := <-nacks:
//             if body := pending[nack]; body != nil {
//                 resends <- body
//                 delete(pending, nack)
//             }
//         }
//     }
// }
