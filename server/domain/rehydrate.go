package domain

import (
	"github.com/ThomasFerro/l-edition-libre/events"
	"golang.org/x/exp/slog"
)

func (m Manuscript) applyManuscriptSubmitted(event events.ManuscriptSubmitted) Manuscript {
	m.Title = event.Title
	m.Author = event.Author
	m.Status = PendingReview
	m.FileURL = event.FileURL
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

func (p Publication) applyPublicationMadeAvailable(e events.PublicationMadeAvailable) Publication {
	p.Status = Available
	return p
}

func RehydratePublication(history []events.Event) Publication {
	publication := Publication{}

	for _, nextEvent := range history {
		switch typedEvent := nextEvent.(type) {
		case events.PublicationMadeAvailable:
			publication = publication.applyPublicationMadeAvailable(typedEvent)
		default:
			slog.Warn("unknown publication event", "event", typedEvent)
		}
	}

	return publication
}
