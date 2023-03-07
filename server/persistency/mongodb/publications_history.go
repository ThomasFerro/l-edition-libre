package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PublicationsHistory struct {
	client *DatabaseClient
}

const publicationsEventsCollectionName = "publications_events"

type PublicationEvent struct {
	PublicationID string `bson:"publicationId"`
	UserID        string `bson:"userId"`
	EventType     string `bson:"eventType"`
}

func (history PublicationsHistory) For(publicationID application.PublicationID) ([]application.ContextualizedEvent, error) {
	// TODO: Passer le context ?
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := Collection(history.client, publicationsEventsCollectionName).Find(ctx, bson.D{primitive.E{Key: "publicationId", Value: publicationID.String()}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	results := []application.ContextualizedEvent{}
	for cur.Next(ctx) {
		var nextDecodedEvent PublicationEvent
		err := cur.Decode(&nextDecodedEvent)
		if err != nil {
			return nil, err
		}
		publicationEvent, err := toPublicationEvent(nextDecodedEvent)
		if err != nil {
			return nil, err
		}
		// TODO: Ne pas passer par des contextualized
		results = append(results, application.ContextualizedEvent{
			OriginalEvent: publicationEvent,
			Context: application.EventContext{
				UserID: application.MustParseUserID(nextDecodedEvent.UserID),
			},
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func toPublicationEvent(nextDecodedEvent PublicationEvent) (events.Event, error) {
	switch nextDecodedEvent.EventType {
	case "PublicationMadeAvailable":
		return events.PublicationMadeAvailable{}, nil
	}
	return nil, fmt.Errorf("unmanaged publication event %v", nextDecodedEvent.EventType)
}

func (history PublicationsHistory) Append(ctx context.Context, newEvents []application.ContextualizedEvent) error {
	userID := ctx.Value(contexts.UserIDContextKey{}).(application.UserID)
	publicationID := ctx.Value(contexts.PublicationIDContextKey{}).(application.PublicationID)
	collection := Collection(history.client, publicationsEventsCollectionName)
	mongoctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	documentsToInsert := []interface{}{}
	for _, newEvent := range newEvents {
		documentsToInsert = append(documentsToInsert, PublicationEvent{
			UserID:        userID.String(),
			PublicationID: publicationID.String(),
			EventType:     newEvent.OriginalEvent.(events.PublicationEvent).PublicationEventName(),
		})
	}
	_, err := collection.InsertMany(mongoctx, documentsToInsert)
	if err != nil {
		return err
	}
	return nil
}

func NewPublicationsHistory(client *DatabaseClient) application.PublicationsHistory {
	return PublicationsHistory{
		client,
	}
}
