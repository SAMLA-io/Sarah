package mongodb

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Contact struct {
	Id          bson.ObjectID          `json:"id" bson:"_id"`
	Name        string                 `json:"name" bson:"name"`
	Email       string                 `json:"email" bson:"email"`
	PhoneNumber string                 `json:"phone_number" bson:"phone_number"`
	Company     string                 `json:"company" bson:"company"`
	Position    string                 `json:"position" bson:"position"`
	Address     string                 `json:"address" bson:"address"`
	Metadata    map[string]interface{} `json:"metadata" bson:"metadata"`
}
