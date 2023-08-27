package protocol

import "errors"

type Ack struct{}

func (this *Ack) Serialize() []byte {
	return buildHeader([]byte{ACK_OP}, 0)
}

func (this *Ack) Deserialize(stream []byte) error {

	if !compareStreams(this.Serialize(), stream) {
		return errors.New("Invalid Ack message")
	}
	return nil
}

func (this *Ack) ShouldAck() bool {
	return false
}
