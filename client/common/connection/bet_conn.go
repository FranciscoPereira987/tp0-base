package connection

import (
	"errors"
	"net"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/protocol"
)

// Implements a socket abstraction to work with Bets
// Works with tcp
type BetConn struct {
	conn net.Conn
}

func (this *BetConn) helloServer() error {
	message := &protocol.Hello{}

	err := this.Write(message)

	if err != nil {
		return err
	}

	return err
}

func NewBetConn(addr string) (*BetConn, error) {

	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return nil, err
	}

	betConn := &BetConn{
		conn,
	}
	betConn.helloServer()
	return betConn, nil
}

func (this *BetConn) Close() error {
	this.Write(&protocol.End{})
	
	return this.conn.Close()
}

func (this *BetConn) peak() ([]byte, error) {
	header := make([]byte, 4)

	//TODO: Avoid short read
	_, err := this.conn.Read(header)

	return header, err
}

func (this *BetConn) readMessage(header []byte) ([]byte, error) {
	length, err := protocol.GetMessageLength(header)

	if err != nil {
		return nil, err
	}

	// TODO: Avoid short read
	message := make([]byte, length-4)
	readed, err := this.conn.Read(message)

	if readed != length-4 {
		return nil, errors.New("Malformed Message")
	}

	if err != nil {
		return nil, err
	}

	return append(header, message...), nil
}

func (this *BetConn) Write(message protocol.Message) error {
	stream := message.Serialize()
	writen, err := this.conn.Write(stream)

	if writen < len(stream) {
		return errors.New("Could not write all bytes into socket")
	}

	if err != nil {
		return err
	}

	if message.ShouldAck() {

		ack := &protocol.Ack{}
		err = this.Read(ack)
	}

	return err
}

func (this *BetConn) Read(message protocol.Message) error {
	header, err := this.peak()
	if err != nil {
		return err
	}

	stream, err := this.readMessage(header)

	if err != nil {
		return err
	}

	return message.Deserialize(stream)
}
