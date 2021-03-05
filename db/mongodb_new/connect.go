package mongodb_new

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Infrastructure struct {
	m      sync.Mutex
	mongos map[string]*mongo.Client
}

func NewInfrastructure() *Infrastructure {
	return &Infrastructure{
		mongos: map[string]*mongo.Client{},
	}
}

func (inf *Infrastructure) ConnectMongo(ctx context.Context, uri, user, pass string) (client *mongo.Client, err error) {
	inf.m.Lock()
	defer inf.m.Unlock()
	client, ok := inf.mongos[uri]
	if ok {
		return
	}
	var optionsClient = options.Client().ApplyURI(uri).SetAuth(options.Credential{
		Username: user,
		Password: pass,
	})
	client, err = mongo.NewClient(optionsClient)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		return
	}
	inf.mongos[uri] = client
	return client, err
}

func (inf *Infrastructure) GetDatabase(client *mongo.Client, dbName string) *mongo.Database {
	return client.Database(dbName)
}
