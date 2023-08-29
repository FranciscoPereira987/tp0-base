package protocol

import "errors"

type End struct{}

func (this *End) Serialize() []byte {
	return buildHeader([]byte{END_OP}, 0)
}

func (this *End) Deserialize(stream []byte) error {
	if !compareStreams(this.Serialize(), stream) {
		return errors.New("Invalid end message")
	}
	return nil
}

func (this *End) ShouldAck() bool {
	return true
}
