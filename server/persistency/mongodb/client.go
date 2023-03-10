package mongodb

import (
	"context"
	"time"

	"github.com/ThomasFerro/l-edition-libre/configuration"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/exp/slog"
)

type DatabaseClient struct {
	MongoClient  *mongo.Client
	databaseName string
}

func (d DatabaseClient) Close() {
	if err := d.MongoClient.Disconnect(context.Background()); err != nil {
		slog.Error("mongo disconnection error", err)
	}
}

func (d DatabaseClient) InitDatabase() error {
	for _, streamIndexex := range STREAMS_INDEXES {
		indexesToAdd := []mongo.IndexModel{
			{
				Keys: bson.M{"streamId": 1},
			},
		}
		indexesToAdd = append(indexesToAdd, streamIndexex.Indexes...)

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		_, err := Collection(&d, streamIndexex.Collection).Indexes().CreateMany(ctx, indexesToAdd)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d DatabaseClient) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return d.MongoClient.Ping(ctx, readpref.Primary())
}

func GetClient(databaseName string) (*DatabaseClient, error) {
	slog.Info("connecting to mongodb", "databaseName", databaseName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(configuration.GetConfiguration(configuration.MONGO_CONNECTION_STRING)))
	if err != nil {
		slog.Error("cannot connect to mongodb", err)
		return nil, err
	}
	return &DatabaseClient{
		mongoClient,
		databaseName,
	}, err
}

func Collection(client *DatabaseClient, collection string) *mongo.Collection {
	return client.MongoClient.Database(client.databaseName).Collection(collection)
}
