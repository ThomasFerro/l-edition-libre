package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoIndex struct {
	Collection string
	Indexes    []mongo.IndexModel
}

var STREAMS_INDEXES = []MongoIndex{
	{
		Collection: manuscriptsEventsCollectionName,
		Indexes: []mongo.IndexModel{
			{
				Keys: bson.M{"event.userId": 1},
			},
		},
	},
	{
		Collection: publicationsEventsCollectionName,
	},
	{
		Collection: usersEventsCollectionName,
	},
}
