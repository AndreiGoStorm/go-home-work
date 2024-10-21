package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "connect timeout")
	flag.Parse()

	if len(flag.Args()) < 2 {
		log.Fatal("wrong number of parameters")
	}

	client := NewTelnetClient(net.JoinHostPort(flag.Arg(0), flag.Arg(1)), timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalf("telnet connection error: %v", err)
	}
	defer client.Close()

	log.Printf("successfully connected to host")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer cancel()

	go func() {
		if err := client.Send(); err != nil {
			log.Printf("telnet sending error: %v", err)
		}
		cancel()
	}()

	go func() {
		if err := client.Receive(); err != nil {
			log.Printf("telnet receiving error: %v", err)
		}
		cancel()
	}()

	<-ctx.Done()
}
