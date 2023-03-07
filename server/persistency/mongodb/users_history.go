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

type UsersHistory struct {
	client *DatabaseClient
}

const usersEventsCollectionName = "users_events"

type UserEvent struct {
	UserId    string `bson:"userId"`
	EventType string `bson:"eventType"`
}

func (history UsersHistory) For(userID application.UserID) ([]application.ContextualizedEvent, error) {
	// TODO: Passer le context ?
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := Collection(history.client, usersEventsCollectionName).Find(ctx, bson.D{primitive.E{Key: "userId", Value: userID.String()}})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	results := []application.ContextualizedEvent{}
	for cur.Next(ctx) {
		var nextDecodedEvent UserEvent
		err := cur.Decode(&nextDecodedEvent)
		if err != nil {
			return nil, err
		}
		userEvent, err := toUserEvent(nextDecodedEvent)
		if err != nil {
			return nil, err
		}
		// TODO: Ne pas passer par des contextualized
		results = append(results, application.ContextualizedEvent{
			OriginalEvent: userEvent,
			Context: application.EventContext{
				UserID: userID,
			},
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func toUserEvent(nextDecodedEvent UserEvent) (events.Event, error) {
	switch nextDecodedEvent.EventType {
	case "AccountCreated":
		return events.AccountCreated{}, nil
	case "UserPromotedToEditor":
		return events.UserPromotedToEditor{}, nil
	}
	return nil, fmt.Errorf("unmanaged user event %v", nextDecodedEvent.EventType)
}

func (history UsersHistory) Append(ctx context.Context, newEvents []application.ContextualizedEvent) error {
	userID := ctx.Value(contexts.UserIDContextKey{}).(application.UserID)
	collection := Collection(history.client, usersEventsCollectionName)
	mongoctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	documentsToInsert := []interface{}{}
	for _, newEvent := range newEvents {
		documentsToInsert = append(documentsToInsert, UserEvent{
			UserId:    userID.String(),
			EventType: newEvent.OriginalEvent.(events.UserEvent).UserEventName(),
		})
	}
	_, err := collection.InsertMany(mongoctx, documentsToInsert)
	if err != nil {
		return err
	}
	return nil
}

func NewUsersHistory(client *DatabaseClient) application.UsersHistory {
	return UsersHistory{
		client,
	}
}
