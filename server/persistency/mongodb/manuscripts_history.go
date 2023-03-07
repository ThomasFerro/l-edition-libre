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

// TODO: Mutualiser ce helper
// TODO: Ranger par ordre de date => l'ordre est important = on ne peut plus passer par une map ?
func (history ManuscriptsHistory) find(query bson.D) (map[application.ManuscriptID][]application.ContextualizedEvent, error) {
	// TODO: Passer le context ?
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := Collection(history.client, manuscriptsEventsCollectionName).Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	results := map[application.ManuscriptID][]application.ContextualizedEvent{}
	for cur.Next(ctx) {
		var nextDecodedEvent ManuscriptEvent
		err := cur.Decode(&nextDecodedEvent)
		if err != nil {
			return nil, err
		}
		manuscriptEvent, err := toManuscriptEvent(nextDecodedEvent)
		if err != nil {
			return nil, err
		}
		manuscriptID := application.MustParseManuscriptID(nextDecodedEvent.ManuscriptID)
		if _, exists := results[manuscriptID]; !exists {
			results[manuscriptID] = []application.ContextualizedEvent{}
		}
		// TODO: Ne pas passer par des contextualized
		results[manuscriptID] = append(results[manuscriptID], application.ContextualizedEvent{
			OriginalEvent: manuscriptEvent,
			Context: application.EventContext{
				UserID: application.MustParseUserID(nextDecodedEvent.UserID),
			},
			ManuscriptID: manuscriptID,
		})
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func (history ManuscriptsHistory) For(manuscriptID application.ManuscriptID) ([]application.ContextualizedEvent, error) {
	results, err := history.find(bson.D{primitive.E{Key: "manuscriptId", Value: manuscriptID.String()}})
	if err != nil {
		return nil, err
	}
	return results[manuscriptID], nil
}

func (history ManuscriptsHistory) ForAll() (map[application.ManuscriptID][]application.ContextualizedEvent, error) {
	return history.find(bson.D{})
}

func (history ManuscriptsHistory) ForAllOfUser(userID application.UserID) (map[application.ManuscriptID][]application.ContextualizedEvent, error) {
	return history.find(bson.D{primitive.E{Key: "userId", Value: userID.String()}})
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
