package mongodb

import (
	"context"
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UsersHistory struct {
	client  *DatabaseClient
	history EventsHistory[string, UserEvent]
}

const usersEventsCollectionName = "users_events"

type UserEvent struct {
	UserId    string `bson:"userId"`
	EventType string `bson:"eventType"`
}

func (u UserEvent) StreamID() (string, error) {
	return u.UserId, nil
}

func (users UsersHistory) For(userID application.UserID) ([]application.ContextualizedEvent, error) {
	events, err := users.history.ForSingleStream(userID.String(), bson.D{primitive.E{Key: "userId", Value: userID.String()}})
	if err != nil {
		return nil, err
	}

	contextualizedEvents := []application.ContextualizedEvent{}
	for _, nextEvent := range events {
		userEvent, err := toUserEvent(nextEvent)
		if err != nil {
			return nil, err
		}
		// TODO: Ne pas passer par des contextualized
		contextualizedEvents = append(contextualizedEvents, application.ContextualizedEvent{
			OriginalEvent: userEvent,
			Context: application.EventContext{
				UserID: userID,
			},
		})
	}
	return contextualizedEvents, nil
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

	documentsToInsert := []UserEvent{}
	for _, newEvent := range newEvents {
		documentsToInsert = append(documentsToInsert, UserEvent{
			UserId:    userID.String(),
			EventType: newEvent.OriginalEvent.(events.UserEvent).UserEventName(),
		})
	}
	return history.history.Append(documentsToInsert)
}

func NewUsersHistory(client *DatabaseClient) application.UsersHistory {
	return UsersHistory{
		history: NewHistory[string, UserEvent](client, usersEventsCollectionName),
	}
}
