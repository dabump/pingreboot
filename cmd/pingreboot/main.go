package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/dabump/pingreboot/internal/ping"
	"github.com/dabump/pingreboot/internal/util"
)

const (
	version string = "0.0.1alpha"
)

func main() {
	log.Printf("running pingreboot %v\n", version)

	isRoot, err := util.IsRoot()
	if err != nil {
		log.Fatalf("error during determining current user: %v", err)
	}
	if !isRoot {
		log.Fatalf("pingreboot must be run as root")
	}

	client, err := ping.NewClient()
	if err != nil {
		log.Fatalf("unable to initialise: %v", err)
	}

	err = client.Ping()
	if err != nil {
		log.Fatalf("failure with initial ping request: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		client.PingUntilCancel(ctx)
	}()

	cancelSignal := make(chan os.Signal, 1)
	signal.Notify(cancelSignal, os.Interrupt)
	<-cancelSignal
	log.Println("Interrupted. Shutting down...")
	cancel()
}
