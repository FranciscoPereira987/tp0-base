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
	agency int
	config BetReaderConfig
	file   *csv.Reader
	fd     *os.File
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
	reader.agency = id

	return reader, nil
}

func (reader *BetReader) BetBatch() (protocol.BetBatch, error) {
	batch := protocol.BetBatch{}
	for i := 0; i < reader.config.BatchSize; i++ {
		record, err := reader.file.Read()
		if err != nil {
			return batch, reader.Close()
		}
		number, _ := strconv.Atoi(record[BET_NUMBER_POS])
		bet := protocol.Bet{
			Agency: uint32(reader.agency),
			Name:        record[NAME_POS],
			Surname:     record[SURNAME_POS],
			PersonalId:  record[ID_POS],
			Birthdate:   record[BIRTHDATE_POS],
			BetedNumber: uint32(number),
		}
		batch.Bets = append(batch.Bets, bet)
	}
	return batch, nil
}

func (reader *BetReader) Close() error {
	reader.open = false
	return reader.fd.Close()
}

func (reader BetReader) BetsLeft() bool {
	return reader.open
}
