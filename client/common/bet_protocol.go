package common

import "errors"

//OP codes definitions for Bet protocol messages
var (
	HELLO_OP = byte(0x01)
	ACK_OP = byte(0x02)
	ERR_OP = byte(0x03)
	BET_OP = byte(0x04)
	BETBATCH_OP = byte(0x05)
	END_OP = byte(0xff)

	EXTRA_BET_BYTES = 8
)


//Interface for serializing and deserializing Bet protocol messages
// Message header: |OP_CODE|MLength| => OP_CODE 1 byte; MLength 3 bytes
type Message interface {
	Serialize() []byte
	//Assumes that the OP code is already known
	Deserialize(stream []byte) error
}

func compareStreams(a []byte, b []byte) (result bool) {
	result = len(a) == len(b)

	for index := 0; result && index < len(a); index++{
		result = a[index] == b[index]
	}

	return
}

type Hello struct {}

func (this *Hello) Serialize() []byte {
	return []byte{HELLO_OP, 0, 0, 1}
}

func (this *Hello) Deserialize(stream []byte) error {
	if !compareStreams(this.Serialize(), stream) {
		return errors.New("Invalid hello message")
	}
	return nil
}

type Ack struct {}


func (this *Ack) Serialize() []byte {
	return []byte{ACK_OP, 0, 0, 1}
}


func (this *Ack) Deserialize(stream []byte) error {
	
	if !compareStreams(this.Serialize(), stream) {
		return errors.New("Invalid Ack message")
	}
	return nil
}


type Err struct {}


func (this *Err) Serialize() []byte {
	return []byte{ERR_OP, 0, 0, 1}
}


func (this *Err) Deserialize(stream []byte) error {
	if !compareStreams(this.Serialize(), stream){
		return errors.New("Invalid Err message")
	}
	return nil
}

type Bet struct {

	name string
	surname string
	personalId string
	birthdate string // yyyy-mm-dd

	betedNumber uint32
}

func (this Bet) makeHeader() []byte {
	length := EXTRA_BET_BYTES +
	 			len(this.name) + 
				len(this.surname) + 
				len(this.personalId) + 
				len(this.birthdate) + 
				4

	header := []byte{BET_OP}

	for shift := 16; shift >= 0; shift -= 8 {
		header = append(header, byte(length >> shift))
	} 

	return header
}

func (this Bet) addField(stream *[]byte, field string) {
	*stream = append(*stream, byte(len(field)))

	*stream = append(*stream, []byte(field)...)
} 

func (this Bet) addBetNumbet(stream *[]byte)  {
	serialized := make([]byte, 0)

	for shift := 24; shift >= 0; shift -= 8 {
		serialized = append(serialized, byte(this.betedNumber >> shift))
	}

	*stream = append(*stream, serialized...)
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
	var betNumber uint32
	
	if len(*stream) !=  4 {
		return betNumber, errors.New("Stream size is diferent than expected")
	}

	for index, value := range *stream {
		betNumber |= (uint32(value) << (8 * (3 - index)))
	}

	return betNumber, nil
 }

func (this *Bet) Deserialize(stream []byte) (err error) {
	fields := make([]string, 0)
	var field string
	var betedNumber uint32
	for i := 0; i < 4 && err == nil; i++{
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

type BetBatch struct {
	bets []Bet
}


func (this BetBatch) makeHeader(betsLength int) []byte {
	header := []byte{BETBATCH_OP}
	betsLength += 4
	for shift := 16; shift >= 0; shift -= 8 {
		header = append(header, byte(betsLength >> shift))
	}

	return header
}

func (this *BetBatch) Serialize() []byte {
	serializedBets := make([]byte, 0)

	for _, bet := range this.bets {
		serializedBets = append(serializedBets, bet.Serialize()...)
	}

	header := this.makeHeader(len(serializedBets))

	return append(header, serializedBets...)
}


func (this BetBatch) getBetStreamSize(stream []byte) (size int, err error) {
	if len(stream) < 4 {
		err = errors.New("Stream is shorter than expected")
		return
	}

	for i := 1; i < 4; i++ {
		size += int(stream[i]) << (8 * (3 - i))
	}

	return

}

func (this *BetBatch) deserializeBet(stream *[]byte) (err error) {
	
	betSize, err := this.getBetStreamSize(*stream)
	if err == nil {
		newBet := new(Bet)
		err = newBet.Deserialize((*stream)[4:betSize])
		this.bets = append(this.bets, *newBet)
		*stream = (*stream)[betSize:]
	}

	return
}

func (this *BetBatch) Deserialize(stream []byte) (err error) {
	
	for ;len(stream) > 0 && err == nil; {
		err = this.deserializeBet(&stream)
	}
	
	return 
}


type End struct {}


func (this *End) Serialize() []byte {
	return []byte{END_OP, 0, 0, 1}
}


func (this *End) Deserialize(stream []byte) error {
	if !compareStreams(this.Serialize(), stream) {
		return errors.New("Invalid end message")
	}
	return nil
}
