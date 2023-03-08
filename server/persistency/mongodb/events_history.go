package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type PersistedEvent[StreamID comparable] interface {
	StreamID() (StreamID, error)
}

type EventsHistory[StreamID comparable, Event PersistedEvent[StreamID]] interface {
	Client() *DatabaseClient
	Append([]Event) error
	ForSingleStream(StreamID, bson.D) ([]Event, error)
	ForMultipleStreams(bson.D) (map[StreamID][]Event, error)
}

type GenericEventsHistory[StreamID comparable, Event PersistedEvent[StreamID]] struct {
	client     *DatabaseClient
	collection string
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
		documentsToInsert = append(documentsToInsert, newEvent)
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

	for nextStreamID := range multipleStreams {
		if nextStreamID != streamID {
			return nil, fmt.Errorf("unexpected stream id %v", nextStreamID)
		}
	}

	return multipleStreams[streamID], nil
}

func (history GenericEventsHistory[StreamID, PersistedEvent]) ForMultipleStreams(query bson.D) (map[StreamID][]PersistedEvent, error) {
	// TODO: Passer le context ?
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cur, err := Collection(history.client, manuscriptsEventsCollectionName).Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	results := map[StreamID][]PersistedEvent{}
	for cur.Next(ctx) {
		var nextDecodedEvent PersistedEvent
		err := cur.Decode(&nextDecodedEvent)
		if err != nil {
			return nil, err
		}
		streamKey, err := nextDecodedEvent.StreamID()
		if err != nil {
			return nil, err
		}
		if _, exists := results[streamKey]; !exists {
			results[streamKey] = []PersistedEvent{}
		}
		results[streamKey] = append(results[streamKey], nextDecodedEvent)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

func NewHistory[StreamID comparable, Event PersistedEvent[StreamID]](client *DatabaseClient, collectionName string) EventsHistory[StreamID, Event] {
	return GenericEventsHistory[StreamID, Event]{
		client:     client,
		collection: collectionName,
	}
}
