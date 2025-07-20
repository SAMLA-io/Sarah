package mongodb

import (
	"context"
	"log"
	"os"
	"sarah/types/mongodb"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var Client *mongo.Client

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(os.Getenv("MONGO_URI")).SetServerAPIOptions(serverAPI)

	var err error
	Client, err = mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	if err := Client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	log.Println("Pinged deployment. Successfully connected to MongoDB!")
}

func GetCampaignByOrgId(orgId string) []mongodb.Campaign {
	coll := Client.Database(orgId).Collection(os.Getenv("MONGO_COLLECTION_CAMPAIGNS"))

	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var campaigns []mongodb.Campaign
	if err := cursor.All(context.Background(), &campaigns); err != nil {
		log.Fatal(err)
	}

	return campaigns
}

func CreateCampaign(orgId string, campaign mongodb.Campaign) *mongo.InsertOneResult {
	coll := Client.Database(orgId).Collection(os.Getenv("MONGO_COLLECTION_CAMPAIGNS"))

	result, err := coll.InsertOne(context.Background(), campaign)
	if err != nil {
		log.Fatal(err)
	}

	return result
}
