package mongodb

import (
	"context"
	"log"
	"os"
	"sarah/types/mongodb"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func GetPhoneNumberByOrgId(orgId string) []mongodb.PhoneNumber {
	coll := Client.Database(orgId).Collection(os.Getenv("MONGO_COLLECTION_PHONE_NUMBERS"))

	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var phoneNumbers []mongodb.PhoneNumber
	if err := cursor.All(context.Background(), &phoneNumbers); err != nil {
		log.Fatal(err)
	}

	return phoneNumbers
}
