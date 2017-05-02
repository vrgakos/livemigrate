package main

import (
	"flag"
	"github.com/vrgakos/livemigrate/tcpapp"
	"fmt"
	"os"
	"os/signal"
	"os/exec"
)

var port int
var bind string

func init() {
	flag.IntVar(&port, "port", 1234, "Listening port")
	flag.StringVar(&bind, "bind", "0.0.0.0", "Binging ip")
}

func main() {
	flag.Parse()

	server := tcpapp.NewTcpServer(fmt.Sprintf("%s:%d", bind, port))
	server.Start()

	stressArgs := os.Getenv("STRESS")
	var cmd *exec.Cmd
	if len(stressArgs) > 0 {
		cmd = exec.Command("/usr/bin/stress", stressArgs)
		cmd.Stdout = os.Stdout
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	for _ = range signalChan {

	}
	server.Stop()

	if cmd != nil {
		cmd.Process.Kill()
	}
}
