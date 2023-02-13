package helpers

import (
	"github.com/cucumber/godog"
)

type tableHeader struct {
	index int
	value string
}

func ExtractData(table *godog.Table) []map[string]string {
	headers := make([]tableHeader, 0)
	for index, header := range table.Rows[0].Cells {
		headers = append(headers, tableHeader{
			index: index,
			value: header.Value,
		})
	}
	extracted := []map[string]string{}
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		extractedRow := make(map[string]string)
		for _, header := range headers {
			extractedRow[header.value] = row.Cells[header.index].Value
		}
		extracted = append(extracted, extractedRow)
	}
	return extracted
}
