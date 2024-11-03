package ping

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

const (
	pingInterval int = 1
	pingTimeout  int = 5
)

type client struct {
	Target     *string
	Interval   *int
	RetryCount *int
	pingClient *probing.Pinger
}

func NewClient() (*client, error) {
	target := flag.String("target", "", "target ip address to ping")
	interval := flag.Int("interval", 1, "interval in minutes between pings\n")
	retryCount := flag.Int("retry-count", 3, "retry count before instructing reboot")
	flag.Parse()

	if *target == "" {
		flag.PrintDefaults()
		return nil, fmt.Errorf("target ip not specified")
	}

	pingClient, err := probing.NewPinger(*target)
	if err != nil {
		return nil, fmt.Errorf("error during initialise: %v", err)
	}
	pingClient.Interval = time.Second * time.Duration(pingInterval)
	pingClient.Timeout = time.Second * time.Duration(pingTimeout)

	return &client{
		Target:     target,
		Interval:   interval,
		RetryCount: retryCount,
		pingClient: pingClient,
	}, nil
}

func (c *client) Ping() error {
	err := c.pingClient.Run()
	if err != nil {
		return err
	}

	if c.pingClient.Statistics().PacketLoss == 100 {
		return fmt.Errorf("100%% packet loss")
	}

	log.Printf("packet loss: %v", c.pingClient.Statistics().PacketLoss)
	return nil
}

func (c *client) PingUntilCancel(ctx context.Context) {
	fmt.Println("ping until cancel started...")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("ping until cancel stopped...")
			return
		case <-time.After(time.Second * 1):
			c.Ping()
		}
	}
}
