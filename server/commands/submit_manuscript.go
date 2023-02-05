package commands

import (
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/events"
)

type SubmitManuscript struct {
	ManuscriptName string
}

func HandleSubmitManuscript(command SubmitManuscript) ([]events.Event, error) {
	return []events.Event{
		events.ManuscriptSubmitted{},
	}, nil
}

func (c SubmitManuscript) String() string {
	return fmt.Sprintf("SubmitManuscript{ManuscriptName %v}", c.ManuscriptName)
}
