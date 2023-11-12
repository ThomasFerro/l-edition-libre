package mongodb

import (
	"context"
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/events"
	"go.mongodb.org/mongo-driver/bson"
)

type PublicationsHistory struct {
	history EventsHistory[string, PublicationEvent]
}

const publicationsEventsCollectionName = "publications_events"

type PublicationEvent struct {
	PublicationID string `bson:"publicationId"`
	UserID        string `bson:"userId"`
	EventType     string `bson:"eventType"`
}

func (event PublicationEvent) StreamID() (string, error) {
	return event.PublicationID, nil
}

func (publications PublicationsHistory) For(publicationID application.PublicationID) ([]application.ContextualizedEvent, error) {
	events, err := publications.history.ForSingleStream(publicationID.String(), bson.D{})
	if err != nil {
		return nil, err
	}

	contextualizedEvents := []application.ContextualizedEvent{}

	for _, nextEvent := range events {
		publicationEvent, err := toPublicationEvent(nextEvent)
		if err != nil {
			return nil, err
		}
		contextualizedEvents = append(contextualizedEvents, application.ContextualizedEvent{
			OriginalEvent: publicationEvent,
			Context: application.EventContext{
				UserID: contexts.UserID(nextEvent.UserID),
			},
		})
	}

	return contextualizedEvents, nil
}

func toPublicationEvent(nextDecodedEvent PublicationEvent) (events.Event, error) {
	switch nextDecodedEvent.EventType {
	case "PublicationMadeAvailable":
		return events.PublicationMadeAvailable{}, nil
	}
	return nil, fmt.Errorf("unmanaged publication event %v", nextDecodedEvent.EventType)
}

func (publications PublicationsHistory) Append(ctx context.Context, newEvents []application.ContextualizedEvent) error {
	userID := ctx.Value(contexts.UserIDContextKey{}).(contexts.UserID)
	publicationID := ctx.Value(contexts.PublicationIDContextKey{}).(application.PublicationID)

	documentsToInsert := []PublicationEvent{}
	for _, newEvent := range newEvents {
		documentsToInsert = append(documentsToInsert, PublicationEvent{
			UserID:        string(userID),
			PublicationID: publicationID.String(),
			EventType:     newEvent.OriginalEvent.(events.PublicationEvent).PublicationEventName(),
		})
	}
	return publications.history.Append(documentsToInsert)
}

func NewPublicationsHistory(client *DatabaseClient) application.PublicationsHistory {
	return PublicationsHistory{
		history: NewHistory[string, PublicationEvent](client, publicationsEventsCollectionName),
	}
}
