package ping

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

const (
	pingTimeoutInSeconds  int = 5
	pingIntervalInSeconds int = 1
)

type client struct {
	Target     *string
	Interval   *int
	RetryCount *int
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

	return &client{
		Target:     target,
		Interval:   interval,
		RetryCount: retryCount,
	}, nil
}

func (c *client) Ping() error {
	pingClient, err := probing.NewPinger(*c.Target)
	if err != nil {
		return fmt.Errorf("error during initialise: %w", err)
	}

	pingClient.Interval = time.Second * time.Duration(pingIntervalInSeconds)
	pingClient.Timeout = time.Second * time.Duration(pingTimeoutInSeconds)

	err = pingClient.Run()
	if err != nil {
		return err
	}

	if pingClient.Statistics().PacketLoss == 100 {
		return fmt.Errorf("100%% packet loss")
	}

	pingClient.Stop()
	return nil
}

func (c *client) PingUntilCancel(ctx context.Context) {
	log.Println("ping until cancel started...")
	failures := 0
	for {
		select {
		case <-ctx.Done():
			log.Println("ping until cancel stopped...")
			return
		case <-time.After(time.Minute * time.Duration(*c.Interval)):
			err := c.Ping()
			if err != nil {
				failures += 1
				log.Printf("failure attempt %v out of %v with error: %v\n", failures, *c.RetryCount, err)
			} else if failures >= 1 {
				log.Printf("connectivity restored after %v failure attempts\n", failures)
				failures = 0
			}

			if failures >= *c.RetryCount {
				reboot()
			}
		}
	}
}

func reboot() {
	log.Println("rebooting system...")
	switch os := runtime.GOOS; os {

	case "linux":
		cmd := exec.Command("systemctl", "reboot")
		err := cmd.Run()
		if err != nil {
			log.Fatalf("error during rebooting of system: %v", err)
		}

	case "darwin":
		cmd := exec.Command("shutdown", "-r", "now")
		err := cmd.Run()
		if err != nil {
			log.Fatalf("error during rebooting of system: %v", err)
		}
	}
	os.Exit(0)
}
