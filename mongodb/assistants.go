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

// GetOrganizationAssistants retrieves all assistants for a specific organization from the database.
// This function queries the assistants collection in the organization's database
// and returns all assistant documents.
//
// Parameters:
//   - orgId: The organization ID to retrieve assistants for
//
// Returns:
//   - []mongodb.Assistant: Array of assistants for the organization
//
// Database Operations:
//   - Database: Uses the organization ID as the database name
//   - Collection: Uses the MONGO_COLLECTION_ASSISTANTS environment variable
//   - Query: Retrieves all documents (no filtering)
//
// Error Handling:
//   - Logs and terminates the application if database operations fail
//
// Example Usage:
//
//	assistants := GetOrganizationAssistants("org_1234567890abcdef")
//	// Returns all assistants for the specified organization
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
