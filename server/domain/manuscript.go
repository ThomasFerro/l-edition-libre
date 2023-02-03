package domain

import (
	"github.com/ThomasFerro/l-edition-libre/events"
	"golang.org/x/exp/slog"
)

type Manuscript struct {
	Status
}

type Status string

const (
	PendingReview Status = "PendingReview"
	Canceled      Status = "Canceled"
)

func (m Manuscript) applyManuscriptSubmitted(event events.ManuscriptSubmitted) Manuscript {
	m.Status = PendingReview
	return m
}
func (m Manuscript) applyManuscriptSubmissionCanceled(event events.ManuscriptSubmissionCanceled) Manuscript {
	m.Status = Canceled
	return m
}

func Rehydrate(history []events.Event) Manuscript {
	manuscript := Manuscript{}

	for _, nextEvent := range history {
		switch typedEvent := nextEvent.(type) {
		case events.ManuscriptSubmitted:
			manuscript = manuscript.applyManuscriptSubmitted(typedEvent)
		case events.ManuscriptSubmissionCanceled:
			manuscript = manuscript.applyManuscriptSubmissionCanceled(typedEvent)
		default:
			slog.Warn("unknown manuscript event", "event", typedEvent)
		}
	}

	return manuscript
}
