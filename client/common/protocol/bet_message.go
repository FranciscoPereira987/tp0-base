package protocol

import (
	"encoding/binary"
	"errors"
)

type Bet struct {
	Agency uint32
	Name       string
	Surname    string
	PersonalId string
	Birthdate  string // yyyy-mm-dd

	BetedNumber uint32
}

func (bet *Bet) ShouldAck() bool {
	return true
}

func (bet Bet) makeHeader() []byte {
	length := EXTRA_BET_BYTES +
		len(bet.Name) +
		len(bet.Surname) +
		len(bet.PersonalId) +
		len(bet.Birthdate)

	header := []byte{BET_OP}

	return buildHeader(header, length)
}

func (bet Bet) addField(stream *[]byte, field string) {
	*stream = append(*stream, byte(len(field)))

	*stream = append(*stream, []byte(field)...)
}



func (bet Bet) addBody(stream *[]byte) {
	bet.addField(stream, bet.Name)
	bet.addField(stream, bet.Surname)
	bet.addField(stream, bet.PersonalId)
	bet.addField(stream, bet.Birthdate)

	serializeUint32(stream, bet.BetedNumber)
	serializeUint32(stream, bet.Agency)
}

func (bet *Bet) Serialize() []byte {
	serialized := bet.makeHeader()
	bet.addBody(&serialized)
	return serialized
}

func (bet Bet) deserializeField(fieldLength int, stream *[]byte) (string, error) {
	field := ""
	if len(*stream) < fieldLength {
		return field, errors.New("stream is shorter than expected")
	}

	field = string((*stream)[:fieldLength])

	*stream = (*stream)[fieldLength:]

	return field, nil
}

func (bet Bet) deserializeFieldLength(stream *[]byte) (int, error) {
	if len(*stream) == 0 {
		return 0, errors.New("stream is shorter than expected")
	}
	fieldSize := []byte{0, (*stream)[0]}
	length := binary.BigEndian.Uint16(fieldSize)
	
	*stream = (*stream)[1:]

	return int(length), nil
}

func (bet Bet) getFieldFromStream(stream *[]byte) (string, error) {
	fieldLength, err := bet.deserializeFieldLength(stream)
	field := ""
	if err == nil {
		field, err = bet.deserializeField(fieldLength, stream)
	}

	return field, err

}


func (bet *Bet) Deserialize(stream []byte) (err error) {
	fields := make([]string, 0)
	var field string
	var betedNumber, agency uint32

	err = checkHeader(stream, BET_OP)

	if err != nil {
		return
	}

	stream = stream[HEADER_SIZE:]

	for i := 0; i < 4 && err == nil; i++ {
		field, err = bet.getFieldFromStream(&stream)

		fields = append(fields, field)
	}

	if err == nil {
		betedNumber, err = deserializeUint32(&stream)
		if err != nil {
			return errors.New("invalid bet structure")
		}
		agency, err = deserializeUint32(&stream)
		bet.Agency = agency
		bet.Name = fields[0]
		bet.Surname = fields[1]
		bet.PersonalId = fields[2]
		bet.Birthdate = fields[3]
		bet.BetedNumber = betedNumber
	}

	return err
}