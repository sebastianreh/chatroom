package mongodb

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/sebastianreh/chatroom/internal/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const NoResultsOnFind = "mongo: no documents in result"
const timeOut = 10 * time.Second

type MongoDBier interface {
	Collection(name string) *mongo.Collection
	CleanCollectionByIds(ctx context.Context, collection string, ids ...primitive.ObjectID)
	PrepareData(ctx context.Context, collection string, documents ...interface{}) []primitive.ObjectID
	ClearCollection(ctx context.Context, collection string)
	PrepareCollectionWithTTL(ctx context.Context, collection string)
}

type MongoDB struct {
	client mongo.Client
	config config.Config
}

func NewMongoDB(c config.Config) MongoDBier {
	opts := options.Client()
	opts.ApplyURI(c.MongoDB.URI)

	client, err := mongo.NewClient(opts)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Print(err)
	}
	err = client.Ping(ctx, readpref.PrimaryPreferred())
	if err != nil {
		log.Print(err)
	}

	return &MongoDB{client: *client, config: c}
}

func (d *MongoDB) Collection(name string) *mongo.Collection {
	return d.client.Database(d.config.MongoDB.Database).Collection(name)
}

func (d *MongoDB) CleanCollectionByIds(ctx context.Context, collection string, ids ...primitive.ObjectID) {
	for i := range ids {
		d.Collection(collection).DeleteOne(ctx, bson.M{"_id": ids[i]})
	}
}

func (d *MongoDB) PrepareData(ctx context.Context, collection string, documents ...interface{}) []primitive.ObjectID {
	result, _ := d.Collection(collection).InsertMany(ctx, documents)
	return d.getObjectIdsFromInterfaceArray(result.InsertedIDs)
}

func (d *MongoDB) ClearCollection(ctx context.Context, collection string) {
	d.Collection(collection).DeleteMany(ctx, bson.D{})
}

func (d *MongoDB) PrepareCollectionWithTTL(ctx context.Context, collection string) {
	d.Collection(collection).Drop(ctx)
	index := d.getExpirationIndex()
	_, err := d.Collection(collection).Indexes().CreateOne(context.Background(), index)
	if err != nil {
		log.Print(err)
	}
	d.Collection(collection).Indexes().List(ctx)
}

func (d *MongoDB) getExpirationIndex() mongo.IndexModel {
	return mongo.IndexModel{
		Keys:    bson.D{{Key: "expires", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0)}
}

func (d *MongoDB) getObjectIdsFromInterfaceArray(interfaceArray []interface{}) []primitive.ObjectID {
	objectIDS := make([]primitive.ObjectID, 0)
	var idBytes []byte
	var id primitive.ObjectID

	for _, operatorMap := range interfaceArray {
		idBytes, _ = json.Marshal(operatorMap)
		json.Unmarshal(idBytes, &id)
		objectIDS = append(objectIDS, id)
	}
	return objectIDS
}
