package mongodb

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Campaign struct {
	Id     bson.ObjectID `json:"id" bson:"_id"`
	VapiId string        `json:"vapi_id" bson:"vapi_id"`
	Type   string        `json:"type" bson:"type"`
}
