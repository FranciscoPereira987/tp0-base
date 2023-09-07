package protocol

import "errors"


type Winners struct {}

func (win *Winners) ShouldAck() bool {
	return false
}

func (win *Winners) Serialize() []byte {
	return buildHeader([]byte{WINN_OP}, 0)
}

func (win *Winners) Deserialize(stream []byte) error{
	if !compareStreams(win.Serialize(), stream) {
		return errors.New("invalid Winners message")
	}
	return nil
}