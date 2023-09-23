package connector

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/archiscope/archiscope-bot/internal/cache"
)

type Competition struct {
	ID              primitive.ObjectID `bson:"_id"`
	Title           string             `bson:"title"`
	Meta            string             `bson:"meta"`
	Tags            []string           `bson:"tags"`
	Register_date   string             `bson:"register_date"`
	Submission_date string             `bson:"submission_date"`
	Link            string             `bson:"link"`
	Hash            string             `bson:"hash"`
}

//	var credentials = options.Credential{
//		AuthSource: os.Getenv("MONGO_AUTH_SOURCE"),
//		Username:   os.Getenv("MONGO_AUTH_SOURCE"),
//		Password:   os.Getenv("MONGO_AUTH_SOURCE"),
//	}
var credentials = options.Credential{
	AuthSource: "admin",
	Username:   "crwl",
	Password:   "Crawly97",
}

func NewItemChan(ctx context.Context) <-chan []Competition {
	var fileCache cache.Cacher
	competitionStream := make(chan []Competition)
	fileCache, err := cache.NewFileCache("./cache.txt")
	if err != nil {
		log.Fatal(err)
	}
	collection := connect(ctx, credentials)

	go loop(ctx, fileCache, collection, competitionStream)
	return competitionStream
}

func connect(ctx context.Context, credentials options.Credential) (collection *mongo.Collection) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	clientOptions.SetAuth(credentials)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("test").Collection("test_collection")
	documentCount, _ := collection.CountDocuments(ctx, bson.D{})
	if documentCount == 0 {
		log.Fatalf("Collection %s appears to be empty", "test_collection")
	}
	return collection
}

func get_last(ctx context.Context, collection *mongo.Collection, filter bson.D) (items []Competition, err error) {
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	err = cur.All(ctx, items)
	if err != nil {
		return items, err
	}
	fmt.Print(items)
	return items, nil
}

func loop(ctx context.Context, cache cache.Cacher, collection *mongo.Collection, competitionStream chan<- []Competition) {
	for {
		lastItemHash, err := cache.GetEntry(0)
		if err != nil {
			log.Fatal(err)
		}
		filter := bson.D{{"_hash", bson.D{{"$gt", lastItemHash}}}}
		item, err := get_last(ctx, collection, filter)
		if err != nil {
			break
		}
		competitionStream <- item
		time.Sleep(2 * time.Second)
	}
}
