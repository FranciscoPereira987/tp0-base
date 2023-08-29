package protocol

type BetBatch struct {
	bets []Bet
}

func (this *BetBatch) Serialize() []byte {
	serializedBets := make([]byte, 0)

	for _, bet := range this.bets {
		serializedBets = append(serializedBets, bet.Serialize()...)
	}

	header := this.makeHeader(len(serializedBets))

	return append(header, serializedBets...)
}

func (this *BetBatch) Deserialize(stream []byte) (err error) {

	err = checkHeader(stream, BETBATCH_OP)

	if err != nil {
		return
	}

	stream = stream[HEADER_SIZE:]
	
	for len(stream) > 0 && err == nil {
		err = this.deserializeBetMessage(&stream)
	}
	return
}

func (this *BetBatch) ShouldAck() bool {
	return true
}

func (this BetBatch) makeHeader(betsLength int) []byte {
	header := []byte{BETBATCH_OP}
	return buildHeader(header, betsLength)
}

func (this BetBatch) getBetStreamSize(stream []byte) (int, error) {
	return GetMessageLength(stream)
}

func (this *BetBatch) deserializeBetMessage(stream *[]byte) (err error) {

	betSize, err := this.getBetStreamSize(*stream)
	if err == nil {
		newBet := new(Bet)
		err = newBet.Deserialize((*stream)[:betSize])
		this.bets = append(this.bets, *newBet)
		*stream = (*stream)[betSize:]
	}
	return
}
