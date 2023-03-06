package mongodb

import (
	"context"
	"time"

	"github.com/ThomasFerro/l-edition-libre/configuration"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slog"
)

func GetClient() (*mongo.Client, error) {
	slog.Info("connecting to mongodb")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(configuration.GetConfiguration(configuration.MONGO_CONNECTION_STRING)))
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
	return client, err
}

// TODO: Migration as code pour créer les collections (exposer sur un /init)
var database = configuration.GetConfiguration(configuration.MONGO_DATABASE_NAME)

func Collection(client *mongo.Client, collection string) *mongo.Collection {
	return client.Database(database).Collection(collection)
}

/*
Manuscripts:

StreamID => ManuscriptID
Creator => UserID

Users:

StreamID => UserID
Creator => UserID

Publications:

StreamID => PublicationID


Une seule table Events
*/
/*
TODO: Heathcheck via
Calling Connect does not block for server discovery. If you wish to know if a MongoDB server has been found and connected to, use the Ping method:

ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
err = client.Ping(ctx, readpref.Primary())
*/
