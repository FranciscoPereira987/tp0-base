package protocol

import "errors"

type Hello struct{
	clientID uint32
}

func (this *Hello) Serialize() []byte {

	stream := buildHeader([]byte{HELLO_OP}, 4)

	this.addClientID(&stream)

	return stream
}

func (this *Hello) Deserialize(stream []byte) error {
	if len(stream) != 8 {
		return errors.New("Invalid hello message")
	}
	stream = stream[HEADER_SIZE:]
	clientID, err := this.deserializeClientID(&stream)

	this.clientID = clientID

	return err
}

func (this *Hello) ShouldAck() bool {
	return true
}

func (this Hello) addClientID(stream *[]byte) {
	serializeUint32(stream, this.clientID)
}
func (this Hello) deserializeClientID(stream *[]byte) (uint32, error) {
	return deserializeUint32(stream)
}