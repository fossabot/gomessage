package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	batchinsert "github.com/rmeharg/gomessage/writer/batchinsert" // UGLY
)

const (
	// EXCHANGENAME is the name of the RabbitMQ deduplication exchange
	EXCHANGENAME = "decoder-exchange"
	// ROUTINGKEY is the name of the RabbitMQ routing key for routing messages
	ROUTINGKEY = "1"
	// QUEUENAME is the name of the RabbitMQ queue
	QUEUENAME = "decoder-queue"
	// CONSUMERNAME is the name of the RabbitMQ consumer
	CONSUMERNAME = "writer"
)

var (
	amqpURI            = flag.String("amqp", "amqp://user:CHANGEME@localhost:5672/", "AMQP URI")
	influxURI          = flag.String("influxdb", "http://localhost:8086", "InfluxDB URI")
	influxToken        = flag.String("influxdb-authtoken", "admin:admin", "InfluxDB authentication token (optional)")
	clickhouseURI      = flag.String("clickhouse", "127.0.0.1:9000", "Clickhouse URI")
	clickhouseUser     = flag.String("clickhouse-user", "clickhouse_operator", "Clickhouse username")
	clickhousePassword = flag.String("clickhouse-password", "clickhouse_operator_password", "Clickhouse password")
)

func init() {
	flag.Parse()
	initLog()
	initAmqp()
	initInfluxDB()
	initClickhouse()
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
	conn    *amqp.Connection
	ch      *amqp.Channel
	replies <-chan amqp.Delivery
)

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
		EXCHANGENAME,        // name of the exchange
		"x-consistent-hash", // type
		true,                // durable
		false,               // delete when complete
		false,               // internal
		false,               // noWait
		nil,                 // arguments
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

	log.Infof("Initialized AMQP RabbitMQ successfully")
}

var (
	client   influxdb2.Client
	writeAPI api.WriteAPI
	errorsCh <-chan error
)

func initInfluxDB() {
	// Create a new client using an InfluxDB server base URL and an authentication token
	client = influxdb2.NewClientWithOptions(*influxURI, *influxToken,
		influxdb2.DefaultOptions().
			SetUseGZip(true).
			SetTLSConfig(&tls.Config{
				InsecureSkipVerify: true,
			}))

	// // Get non-blocking write client for database
	writeAPI = client.WriteAPI("", "testdata")

	errorsCh = writeAPI.Errors()

	// Create go proc for reading and logging errors
	go func() {
		for err := range errorsCh {
			log.Errorf("write error: %s\n", err.Error())
			// fmt.Printf("write error: %s\n", err.Error())
		}
	}()
	log.Infof("Initialized InfluxDB successfully")
}

func parseTime(s string) time.Time {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(i, 0)
}

func writeDataMessage(b []byte) {
	message := strings.Split(string(b), "-")
	stamp := message[0]
	data := message[1]

	// create point
	p := influxdb2.NewPoint("stat",
		map[string]string{
			"id":     "test",
			"name":   "test",
			"source": "255.255.255.255",
		},
		map[string]interface{}{
			"message": data,
		},
		parseTime(stamp))

	// write asynchronously
	writeAPI.WritePoint(p)
	log.Infof("[x] Writing to TSDB: %s %s", stamp, data)

	// Force all unwritten data to be sent
	writeAPI.Flush()

	log.Infof("Wrote data message %s to InfluxDB successfully", b)
}

// var (
// 	connect *sql.DB
// 	tx      *sql.Tx
// 	stmt    *sql.Stmt
// )

var batch batchinsert.BatchInsert

func initClickhouse() {
	createTable := `CREATE TABLE IF NOT EXISTS element (
		id      String,
		name    String,
		source  String,
		message String,
		time    DateTime
	) engine=Memory`

	insertSQL := "INSERT INTO element (id, name, source, message, time) VALUES (?, ?, ?, ?, ?)"

	batch, err := batchinsert.New(
		insertSQL,
		batchinsert.WithHost(*clickhouseURI),
		batchinsert.WithUserInfo(*clickhouseUser, *clickhousePassword),
	)
	if err != nil {
		panic(err)
	}

	_, err = batch.DB().Exec(createTable)
	if err != nil {
		panic(err)
	}

	// connect, err := sql.Open("clickhouse", *clickhouseURI+"&username="+*clickhouseUser+"&password="+*clickhousePassword)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if err := connect.Ping(); err != nil {
	// 	if exception, ok := err.(*clickhouse.Exception); ok {
	// 		fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
	// 	} else {
	// 		fmt.Println(err)
	// 	}
	// 	return
	// }

	// _, err = connect.Exec(`
	// 	CREATE TABLE IF NOT EXISTS element (
	// 		id      String,
	// 		name    String,
	// 		source  String,
	// 		message String,
	// 		time    DateTime
	// 	) engine=Memory
	// `)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// tx, _ = connect.Begin()
	// stmt, _ = tx.Prepare("INSERT INTO element (id, name, source, message, time) VALUES (?, ?, ?, ?, ?)")

	log.Infof("Initialized Clickhouse successfully")
}

func writeDataMessageClickhouse(b []byte) {
	message := strings.Split(string(b), "-")
	stamp := message[0]
	data := message[1]

	// Execute statement
	batch.Insert(
		"test",
		"test",
		"255.255.255.255",
		data,
		parseTime(stamp),
	)

	// if _, err := tx.Stmt(stmt).Exec(
	// 	"test",
	// 	"test",
	// 	"255.255.255.255",
	// 	data,
	// 	parseTime(stamp),
	// ); err != nil {
	// 	log.Fatal(err)
	// }

	// // Commit statement to SQL DB
	// if err := tx.Commit(); err != nil {
	// 	log.Fatal(err)
	// }

	log.Infof("Wrote data message %s to Clickhouse successfully", b)

	// TODO: Move query to reporter
	// rows, err := connect.Query("SELECT id, name, source, message, time FROM element")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()

	// for rows.Next() {
	// 	var (
	// 		id      string
	// 		name    string
	// 		source  string
	// 		message string
	// 		time    time.Time
	// 	)
	// 	if err := rows.Scan(&id, &name, &source, &message, &time); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.Printf("id: %s, name: %d, source: %d, message: %v, time: %s", id, name, source, message, time)
	// }
	// if err := rows.Err(); err != nil {
	// 	log.Fatal(err)
	// }
	// if _, err := connect.Exec("DROP TABLE example"); err != nil {
	// 	log.Fatal(err)
	// }
}

func main() {
	log.Infoln("Start writer...")

	var count int = 1
	for r := range replies {
		log.Debugf("Consuming message number %d", count)
		log.Infof("[x] Received a message: %s", r.Body)
		r.Ack(false)

		//TODO: Decoding business logic

		writeDataMessage(r.Body)
		writeDataMessageClickhouse(r.Body)

		count++
	}

	// Close AMQP RabbitMQ Channel
	defer ch.Close()

	// Close AMQP RabbitMQ Connection
	defer conn.Close()

	// Close InfluxDB Client
	defer client.Close()

	// Close batch Clickhouse client
	defer batch.Close()

	// Close Clickhouse statement
	// defer stmt.Close()

	// // Close Clickhouse Connection
	// defer connect.Close()
}
