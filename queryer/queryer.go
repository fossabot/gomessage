package main

import (
	"crypto/tls"
	"flag"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	log "github.com/sirupsen/logrus"
)

var (
	influxURI   = flag.String("influxdb", "http://localhost:8086", "InfluxDB URI")
	influxToken = flag.String("influxdb-authtoken", "admin:admin", "InfluxDB authentication token (optional)")
)

func init() {
	initLog()
}

func initLog() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel) // TODO: configurable
}

var client influxdb2.Client

func configureInfluxDB() {
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

func cleanup() {
	// Close Client
	log.Infof("Closing client connection")
	defer client.Close()
}

func main() {
	flag.Parse()
	log.Infoln("Start reporter...")

	configureInfluxDB()

	// Watch for CTRL+C / SIGTERM
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup()
		os.Exit(1)
	}()

	// TODO: API Server
}
