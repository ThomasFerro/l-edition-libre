package commands

import (
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/events"
)

type CancelManuscriptSubmission struct {
	events.ManuscriptID
}

func HandleCancelManuscriptSubmission(command CancelManuscriptSubmission) ([]events.Event, error) {
	return []events.Event{
		events.ManuscriptSubmissionCanceled{},
	}, nil
}

func (c CancelManuscriptSubmission) String() string {
	return fmt.Sprintf("CancelManuscriptSubmission{ManuscriptID %v}", c.ManuscriptID)
}
