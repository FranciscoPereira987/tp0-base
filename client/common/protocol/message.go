package protocol

import "errors"

// OP codes definitions for Bet protocol messages
var (
	HELLO_OP    = byte(0x01)
	ACK_OP      = byte(0x02)
	ERR_OP      = byte(0x03)
	BET_OP      = byte(0x04)
	BETBATCH_OP = byte(0x05)
	END_OP      = byte(0xff)

	EXTRA_BET_BYTES = 8

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
	for shift := 16; shift >= 0; shift -= 8 {
		header = append(header, byte(size>>shift))
	}

	return header
}

func GetMessageLength(stream []byte) (size int, err error) {
	if len(stream) < HEADER_SIZE {
		err = errors.New("Invalid stream for message")
		return
	}

	for i := 1; i < 4; i++ {
		size += int(stream[i]) << (8 * (3 - i))
	}

	return
}

func checkHeader(stream []byte, op_code byte) (err error) {
	length, err := GetMessageLength(stream)

	if err != nil {
		return
	}

	if stream[0] != op_code || length != len(stream) {
		err = errors.New("Invalid header for message")
	}

	return
}

func deserializeUint32(stream *[]byte) (uint32, error) {
	var number uint32

	if len(*stream) != 4 {
		return number, errors.New("Stream size is diferent than expected")
	}

	for index, value := range *stream {
		number |= (uint32(value) << (8 * (3 - index)))
	}

	return number, nil
}

func serializeUint32(stream *[]byte, number uint32) {
	serialized := make([]byte, 0)

	for shift := 24; shift >= 0; shift -= 8 {
		serialized = append(serialized, byte(number>>shift))
	}

	*stream = append(*stream, serialized...)
}