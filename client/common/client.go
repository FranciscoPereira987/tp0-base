package common

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/connection"
	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            int
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
	Reader *BetReader
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   *connection.BetConn
	stopNotify chan os.Signal
	running bool
	
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config:    config,
		running:   false,
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


func (c *Client) stop() {
	log.Infof("action: stop | result: in_progress | comment: closing_stop_channel")
	close(c.stopNotify)
	log.Infof("action: stop | result: in_progress | comment: clossing_connection")
	c.conn.Close()
	log.Infof("action: stop | result: in_progress | comment: clossing_reader")
	c.config.Reader.Close()
	log.Infof("action stop | result: success")
 	c.running = false
}

func (c *Client) stopIfRunning() {
	if c.isRunning(){
		c.stop()
	}
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
		default:
		}
	}
	return c.running
}


// Sets the c.stopNotify channel and starts up manageStatus
func (c *Client) setStatusManager() {

	
	stopNotify := make(chan os.Signal, 1)

	signal.Notify(stopNotify, syscall.SIGTERM)

	c.stopNotify = stopNotify
	c.running = true
}


func (c *Client) waitForWinners() {
	c.conn.Close()
	waiter := connection.NewWinnersConn(c.config.ServerAddress, c.config.ID, c.config.LoopLapse)
	response, err := waiter.WaitForWinners()
	if err == nil {
		log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %d", response.TotalWinners())
	}else {
		log.Infof("action: consulta_ganadores | result: failed | error: %s", err)
	}
}

func (c *Client) sendBetBatch() error{
	batch, err := c.config.Reader.BetBatch()

	if err != nil {
		log.Errorf("action: batch_read | result: error | info: %s", err)
			
	}

	return c.conn.Write(&batch)
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	// autoincremental msgID to identify every message sent
	msgID := 1

	c.setStatusManager()
	c.createClientSocket()

	// Send messages if the loopLapse threshold has not been surpassed
	for c.isRunning() && c.config.Reader.BetsLeft(){
		
		
		// Create the connection the server in every loop iteration. Send an
		err := c.sendBetBatch()
		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			c.stopIfRunning()
			return
		}
		log.Infof("action: apuesta_enviada | result: success | batch_number: %d",
			msgID)				
		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)
		msgID++
	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
	if c.isRunning() {
		c.waitForWinners()
	}
	c.stopIfRunning()
}
