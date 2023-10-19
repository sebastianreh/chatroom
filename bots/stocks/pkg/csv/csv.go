package csv

import (
	"bytes"
	"encoding/csv"
	"errors"
)

func ReadStockCsv(CSV []byte) ([][]string, error) {
	var records [][]string
	reader := csv.NewReader(bytes.NewReader(CSV))
	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() != "EOF" {
				return records, err
			}
			break
		}
		records = append(records, record)
	}
	if len(records) < 2 {
		return records, errors.New("records len != 2")
	}
	return records, nil
}
