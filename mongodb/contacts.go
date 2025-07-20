package mongodb

import (
	"context"
	"log"
	"os"
	"sarah/types/mongodb"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetContactByOrgId(orgId string) []mongodb.Contact {
	coll := Client.Database(orgId).Collection(os.Getenv("MONGO_COLLECTION_CONTACTS"))

	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var contacts []mongodb.Contact
	if err := cursor.All(context.Background(), &contacts); err != nil {
		log.Fatal(err)
	}

	return contacts
}
