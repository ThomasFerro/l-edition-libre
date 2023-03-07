package mongodb

import (
	"context"
	"time"

	"github.com/ThomasFerro/l-edition-libre/configuration"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slog"
)

type DatabaseClient struct {
	MongoClient  *mongo.Client
	databaseName string
}

// TODO: Exposer et utiliser une méthode pour fermer le client

func GetClient(databaseName string) (*DatabaseClient, error) {
	// TODO: Un seul client ? Un pool ? Qui s'occupe de les lancer / les couper ?
	slog.Info("connecting to mongodb", "databaseName", databaseName)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(configuration.GetConfiguration(configuration.MONGO_CONNECTION_STRING)))
	if err != nil {
		slog.Error("cannot connect to mongodb", err)
		return nil, err
	}
	// defer func() {
	// 	if err = client.Disconnect(ctx); err != nil {
	// 		slog.Error("cannot disconnect the mongodb client", err)
	// 		panic(err)
	// 	}
	// }()
	return &DatabaseClient{
		mongoClient,
		databaseName,
	}, err
}

// TODO: Migration as code pour créer les collections (exposer sur un /init)
func Collection(client *DatabaseClient, collection string) *mongo.Collection {
	return client.MongoClient.Database(client.databaseName).Collection(collection)
}

/*
TODO: Heathcheck via
Calling Connect does not block for server discovery. If you wish to know if a MongoDB server has been found and connected to, use the Ping method:

ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
err = client.Ping(ctx, readpref.Primary())
*/
