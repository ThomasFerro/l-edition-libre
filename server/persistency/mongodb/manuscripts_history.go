package mongodb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/events"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ManuscriptsHistory struct {
	client *DatabaseClient
}

const manuscriptsEventsCollectionName = "manuscripts_events"

type ManuscriptEvent struct {
	ManuscriptID string `bson:"manuscriptId"`
	UserID       string `bson:"userId"`
	EventType    string `bson:"eventType"`
	EventPayload string `bson:"eventPayload"`
}

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

func manuscriptToStreamKey(manuscriptEvent ManuscriptEvent) application.ManuscriptID {
	return application.MustParseManuscriptID(manuscriptEvent.ManuscriptID)
}

func (history ManuscriptsHistory) findManuscripts(query bson.D) (map[application.ManuscriptID][]application.ContextualizedEvent, error) {
	return find(history.client, query, manuscriptEventMapper, manuscriptToStreamKey)
}

func (history ManuscriptsHistory) For(manuscriptID application.ManuscriptID) ([]application.ContextualizedEvent, error) {
	query := bson.D{primitive.E{Key: "manuscriptId", Value: manuscriptID.String()}}
	results, err := history.findManuscripts(query)
	if err != nil {
		return nil, err
	}
	return results[manuscriptID], nil
}

func (history ManuscriptsHistory) ForAll() (map[application.ManuscriptID][]application.ContextualizedEvent, error) {
	return history.findManuscripts(bson.D{})
}

func (history ManuscriptsHistory) ForAllOfUser(userID application.UserID) (map[application.ManuscriptID][]application.ContextualizedEvent, error) {
	return history.findManuscripts(bson.D{primitive.E{Key: "userId", Value: userID.String()}})
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
		fmt.Printf("\n\n\n %v unmarshal %v err %v \n\n", nextDecodedEvent.EventPayload, payload, err)
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

type ManuscriptSubmittedDto struct {
	Title    string
	Author   string
	FileName string
	FileURL  string
}

func toDto(manuscriptEvent events.ManuscriptEvent) (interface{}, error) {
	switch manuscriptEvent.ManuscriptEventName() {
	case "ManuscriptReviewed":
		return nil, nil
	case "ManuscriptSubmissionCanceled":
		return nil, nil
	case "ManuscriptSubmitted":
		manuscriptSubmitted := manuscriptEvent.(events.ManuscriptSubmitted)
		return ManuscriptSubmittedDto{
			Title:    manuscriptSubmitted.Title,
			Author:   manuscriptSubmitted.Author,
			FileName: manuscriptSubmitted.FileName,
			FileURL:  manuscriptSubmitted.FileURL.String(),
		}, nil
	}
	return nil, fmt.Errorf("unmanaged manuscript event %v", manuscriptEvent.ManuscriptEventName())
}

// TODO: Mutualiser le append
// TODO: Ajouter la date
func (history ManuscriptsHistory) Append(ctx context.Context, newEvents []application.ContextualizedEvent) error {
	userID := ctx.Value(contexts.UserIDContextKey{}).(application.UserID)
	manuscriptID := ctx.Value(contexts.ManuscriptIDContextKey{}).(application.ManuscriptID)
	collection := Collection(history.client, manuscriptsEventsCollectionName)
	mongoctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	documentsToInsert := []interface{}{}
	for _, newEvent := range newEvents {
		manuscriptEvent := newEvent.OriginalEvent.(events.ManuscriptEvent)
		dto, err := toDto(manuscriptEvent)
		if err != nil {
			return err
		}
		payload, err := json.Marshal(dto)
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
	_, err := collection.InsertMany(mongoctx, documentsToInsert)
	if err != nil {
		return err
	}
	return nil
}

func NewManuscriptsHistory(client *DatabaseClient) application.ManuscriptsHistory {
	return ManuscriptsHistory{
		client,
	}
}
