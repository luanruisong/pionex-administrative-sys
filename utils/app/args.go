package app

import (
	"flag"
	"fmt"
)

var (
	port    = flag.String("p", "8080", "port to listen on")
	daemon  = flag.Bool("d", false, "daemon process")
	fileLog = flag.Bool("fl", false, "file log")
)

func Parse() {
	flag.Parse()
}

func Port() string {
	return fmt.Sprintf(":%s", *port)
}

func Daemon() bool {
	return *daemon
}

func FileLog() bool {
	return *fileLog
}
