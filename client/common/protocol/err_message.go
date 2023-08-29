package protocol

import "errors"

type Err struct{}

func (err *Err) Serialize() []byte {
	return buildHeader([]byte{ERR_OP}, 0)
}

func (err *Err) Deserialize(stream []byte) error {
	if !compareStreams(err.Serialize(), stream) {
		return errors.New("invalid Err message")
	}
	return nil
}

func (err *Err) ShouldAck() bool {
	return false
}
