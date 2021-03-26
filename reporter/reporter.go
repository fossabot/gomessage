package main

import (
	"context"
	"crypto/tls"
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/jung-kurt/gofpdf"
	log "github.com/sirupsen/logrus"
)

var (
	influxURI   = flag.String("influxdb", "http://localhost:8086", "InfluxDB URI")
	influxToken = flag.String("influxdb-authtoken", "admin:admin", "InfluxDB authentication token (optional)")
)

func init() {
	flag.Parse()
	initLog()
	initInfluxDB()
}

func initLog() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel) // TODO: configurable
}

var client influxdb2.Client

func initInfluxDB() {
	// Create a new client using an InfluxDB server base URL and an authentication token
	client = influxdb2.NewClientWithOptions(*influxURI, *influxToken,
		influxdb2.DefaultOptions().
			SetUseGZip(true).
			SetTLSConfig(&tls.Config{
				InsecureSkipVerify: true,
			}))
}

func parseTime(s string) time.Time {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(i, 0)
}

func readDataMessage() []string {
	queryAPI := client.QueryAPI("")
	// get QueryTableResult

	var rows []string
	result, err := queryAPI.Query(context.Background(), `from(bucket:"testdata")|> range(start: -24h) |> filter(fn: (r) => r._measurement == "stat")`)
	if err == nil {
		// Iterate over query response
		for result.Next() {
			// Notice when group key has changed
			if result.TableChanged() {
				log.Infof("table: %s\n", result.TableMetadata().String())
			}
			// Access data
			log.Infof("value: %v\n", result.Record().Value())
			rows = append(rows, result.Record().String())
		}
		// check for an error
		if result.Err() != nil {
			log.Infof("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}

	return rows
}

func createPDF(rows []string) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	for _, row := range rows { 
		pdf.Cell(40, 10, row)
	}
	err := pdf.OutputFileAndClose("hello.pdf")
	if err != nil {
		log.Fatalf("error creating PDF file: %s\n", err.Error())
	}
}

func cleanup() {
	// Close Client
	log.Infof("Closing client connection")
	defer client.Close()
}

func main() {
	log.Infoln("Start reporter...")

	// Watch for CTRL+C / SIGTERM
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	for {
		time.Sleep(1 * time.Minute)
		data := readDataMessage()
		createPDF(data)
	}
}
