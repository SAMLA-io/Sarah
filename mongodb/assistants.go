package mongodb

import (
	"context"
	"log"
	"os"
	"sarah/types/mongodb"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetOrganizationAssistants(orgId string) []mongodb.Assistant {
	coll := Client.Database(orgId).Collection(os.Getenv("MONGO_COLLECTION_ASSISTANTS"))

	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var assistants []mongodb.Assistant
	if err := cursor.All(context.Background(), &assistants); err != nil {
		log.Fatal(err)
	}

	return assistants
}
