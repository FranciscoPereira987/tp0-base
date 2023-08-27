package protocol

import "errors"

type Err struct{}

func (this *Err) Serialize() []byte {
	return buildHeader([]byte{ERR_OP}, 0)
}

func (this *Err) Deserialize(stream []byte) error {
	if !compareStreams(this.Serialize(), stream) {
		return errors.New("Invalid Err message")
	}
	return nil
}

func (this *Err) ShouldAck() bool {
	return false
}
