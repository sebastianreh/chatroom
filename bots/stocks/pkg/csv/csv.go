package csv

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/sebastianreh/chatroom-bots/stocks/entities"
)

func ReadStockCsv(CSV []byte) (entities.Stock, error) {
	var records [][]string
	reader := csv.NewReader(bytes.NewReader(CSV))
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Println("error reading CSV record", err)
			}
			break
		}
		records = append(records, record)
	}
	return entities.MapRecordsToStock(records)
}
