package protocol

import "testing"

type testCase struct {
	name     string
	message  Message
	expected []byte
}

func compareBet(a *Bet, b *Bet) (result bool) {

	result = a.BetedNumber == b.BetedNumber
	result = a.Agency == b.Agency
	result = result && a.Birthdate == b.Birthdate
	result = result && a.Name == b.Name
	result = result && a.PersonalId == b.PersonalId
	result = result && a.Surname == b.Surname

	return
}

func compareBetBatch(a *BetBatch, b *BetBatch) (result bool) {
	result = true
	for i := 0; i < len(a.Bets) && result; i++ {
		result = result && compareBet(&a.Bets[i], &b.Bets[i])
	}

	return
}

func Test4ByteMessagesSerialization(t *testing.T) {
	tests := []testCase{
		{
			"Hello Message",
			&Hello{
				3,
			},
			[]byte{HELLO_OP, 0, 0, 8, 0, 0, 0, 3},
		},
		{
			"Ack Message",
			&Ack{},
			[]byte{ACK_OP, 0, 0, 4},
		},
		{
			"Err Message",
			&Err{},
			[]byte{ERR_OP, 0, 0, 4},
		},
		{
			"End Message",
			&End{},
			[]byte{END_OP, 0, 0, 4},
		},
	}

	for _, testcase := range tests {
		serialized := testcase.message.Serialize()

		if !compareStreams(serialized, testcase.expected) {
			t.Errorf("FAILED: %s", testcase.name)
		}
	}
}

func Test4ByteMessagesDeserialization(t *testing.T) {
	tests := []testCase{
		{
			"Hello Message",
			&Hello{},
			[]byte{HELLO_OP, 0, 0, 1},
		},
		{
			"Ack Message",
			&Ack{},
			[]byte{ACK_OP, 0, 2, 0},
		},
		{
			"Err Message",
			&Err{},
			[]byte{ERR_OP, 3, 0, 0},
		},
		{
			"End Message",
			&End{},
			[]byte{END_OP, 1, 0, 1},
		},
	}

	for _, testcase := range tests {
		if testcase.message.Deserialize(testcase.expected) == nil {
			t.Errorf("FAILED: %s", testcase.name)
		}
	}
}

func TestBetAndBetBatchHeaders(t *testing.T) {
	bet := &Bet{
		2,
		"Francisco",
		"Pereira",
		"41797243",
		"1998-12-17",
		12345,
	}

	tests := []testCase{
		{
			"Bet Header",
			bet,
			[]byte{BET_OP, 0, 0, 50},
		},
		{
			"Bet batch header",
			&BetBatch{
				[]Bet{*bet, *bet, *bet, *bet, *bet, *bet, *bet, *bet},
			},
			[]byte{BETBATCH_OP, 0, 1, 148},
		},
	}

	for _, testcase := range tests {
		header := testcase.message.Serialize()

		if !compareStreams(header[:4], testcase.expected) {
			t.Errorf("FAILED: %s", testcase.name)
		}
	}
}

func TestBetDeserialization(t *testing.T) {
	betTest := &Bet{
		2,
		"Francisco",
		"Pereira",
		"41797243",
		"1998-12-17",
		12345,
	}

	serialized := betTest.Serialize()

	deserialized := new(Bet)

	deserialized.Deserialize(serialized)

	if !compareBet(betTest, deserialized) {
		t.Errorf("FAILED: Bet deserialization",)
	}
}

func TestBetBatchDeserialization(t *testing.T) {
	bet := Bet{
		3,
		"Francisco",
		"Pereira",
		"41797243",
		"1998-12-17",
		12345,
	}

	batchTest := &BetBatch{
		[]Bet{bet, bet, bet, bet, bet},
	}

	serialized := batchTest.Serialize()

	deserialized := new(BetBatch)

	deserialized.Deserialize(serialized)

	if len(deserialized.Bets) != len(batchTest.Bets) {
		t.Errorf("FAILED: Batch because of size: %d != %d",
			len(deserialized.Bets),
			len(batchTest.Bets))
	}

	if !compareBetBatch(batchTest, deserialized) {
		t.Errorf("FAILED: Batch")
	}
}
