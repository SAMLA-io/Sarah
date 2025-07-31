package mongodb

import (
	"context"
	"log"
	"os"
	"sarah/types/mongodb"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// GetContactByOrgId retrieves all contacts for a specific organization from the database.
// This function queries the contacts collection in the organization's database
// and returns all contact documents.
//
// Parameters:
//   - orgId: The organization ID to retrieve contacts for
//
// Returns:
//   - []mongodb.Contact: Array of contacts for the organization
//
// Database Operations:
//   - Database: Uses the organization ID as the database name
//   - Collection: Uses the MONGO_COLLECTION_CONTACTS environment variable
//   - Query: Retrieves all documents (no filtering)
//
// Error Handling:
//   - Logs and terminates the application if database operations fail
//
// Example Usage:
//
//	contacts := GetContactByOrgId("org_1234567890abcdef")
//	// Returns all contacts for the specified organization
func GetContactByOrgId(orgId string) ([]mongodb.Contact, error) {
	coll := Client.Database(orgId).Collection(os.Getenv("MONGO_COLLECTION_CONTACTS"))

	cursor, err := coll.Find(context.Background(), bson.M{})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var contacts []mongodb.Contact
	if err := cursor.All(context.Background(), &contacts); err != nil {
		log.Println(err)
		return nil, err
	}

	return contacts, nil
}

// CreateContact creates a new contact in the database.
// This function inserts a new contact document into the contacts collection
// for a specific organization.
//
// Parameters:
//   - orgId: The organization ID to create the contact for
//   - contact: The contact to create
func CreateContact(orgId string, contact mongodb.Contact) (*mongo.InsertOneResult, error) {
	coll := Client.Database(orgId).Collection(os.Getenv("MONGO_COLLECTION_CONTACTS"))

	result, err := coll.InsertOne(context.Background(), contact)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

// UpdateContact updates an existing contact in the database.
// This function updates an existing contact document in the contacts collection
// for a specific organization.
//
// Parameters:
//   - orgId: The organization ID to update the contact for
//   - contact: The contact to update
func UpdateContact(orgId string, contact mongodb.Contact) (*mongo.UpdateResult, error) {
	coll := Client.Database(orgId).Collection(os.Getenv("MONGO_COLLECTION_CONTACTS"))

	result, err := coll.UpdateOne(context.Background(), bson.M{"_id": contact.Id}, bson.M{"$set": contact})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

// DeleteContact deletes an existing contact from the database.
// This function removes a contact document from the contacts collection
// for a specific organization.
//
// Parameters:
//   - orgId: The organization ID to delete the contact for
//   - contactId: The object ID of the contact to delete
func DeleteContact(orgId string, contactId string) (*mongo.DeleteResult, error) {
	coll := Client.Database(orgId).Collection(os.Getenv("MONGO_COLLECTION_CONTACTS"))

	result, err := coll.DeleteOne(context.Background(), bson.M{"_id": contactId})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}
