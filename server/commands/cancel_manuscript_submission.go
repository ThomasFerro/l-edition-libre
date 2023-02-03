package commands

import "github.com/ThomasFerro/l-edition-libre/events"

type CancelManuscriptSubmission struct {
	events.ManuscriptID
}

func HandleCancelManuscriptSubmission(command CancelManuscriptSubmission) ([]events.Event, error) {
	return []events.Event{
		events.ManuscriptSubmissionCanceled{},
	}, nil
}
