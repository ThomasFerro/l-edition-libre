package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/events"
	"github.com/ThomasFerro/l-edition-libre/persistency/mongodb/dtos"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ManuscriptsHistory struct {
	history EventsHistory[string, ManuscriptEvent]
}

const manuscriptsEventsCollectionName = "manuscripts_events"

type ManuscriptEvent struct {
	ManuscriptID string `bson:"manuscriptId"`
	UserID       string `bson:"userId"`
	EventType    string `bson:"eventType"`
	EventPayload string `bson:"eventPayload"`
}

func (m ManuscriptEvent) StreamID() (string, error) {
	return m.ManuscriptID, nil
}

// TODO: Ne pas passer par des contextualized ?
func manuscriptEventMapper(manuscriptEvent ManuscriptEvent) (application.ContextualizedEvent, error) {
	manuscriptID, err := application.ParseManuscriptID(manuscriptEvent.ManuscriptID)
	if err != nil {
		return application.ContextualizedEvent{}, err
	}
	userID, err := application.ParseUserID(manuscriptEvent.UserID)
	if err != nil {
		return application.ContextualizedEvent{}, err
	}
	originalEvent, err := toManuscriptEvent(manuscriptEvent)
	if err != nil {
		return application.ContextualizedEvent{}, err
	}
	return application.ContextualizedEvent{
		OriginalEvent: originalEvent,
		Context: application.EventContext{
			UserID: userID,
		},
		ManuscriptID: manuscriptID,
	}, nil
}

func (manuscripts ManuscriptsHistory) findManuscripts(query bson.D) (map[application.ManuscriptID][]application.ContextualizedEvent, error) {
	events, err := manuscripts.history.ForMultipleStreams(query)
	if err != nil {
		return nil, err
	}

	contextualizedEvents := map[application.ManuscriptID][]application.ContextualizedEvent{}

	for _, keyValue := range events.Map() {
		rawManuscriptID := keyValue.Key
		events := keyValue.Value
		manuscriptID := application.MustParseManuscriptID(rawManuscriptID)
		contextualizedEvents[manuscriptID] = []application.ContextualizedEvent{}
		for _, nextEvent := range events {
			manuscriptEvent, err := manuscriptEventMapper(nextEvent)
			if err != nil {
				return nil, err
			}
			contextualizedEvents[manuscriptID] = append(contextualizedEvents[manuscriptID], manuscriptEvent)
		}
	}

	return contextualizedEvents, nil
}

func (manuscripts ManuscriptsHistory) For(manuscriptID application.ManuscriptID) ([]application.ContextualizedEvent, error) {
	query := bson.D{primitive.E{Key: "manuscriptId", Value: manuscriptID.String()}}
	results, err := manuscripts.findManuscripts(query)
	if err != nil {
		return nil, err
	}
	return results[manuscriptID], nil
}

func (manuscripts ManuscriptsHistory) ForAll() (map[application.ManuscriptID][]application.ContextualizedEvent, error) {
	return manuscripts.findManuscripts(bson.D{})
}

func (manuscripts ManuscriptsHistory) ForAllOfUser(userID application.UserID) (map[application.ManuscriptID][]application.ContextualizedEvent, error) {
	return manuscripts.findManuscripts(bson.D{primitive.E{Key: "userId", Value: userID.String()}})
}

func toManuscriptEvent(nextDecodedEvent ManuscriptEvent) (events.Event, error) {
	switch nextDecodedEvent.EventType {
	case "ManuscriptReviewed":
		return events.ManuscriptReviewed{}, nil
	case "ManuscriptSubmissionCanceled":
		return events.ManuscriptSubmissionCanceled{}, nil
	case "ManuscriptSubmitted":
		var payload map[string]string
		err := json.Unmarshal([]byte(nextDecodedEvent.EventPayload), &payload)
		if err != nil {
			return nil, err
		}
		url, err := url.Parse(payload["FileURL"])
		if err != nil {
			return nil, err
		}
		return events.ManuscriptSubmitted{
			Title:    payload["Title"],
			Author:   payload["Author"],
			FileName: payload["FileName"],
			FileURL:  *url,
		}, nil
	}
	return nil, fmt.Errorf("unmanaged manuscript event %v", nextDecodedEvent.EventType)
}

func (manuscripts ManuscriptsHistory) Append(ctx context.Context, newEvents []application.ContextualizedEvent) error {
	userID := ctx.Value(contexts.UserIDContextKey{}).(application.UserID)
	manuscriptID := ctx.Value(contexts.ManuscriptIDContextKey{}).(application.ManuscriptID)

	documentsToInsert := []ManuscriptEvent{}
	for _, newEvent := range newEvents {
		manuscriptEvent := newEvent.OriginalEvent.(events.ManuscriptEvent)
		payload, err := dtos.ToPayload(manuscriptEvent)
		if err != nil {
			return err
		}
		documentsToInsert = append(documentsToInsert, ManuscriptEvent{
			UserID:       userID.String(),
			ManuscriptID: manuscriptID.String(),
			EventType:    manuscriptEvent.ManuscriptEventName(),
			EventPayload: string(payload),
		})
	}
	return manuscripts.history.Append(documentsToInsert)
}

func NewManuscriptsHistory(client *DatabaseClient) application.ManuscriptsHistory {
	return ManuscriptsHistory{
		history: NewHistory[string, ManuscriptEvent](client, manuscriptsEventsCollectionName),
	}
}
