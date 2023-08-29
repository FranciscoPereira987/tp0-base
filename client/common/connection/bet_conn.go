package connection

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/protocol"
)

// Implements a socket abstraction to work with Bets
// Works with tcp
type BetConn struct {
	conn net.Conn
	active bool
	id int
}



func NewBetConn(addr string, id string) (*BetConn, error) {

	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return nil, err
	}

	parsedId, err := strconv.Atoi(id)
	if err != nil{
		return nil, err
	}

	betConn := &BetConn{
		conn,
		false,
		parsedId,
	}
	betConn.helloServer()
	return betConn, nil
}

func (this *BetConn) Close() error {
	this.Write(&protocol.End{})
	return this.shutdown()
}


func (this *BetConn) Write(message protocol.Message) error {
	if !this.active{
		return errors.New("Connection closed")
	}

	stream := message.Serialize()
	
	err := this.writeBytes(stream)

	if message.ShouldAck() {

		ack := &protocol.Ack{}
		err = this.Read(ack)
	}

	return err
}

func (this *BetConn) Read(message protocol.Message) error {
	if !this.active{
		return errors.New("Connection closed")
	}
	header, err := this.peak()
	if err != nil {
		return err
	}

	err = this.readMessage(header, message)

	return err 
}

func (this *BetConn) helloServer() error {
	message := &protocol.Hello{
		ClientID: uint32(this.id),
	}
	this.active = true

	err := this.Write(message)

	if err != nil {
		return err
	}


	return err
}


func (this *BetConn) shutdown() error {
	if !this.active{
		return nil
	}
	this.active = false
	return this.conn.Close()
}

func (this *BetConn) peak() ([]byte, error) {
	return this.readBytes(4)
}

func (this *BetConn) readBytes(bytes int) ([]byte, error) {
	if bytes < 0 {
		return nil, errors.New(fmt.Sprintf("Invalid read amount: %d", bytes))
	} 

	buff := make([]byte, bytes)
	readed, err := this.conn.Read(buff)
	
	var chunk_size int
	for readed < bytes && err == nil{
		chunk_size, err = this.conn.Read(buff[readed:])
		if chunk_size == 0 {
			err = errors.New("Broken connection")
		}
		readed += chunk_size
	}
	

	return buff, err
}

func (this *BetConn) writeBytes(bytes []byte) error {
	var chunk_size int
	writen, err := this.conn.Write(bytes)

	for writen < len(bytes) && err == nil {
		chunk_size, err = this.conn.Write(bytes[writen:])
		if chunk_size == 0{
			err = errors.New("Broken connection")
		}
		writen += chunk_size
	}

	return err
}

func (this *BetConn) manageInvalidMessage(message []byte, original error) error {
	end := new(protocol.End)
	if end.Deserialize(message) == nil {
		ack := new(protocol.Ack)
		this.Write(ack)
		this.shutdown()
		return errors.New("Connection closed by server")
	}
	err := new(protocol.Err)
	if err.Deserialize(message) == nil {
		return errors.New("Recieved Err message")
	}
	return original
}

func (this *BetConn) readMessage(header []byte, expected protocol.Message) error {
	length, err := protocol.GetMessageLength(header)

	if err != nil {
		return err
	}

	if length < 4 {
		return errors.New("Malformed message read")
	}

	message, err := this.readBytes(length-4) 

	if err != nil {
		return  err
	}
	message = append(header, message...)
	
	err = expected.Deserialize(message)
	
	if err != nil {
		err = this.manageInvalidMessage(message, err)
	}

	return err
}
