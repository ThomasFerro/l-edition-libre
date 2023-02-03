package commands

import "github.com/ThomasFerro/l-edition-libre/events"

type CancelManuscriptSubmission struct {
	ManuscriptName string
}

func HandleCancelManuscriptSubmission(command CancelManuscriptSubmission) ([]events.Event, error) {
	return []events.Event{
		events.ManuscriptSubmissionCanceled{},
	}, nil
}
