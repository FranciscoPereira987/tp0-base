package protocol

import (
	"encoding/binary"
	"errors"
)

type WinnersResponse struct {
	winners []string
}

func (wr *WinnersResponse) ShouldAck() bool {
	return false
}

func (wr WinnersResponse) addField(stream *[]byte, field string) {
	*stream = append(*stream, byte(len(field)))

	*stream = append(*stream, []byte(field)...)
}

func (wr *WinnersResponse) Serialize() []byte {
	body := make([]byte, 0)
	for _, document := range wr.winners{
		wr.addField(&body, document)
	}
	header := buildHeader([]byte{WINNRESP_OP}, len(body))
	return append(header, body...)
}

func (wr WinnersResponse) deserializeFieldLength(stream *[]byte) (int, error) {
	if len(*stream) == 0 {
		return 0, errors.New("stream is shorter than expected")
	}
	fieldSize := []byte{0, (*stream)[0]}
	length := binary.BigEndian.Uint16(fieldSize)
	
	*stream = (*stream)[1:]

	return int(length), nil
}

func (wr WinnersResponse) deserializeField(fieldLength int, stream *[]byte) (string, error) {
	field := ""
	if len(*stream) < fieldLength {
		return field, errors.New("stream is shorter than expected")
	}

	field = string((*stream)[:fieldLength])

	*stream = (*stream)[fieldLength:]

	return field, nil
}


func (wr WinnersResponse) getFieldFromStream(stream *[]byte) (string, error) {
	fieldLength, err := wr.deserializeFieldLength(stream)
	field := ""
	if err == nil {
		field, err = wr.deserializeField(fieldLength, stream)
	}

	return field, err

}

func (wr *WinnersResponse) Deserialize(stream []byte) error {
	err := checkHeader(stream, WINNRESP_OP)

	if err != nil {
		return err
	}

	wr.winners = nil
	stream = stream[HEADER_SIZE:]
	for len(stream) > 0 {
		field, err := wr.getFieldFromStream(&stream)
		if err != nil {
			return err
		}
		wr.winners = append(wr.winners, field)
	}

	return nil
}