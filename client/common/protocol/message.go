package protocol

import (
	"encoding/binary"
	"errors"
)

// OP codes definitions for Bet protocol messages
var (
	HELLO_OP    = byte(0x01)
	ACK_OP      = byte(0x02)
	ERR_OP      = byte(0x03)
	BET_OP      = byte(0x04)
	BETBATCH_OP = byte(0x05)
	END_OP      = byte(0xff)

	EXTRA_BET_BYTES = 12

	HEADER_SIZE = 4
)

// Interface for serializing and deserializing Bet protocol messages
// Message header: |OP_CODE|MLength| => OP_CODE 1 byte; MLength 3 bytes
type Message interface {
	Serialize() []byte
	//Assumes that the OP code is already known
	Deserialize(stream []byte) error

	//Returns wether the message type should be acked
	ShouldAck() bool
}

func compareStreams(a []byte, b []byte) (result bool) {
	result = len(a) == len(b)

	for index := 0; result && index < len(a); index++ {
		result = a[index] == b[index]
	}

	return
}

func buildHeader(header []byte, size int) []byte {
	size += HEADER_SIZE
	totalSize := make([]byte, 4)
	binary.BigEndian.PutUint32(totalSize, uint32(size))
	totalSize[0] = header[0]
	return totalSize
}

func GetMessageLength(stream []byte) (size int, err error) {
	if len(stream) < HEADER_SIZE {
		err = errors.New("invalid stream for message")
		return
	}
	header := stream[0]
	stream[0] = 0
	size = int(binary.BigEndian.Uint32(stream))
	stream[0] = header
	return
}

func checkHeader(stream []byte, op_code byte) (err error) {
	length, err := GetMessageLength(stream)

	if err != nil {
		return
	}

	if stream[0] != op_code || length != len(stream) {
		err = errors.New("invalid header for message")
	}

	return
}

func deserializeUint32(stream *[]byte) (uint32, error) {
	var number uint32

	if len(*stream) < 4 {
		return number, errors.New("stream size is diferent than expected")
	}

	number = binary.BigEndian.Uint32((*stream)[:HEADER_SIZE])
	*stream = (*stream)[HEADER_SIZE:]
	return number, nil
}

func serializeUint32(stream *[]byte, number uint32) {
	bytenumber := make([]byte, 4)
	binary.BigEndian.PutUint32(bytenumber, number)
	*stream = append(*stream, bytenumber...)
}
