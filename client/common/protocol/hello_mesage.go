package protocol

import "errors"

type Hello struct {
	ClientID uint32
}

func (h *Hello) Serialize() []byte {

	stream := buildHeader([]byte{HELLO_OP}, 4)

	h.addClientID(&stream)

	return stream
}

func (h *Hello) Deserialize(stream []byte) error {
	if len(stream) != 8 {
		return errors.New("invalid hello message")
	}
	stream = stream[HEADER_SIZE:]
	clientID, err := h.deserializeClientID(&stream)

	h.ClientID = clientID

	return err
}

func (h *Hello) ShouldAck() bool {
	return true
}

func (h Hello) addClientID(stream *[]byte) {
	serializeUint32(stream, h.ClientID)
}
func (h Hello) deserializeClientID(stream *[]byte) (uint32, error) {
	return deserializeUint32(stream)
}
