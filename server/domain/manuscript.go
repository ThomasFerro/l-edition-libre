package domain

import (
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/events"
	"golang.org/x/exp/slog"
)

type Manuscript struct {
	Status ManuscriptStatus
	Title  string
	Author string
}

type ManuscriptStatus string

const (
	PendingReview ManuscriptStatus = "PendingReview"
	Canceled      ManuscriptStatus = "Canceled"
	Reviewed      ManuscriptStatus = "Reviewed"
)

func (m Manuscript) applyManuscriptSubmitted(event events.ManuscriptSubmitted) Manuscript {
	m.Title = event.Title
	m.Author = event.Author
	m.Status = PendingReview
	return m
}
func (m Manuscript) applyManuscriptSubmissionCanceled(event events.ManuscriptSubmissionCanceled) Manuscript {
	m.Status = Canceled
	return m
}
func (m Manuscript) applyManuscriptReviewed(event events.ManuscriptReviewed) Manuscript {
	m.Status = Reviewed
	return m
}
func (m Manuscript) String() string {
	return fmt.Sprintf("Manuscript{Title %v, Author %v, Status %v}", m.Title, m.Author, m.Status)
}

func RehydrateManuscript(history []events.Event) Manuscript {
	manuscript := Manuscript{}

	for _, nextEvent := range history {
		switch typedEvent := nextEvent.(type) {
		case events.ManuscriptSubmitted:
			manuscript = manuscript.applyManuscriptSubmitted(typedEvent)
		case events.ManuscriptSubmissionCanceled:
			manuscript = manuscript.applyManuscriptSubmissionCanceled(typedEvent)
		case events.ManuscriptReviewed:
			manuscript = manuscript.applyManuscriptReviewed(typedEvent)
		default:
			slog.Warn("unknown manuscript event", "event", typedEvent)
		}
	}

	return manuscript
}
