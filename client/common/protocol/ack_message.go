package protocol

import "errors"

type Ack struct{}

func (ack *Ack) Serialize() []byte {
	return buildHeader([]byte{ACK_OP}, 0)
}

func (ack *Ack) Deserialize(stream []byte) error {

	if !compareStreams(ack.Serialize(), stream) {
		return errors.New("invalid Ack message")
	}
	return nil
}

func (ack *Ack) ShouldAck() bool {
	return false
}
