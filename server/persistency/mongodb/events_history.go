package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/ThomasFerro/l-edition-libre/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PersistedEvent[StreamID comparable] interface {
	StreamID() (StreamID, error)
}

type EventsHistory[StreamID comparable, Event PersistedEvent[StreamID]] interface {
	Client() *DatabaseClient
	Append([]Event) error
	ForSingleStream(StreamID, bson.D) ([]Event, error)
	ForMultipleStreams(bson.D) (utils.OrderedMap[StreamID, []Event], error)
}

type GenericEventsHistory[StreamID comparable, Event PersistedEvent[StreamID]] struct {
	client     *DatabaseClient
	collection string
}

type InsertedDocument[StreamID comparable, Event PersistedEvent[StreamID]] struct {
	EmittedOn string   `bson:"emittedOn"`
	Event     Event    `bson:"event"`
	StreamID  StreamID `bson:"streamId"`
}

func (history GenericEventsHistory[StreamID, PersistedEvent]) Client() *DatabaseClient {
	return history.client
}
func (history GenericEventsHistory[StreamID, PersistedEvent]) Append(newEvents []PersistedEvent) error {
	collection := Collection(history.client, history.collection)
	mongoctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	documentsToInsert := []interface{}{}
	for _, newEvent := range newEvents {
		streamID, err := newEvent.StreamID()
		if err != nil {
			return err
		}
		eventToInsert := InsertedDocument[StreamID, PersistedEvent]{
			EmittedOn: time.Now().UTC().String(),
			StreamID:  streamID,
			Event:     newEvent,
		}
		documentsToInsert = append(documentsToInsert, eventToInsert)
	}
	_, err := collection.InsertMany(mongoctx, documentsToInsert)
	if err != nil {
		return err
	}
	return nil
}

func (history GenericEventsHistory[StreamID, PersistedEvent]) ForSingleStream(streamID StreamID, query bson.D) ([]PersistedEvent, error) {
	multipleStreams, err := history.ForMultipleStreams(query)
	if err != nil {
		return nil, err
	}

	for _, keyValue := range multipleStreams.Map() {
		nextStreamID := keyValue.Key
		if nextStreamID != streamID {
			return nil, fmt.Errorf("unexpected stream id %v", nextStreamID)
		}
	}

	return multipleStreams.Of(streamID), nil
}

func (history GenericEventsHistory[StreamID, PersistedEvent]) ForMultipleStreams(query bson.D) (utils.OrderedMap[StreamID, []PersistedEvent], error) {
	// TODO: Passer le context ?
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	findQuery := bson.D{}
	for _, v := range query {
		findQuery = append(findQuery, primitive.E{Key: fmt.Sprintf("event.%v", v.Key), Value: v.Value})
	}
	cur, err := Collection(history.client, history.collection).Find(ctx, findQuery)
	if err != nil {
		return utils.OrderedMap[StreamID, []PersistedEvent]{}, err
	}
	defer cur.Close(ctx)
	results := utils.OrderedMap[StreamID, []PersistedEvent]{}
	for cur.Next(ctx) {
		var nextDocument InsertedDocument[StreamID, PersistedEvent]
		err := cur.Decode(&nextDocument)
		if err != nil {
			return utils.OrderedMap[StreamID, []PersistedEvent]{}, err
		}
		streamKey := nextDocument.StreamID
		if err != nil {
			return utils.OrderedMap[StreamID, []PersistedEvent]{}, err
		}

		persistedEvents := []PersistedEvent{}
		if results.HasKey(streamKey) {
			persistedEvents = results.Of(streamKey)
		}
		results = results.Upsert(streamKey, append(persistedEvents, nextDocument.Event))
	}
	if err := cur.Err(); err != nil {
		return utils.OrderedMap[StreamID, []PersistedEvent]{}, err
	}
	return results, nil
}

func NewHistory[StreamID comparable, Event PersistedEvent[StreamID]](client *DatabaseClient, collectionName string) EventsHistory[StreamID, Event] {
	return GenericEventsHistory[StreamID, Event]{
		client:     client,
		collection: collectionName,
	}
}
