package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/dmichael/go-multicast/multicast"
	log "github.com/sirupsen/logrus"
)

var (
	multicastURI = flag.String("multicast", "239.0.0.0:9002", "UDP Multicast URI")
)

func init() {
	flag.Parse()
	initLog()
}

func initLog() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel) // TODO: configurable
}

func ping(addr string) {
	conn, err := multicast.NewBroadcaster(addr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		ts := time.Now().Unix()
		stamp := fmt.Sprint(ts)
		conn.Write([]byte(stamp + "-datatest"))
		time.Sleep(1 * time.Second)
		log.Infof("[x] Broadcast a message:  %s", stamp+"-datatest")
	}
}

func main() {
	log.Infoln("Start broadcaster...")
	log.Infof("Broadcasting to %s\n", *multicastURI)
	ping(*multicastURI)
}
