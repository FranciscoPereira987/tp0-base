package common

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
	stopNotify chan os.Signal
	timerNotify <-chan time.Time
	running bool
	
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
		running: false,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Fatalf(
	        "action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	
	c.conn = conn
	return nil
}

func (c *Client) stop() {
	log.Infof("action: closing_notify_channel | result: in_progress")
	close(c.stopNotify)
	log.Infof("action: closing_notify_channel | result: success")
	c.running = false
}

//Returns if the client is running
//If the client is running, checks if a signal has been recieved to shut down the client
func (c *Client) isRunning() bool {
	if c.running {
		select {
		case <-c.stopNotify:
			log.Infof("action: SIGTERM_detected | result: success | client_id: %v",
				c.config.ID)
			c.stop()
		case <-c.timerNotify:
			log.Infof("action: timeout_detected | result: success | client_id: %v",
				c.config.ID,
			)
			c.stop()
		default:

		}
	}
	return c.running
}


//Sets the c.stopNotify channel and starts up manageStatus 
func (c *Client) setStatusManager() {
	
	stopNotify := make(chan os.Signal, 1)
	timerNotify := time.After(c.config.LoopLapse)

	signal.Notify(stopNotify, syscall.SIGTERM)
	c.stopNotify = stopNotify
	c.timerNotify = timerNotify
	c.running = true
	
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	// autoincremental msgID to identify every message sent
	msgID := 1
	c.setStatusManager()
	// Send messages if the loopLapse threshold has not been surpassed
	for c.isRunning() {

		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()

		// TODO: Modify the send to avoid short-write
		fmt.Fprintf(
			c.conn,
			"[CLIENT %v] Message N°%v\n",
			c.config.ID,
			msgID,
		)
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		msgID++
		log.Infof("action: closing_socker | result: in_progress")
		c.conn.Close()
		log.Infof("action: closing_socker | result: success")

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
                c.config.ID,
				err,
			)
			return
		}
		log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
            c.config.ID,
            msg,
        )

		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
