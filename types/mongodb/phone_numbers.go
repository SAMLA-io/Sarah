package mongodb

import "go.mongodb.org/mongo-driver/v2/bson"

type PhoneNumber struct {
	Id            bson.ObjectID `json:"id" bson:"_id"`
	Name          string        `json:"name" bson:"name"`
	PhoneNumberId string        `json:"phone_number_id" bson:"phone_number_id"`
	PhoneNumber   string        `json:"phone_number" bson:"phone_number"`
}
