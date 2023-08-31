package protocol

import "errors"

type End struct{}

func (end *End) Serialize() []byte {
	return buildHeader([]byte{END_OP}, 0)
}

func (end *End) Deserialize(stream []byte) error {
	if !compareStreams(end.Serialize(), stream) {
		return errors.New("invalid end message")
	}
	return nil
}

func (end *End) ShouldAck() bool {
	return true
}
