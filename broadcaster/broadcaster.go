package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmichael/go-multicast/multicast"
	log "github.com/sirupsen/logrus"
)

var (
	multicastURI = flag.String("multicast", "239.0.0.0:9002", "UDP Multicast URI")
	messageCount = flag.Int("count", 0, "Amount of messages to broadcast, 0 for infinite")
	duplicate    = flag.Bool("duplicate", true, "Send the same messages twice to simulate deplication")
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

	var i int = 0
	for {
		ts := time.Now().Unix()
		stamp := fmt.Sprint(ts)
		msg := stamp + "-datatest" + fmt.Sprint(i)
		conn.Write([]byte(msg))
		if *duplicate {
			conn.Write([]byte(msg))
		}
		time.Sleep(1 * time.Second)
		log.Infof("[x] Broadcast a message: %s", msg)
		i++
		if i > *messageCount && *messageCount != 0 {
			break
		}
	}
}

func wait() {
	exitSignal := make(chan os.Signal)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
}

func main() {
	log.Infoln("Start broadcaster...")
	log.Infof("Broadcasting to %s\n", *multicastURI)
	ping(*multicastURI)
	wait()
}
