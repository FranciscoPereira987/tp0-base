package common

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/protocol"
)

var (
	BET_NUMBER_POS = 4
	NAME_POS       = 0
	SURNAME_POS    = 1
	ID_POS         = 2
	BIRTHDATE_POS  = 3
)

type BetReaderConfig struct {
	BetPath   string
	BetFile   string
	BatchSize int
}

type BetReader struct {
	config BetReaderConfig
	file   *csv.Reader
	open   bool
}

func NewBetReader(config BetReaderConfig, id int) (*BetReader, error) {
	fileName := fmt.Sprintf(config.BetPath+config.BetFile, id)
	file, err := os.Open(fileName)

	if err != nil {
		return nil, err
	}

	reader := new(BetReader)
	reader.config = config
	reader.file = csv.NewReader(file)
	reader.open = true

	return reader, nil
}

func (reader *BetReader) BetBatch() protocol.BetBatch {
	batch := protocol.BetBatch{}
	for i := 0; i < reader.config.BatchSize; i++ {
		record, err := reader.file.Read()
		if err != nil {
			reader.close()
			return batch
		}
		number, _ := strconv.Atoi(record[BET_NUMBER_POS])
		bet := protocol.Bet{
			Name:        record[NAME_POS],
			Surname:     record[SURNAME_POS],
			PersonalId:  record[ID_POS],
			Birthdate:   record[BIRTHDATE_POS],
			BetedNumber: uint32(number),
		}
		batch.Bets = append(batch.Bets, bet)
	}
	return batch
}

func (reader *BetReader) close() {
	reader.open = false
}

func (reader BetReader) BetsLeft() bool {
	return reader.open
}
