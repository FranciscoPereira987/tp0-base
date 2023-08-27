package protocol

import "errors"

type Bet struct {
	name       string
	surname    string
	personalId string
	birthdate  string // yyyy-mm-dd

	betedNumber uint32
}

func (this *Bet) ShouldAck() bool {
	return true
}

func (this Bet) makeHeader() []byte {
	length := EXTRA_BET_BYTES +
		len(this.name) +
		len(this.surname) +
		len(this.personalId) +
		len(this.birthdate)

	header := []byte{BET_OP}

	return buildHeader(header, length)
}

func (this Bet) addField(stream *[]byte, field string) {
	*stream = append(*stream, byte(len(field)))

	*stream = append(*stream, []byte(field)...)
}

func (this Bet) addBetNumbet(stream *[]byte) {
	serializeUint32(stream, this.betedNumber)
}

func (this Bet) addBody(stream *[]byte) {
	this.addField(stream, this.name)
	this.addField(stream, this.surname)
	this.addField(stream, this.personalId)
	this.addField(stream, this.birthdate)

	this.addBetNumbet(stream)

}

func (this *Bet) Serialize() []byte {
	serialized := this.makeHeader()
	this.addBody(&serialized)
	return serialized
}

func (this Bet) deserializeField(fieldLength int, stream *[]byte) (string, error) {
	field := ""
	if len(*stream) < fieldLength {
		return field, errors.New("Stream is shorter than expected")
	}

	field = string((*stream)[:fieldLength])

	*stream = (*stream)[fieldLength:]

	return field, nil
}

func (this Bet) deserializeFieldLength(stream *[]byte) (int, error) {
	if len(*stream) == 0 {
		return 0, errors.New("Stream is shorter than expected")
	}
	length := int((*stream)[0])

	*stream = (*stream)[1:]

	return length, nil
}

func (this Bet) getFieldFromStream(stream *[]byte) (string, error) {
	fieldLength, err := this.deserializeFieldLength(stream)
	field := ""
	if err == nil {
		field, err = this.deserializeField(fieldLength, stream)
	}

	return field, err

}

func (this Bet) deserializeBet(stream *[]byte) (uint32, error) {
	return deserializeUint32(stream)
}

func (this *Bet) Deserialize(stream []byte) (err error) {
	fields := make([]string, 0)
	var field string
	var betedNumber uint32

	err = checkHeader(stream, BET_OP)

	if err != nil {
		return
	}

	stream = stream[HEADER_SIZE:]

	for i := 0; i < 4 && err == nil; i++ {
		field, err = this.getFieldFromStream(&stream)

		fields = append(fields, field)
	}

	if err == nil {
		betedNumber, err = this.deserializeBet(&stream)

		this.name = fields[0]
		this.surname = fields[1]
		this.personalId = fields[2]
		this.birthdate = fields[3]
		this.betedNumber = betedNumber
	}

	return err
}
