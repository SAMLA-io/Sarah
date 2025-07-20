package mongodb

import "go.mongodb.org/mongo-driver/v2/bson"

type Assistant struct {
	Id              bson.ObjectID `json:"id" bson:"_id"`
	Name            string        `json:"name" bson:"name"`
	VapiAssistantId string        `json:"vapi_assistant_id" bson:"vapi_assistant_id"`
	Type            string        `json:"type" bson:"type"`
}
