package common

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common/protocol"
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
		number, _ := strconv.Atoi(record[4])
		bet := protocol.Bet{
			Name:        record[0],
			Surname:     record[1],
			PersonalId:  record[2],
			Birthdate:   record[3],
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
