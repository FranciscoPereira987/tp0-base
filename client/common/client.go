package common

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/connection"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/protocol"
	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
	Bet protocol.Bet
}

// Client Entity that encapsulates how
type Client struct {
	config     ClientConfig
	conn       *connection.BetConn
	stopNotify <-chan bool
	running bool
	waitGroup sync.WaitGroup
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config:  config,
		running: false,
		waitGroup: sync.WaitGroup{},
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := connection.NewBetConn(c.config.ServerAddress, c.config.ID)
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

// Returns if the client is running
// If the client is running, checks if a signal has been recieved to shut down the client
func (c *Client) isRunning() bool {
	if c.running {
		select {
		case c.running = <-c.stopNotify:
		default:

		}
	}
	return c.running
}

// Sets the c.stopNotify channel and starts up manageStatus
func (c *Client) setStatusManager() {

	stopNotify := make(chan bool, 1)

	c.stopNotify = stopNotify
	c.running = true

	go c.manageStatus(stopNotify)
}

// Waits for either a message on listener or a timeout, then writes into stopNotify
func (c *Client) manageStatus(stopNotify chan<- bool) {
	c.waitGroup.Add(1)
	listener := make(chan os.Signal)
	timeout := time.After(c.config.LoopLapse)

	signal.Notify(listener, syscall.SIGTERM)

	defer c.waitGroup.Done()
	defer close(listener)
	defer close(stopNotify)

	select {
	case <-listener:
		log.Infof("action: SIGTERM_detected | result: success | client_id: %v",
			c.config.ID)
	case <-timeout:
		log.Infof("action: timeout_detected | result: success | client_id: %v",
			c.config.ID,
		)
	}
	stopNotify <- true
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	// autoincremental msgID to identify every message sent
	msgID := 1
	c.setStatusManager()
	// Send messages if the loopLapse threshold has not been surpassed
	for c.isRunning() {
		
		c.createClientSocket()
		// Create the connection the server in every loop iteration. Send an
		//c.createClientSocket()

		// TODO: Modify the send to avoid short-write
		err := c.conn.Write(&c.config.Bet) 
		msgID++
		c.conn.Close()

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}
		log.Infof("action: apuesta_enviada | result: success | dni: %s | numero: %d",
			c.config.Bet.PersonalId, c.config.Bet.BetedNumber)		
		c.conn.Close()
		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
	c.waitGroup.Wait()
}
