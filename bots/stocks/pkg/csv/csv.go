package csv

import (
	"bytes"
	"encoding/csv"
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
	return records, nil
}
