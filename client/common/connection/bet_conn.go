package connection

import (
	"errors"
	"fmt"
	"net"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/protocol"
)

// Implements a socket abstraction to work with Bets
// Works with tcp
type BetConn struct {
	conn   net.Conn
	active bool
	id     int
}

func NewBetConn(addr string, id int) (*BetConn, error) {

	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return nil, err
	}

	betConn := &BetConn{
		conn,
		false,
		id,
	}
	betConn.helloServer()
	return betConn, nil
}

func (conn *BetConn) Close() error {
	conn.Write(&protocol.End{})
	return conn.shutdown()
}

func (conn *BetConn) Write(message protocol.Message) error {
	if !conn.active {
		return errors.New("connection closed")
	}

	stream := message.Serialize()

	err := conn.writeBytes(stream)

	if message.ShouldAck() {

		ack := &protocol.Ack{}
		err = conn.Read(ack)
	}

	return err
}

func (conn *BetConn) Read(message protocol.Message) error {
	if !conn.active {
		return errors.New("connection closed")
	}
	header, err := conn.peak()
	if err != nil {
		return err
	}

	err = conn.readMessage(header, message)

	return err
}

func (conn *BetConn) helloServer() error {
	message := &protocol.Hello{
		ClientID: uint32(conn.id),
	}
	conn.active = true

	err := conn.Write(message)

	if err != nil {
		return err
	}

	return err
}

func (conn *BetConn) shutdown() error {
	if !conn.active {
		return nil
	}
	conn.active = false
	return conn.conn.Close()
}

func (conn *BetConn) peak() ([]byte, error) {
	return conn.readBytes(protocol.HEADER_SIZE)
}

func (conn *BetConn) readBytes(bytes int) ([]byte, error) {
	if bytes < 0 {
		return nil, fmt.Errorf("invalid read amount: %d", bytes)
	}

	buff := make([]byte, bytes)
	readed, err := conn.conn.Read(buff)

	var chunk_size int
	for readed < bytes && err == nil {
		chunk_size, err = conn.conn.Read(buff[readed:])
		if chunk_size == 0 {
			err = errors.New("broken connection")
		}
		readed += chunk_size
	}

	return buff, err
}

func (conn *BetConn) writeBytes(bytes []byte) error {
	var chunk_size int
	writen, err := conn.conn.Write(bytes)

	for writen < len(bytes) && err == nil {
		chunk_size, err = conn.conn.Write(bytes[writen:])
		if chunk_size == 0 {
			err = errors.New("broken connection")
		}
		writen += chunk_size
	}

	return err
}

func (conn *BetConn) manageInvalidMessage(message []byte, original error) error {

	if conn.isEndMessage(message) {
		conn.sendAck()
		conn.shutdown()
		return errors.New("connection closed by server")
	}

	if conn.isErrMessage(message) {
		return errors.New("recieved Err message")
	}
	return original
}

func (conn *BetConn) isEndMessage(message []byte) bool {
	end := new(protocol.End)
	return conn.isMessage(message, end)
}

func (conn *BetConn) isErrMessage(message []byte) bool {
	err := new(protocol.Err)
	return conn.isMessage(message, err)
}

func (conn *BetConn) isMessage(message []byte, messageType protocol.Message) bool {
	return messageType.Deserialize(message) == nil
}

func (conn *BetConn) sendAck() error {
	ack := new(protocol.Ack)
	return conn.Write(ack)
}

func (conn *BetConn) readMessage(header []byte, expected protocol.Message) error {
	length, err := protocol.GetMessageLength(header)

	if err != nil {
		return err
	}

	if length < protocol.HEADER_SIZE {
		return errors.New("malformed message read")
	}

	message, err := conn.readBytes(length - protocol.HEADER_SIZE)

	if err != nil {
		return err
	}
	message = append(header, message...)

	err = expected.Deserialize(message)

	if err != nil {
		err = conn.manageInvalidMessage(message, err)
	}

	if err == nil && expected.ShouldAck() {
		err = conn.sendAck()
	}

	return err
}
