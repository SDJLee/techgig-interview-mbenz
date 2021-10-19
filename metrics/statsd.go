package metrics

import (
	"fmt"
	"net"
	"os"
	"time"

	log "github.com/SDJLee/mercedes-benz/logger"
)

var statsLogger = log.SubLogger("statsd")

// makes a buffered queue of type string of max capacity = 100
var queue = make(chan string, 100)

func init() {
	go statsdSender()
}

// counter metric
func StatCount(metric string, value int) {
	queue <- fmt.Sprintf("%s:%d|c", metric, value)
}

// timer metric
func StatTime(metric string) func() {
	start := time.Now()
	return func() {
		queue <- fmt.Sprintf("%s:%d|ms", metric, time.Since(start)/1e6)
	}
}

// pushes metrics into graphite's udp port
func statsdSender() {
	for msg := range queue {
		conn, err := net.Dial("udp", os.Getenv("GRAPHITE_URL"))
		if err != nil {
			statsLogger.Error("could not connect to statsd")
			continue
		}

		_, err = conn.Write([]byte(msg + "\n"))
		if err != nil {
			statsLogger.Error(err)
		}
		_ = conn.Close()
	}
}
