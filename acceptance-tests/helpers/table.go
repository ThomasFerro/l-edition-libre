package helpers

import (
	"fmt"

	"github.com/cucumber/godog"
)

// FIXME: Laisser tomber cette façon de faire, on ne peut pas ajouter une méthode sur le code de prod.
// Essayer de faire avec un json.marshal pour le moment
type ComparableWithTableRow interface {
	Value(attribulte string) string
}

type tableHeader struct {
	index int
	value string
}

func ShouldMatch[T ComparableWithTableRow](actuals []T, expected *godog.Table) error {
	headers := make([]tableHeader, 0)
	for index, header := range expected.Rows[0].Cells {
		headers = append(headers, tableHeader{
			index: index,
			value: header.Value,
		})
	}
	for i := 1; i < len(expected.Rows); i++ {
		row := expected.Rows[i]
		actual := actuals[i-1]
		// TODO: Une version avec toutes les erreurs d'un coup
		// TODO: Essayer de faire matcher toutes les lignes sans se soucier de l'ordre
		for _, header := range headers {
			actualValue := actual.Value(header.value)
			if actualValue != row.Cells[header.index].Value {
				return fmt.Errorf("cell %v mismatch, expected %v but got %v", header.value, row.Cells[header.index].Value, actualValue)
			}
		}
	}
	return nil
}
