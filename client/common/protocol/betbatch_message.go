package protocol

type BetBatch struct {
	Bets []Bet
}

func (betb *BetBatch) Serialize() []byte {
	serializedBets := make([]byte, 0)

	for _, bet := range betb.Bets {
		serializedBets = append(serializedBets, bet.Serialize()...)
	}

	header := betb.makeHeader(len(serializedBets))

	return append(header, serializedBets...)
}

func (betb *BetBatch) Deserialize(stream []byte) (err error) {

	err = checkHeader(stream, BETBATCH_OP)

	if err != nil {
		return
	}

	stream = stream[HEADER_SIZE:]

	for len(stream) > 0 && err == nil {
		err = betb.deserializeBetMessage(&stream)
	}
	return
}

func (betb *BetBatch) ShouldAck() bool {
	return true
}

func (betb BetBatch) makeHeader(betsLength int) []byte {
	header := []byte{BETBATCH_OP}
	return buildHeader(header, betsLength)
}

func (betb BetBatch) getBetStreamSize(stream []byte) (int, error) {
	return GetMessageLength(stream)
}

func (betb *BetBatch) deserializeBetMessage(stream *[]byte) (err error) {

	betSize, err := betb.getBetStreamSize(*stream)
	if err == nil {
		newBet := new(Bet)
		err = newBet.Deserialize((*stream)[:betSize])
		betb.Bets = append(betb.Bets, *newBet)
		*stream = (*stream)[betSize:]
	}
	return
}
