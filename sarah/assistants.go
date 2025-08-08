package sarah

import (
	"context"
	"errors"
	"log"
	"os"
	"sarah/mongodb"
	mongodbTypes "sarah/types/mongodb"

	vapiApi "github.com/VapiAI/server-sdk-go"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	apiKey := os.Getenv("VAPI_API_KEY")
	if apiKey == "" {
		log.Printf("ERROR: VAPI_API_KEY environment variable is not set")
		return
	}

	VapiClient = createClient(apiKey)
	if VapiClient == nil {
		log.Printf("ERROR: Failed to create VapiClient")
	}
}

// ensureVapiClient checks if VapiClient is properly initialized
func ensureVapiClient() error {
	if VapiClient == nil {
		return errors.New("VapiClient is not initialized - check VAPI_API_KEY environment variable")
	}
	return nil
}

func CreateAsisstant(orgId string, assistantCreateDto vapiApi.CreateAssistantDto) (*mongo.InsertOneResult, error) {
	if err := ensureVapiClient(); err != nil {
		return nil, err
	}

	assistant, err := VapiClient.Assistants.Create(context.Background(), &assistantCreateDto)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result, err := mongodb.CreateAssistant(orgId, mongodbTypes.Assistant{
		Name:            *assistant.Name,
		VapiAssistantId: assistant.Id,
		Type:            "placeholder type",
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func UpdateAssistant(assistantId string, assistantUpdateDto vapiApi.UpdateAssistantDto) (*vapiApi.Assistant, error) {
	if err := ensureVapiClient(); err != nil {
		return nil, err
	}

	result, err := VapiClient.Assistants.Update(context.Background(), assistantId, &assistantUpdateDto)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func DeleteAssistant(orgId string, assistantId string) (*mongo.DeleteResult, error) {
	if err := ensureVapiClient(); err != nil {
		return nil, err
	}

	_, err := VapiClient.Assistants.Delete(context.Background(), assistantId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result, err := mongodb.DeleteAssistant(orgId, assistantId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func ExistsAssistant(assistantId string) bool {
	if err := ensureVapiClient(); err != nil {
		log.Printf("ExistsAssistant: VapiClient not initialized: %v", err)
		return false
	}

	log.Printf("ExistsAssistant: Checking existence of assistant with ID: %s", assistantId)
	assistant, err := VapiClient.Assistants.Get(context.Background(), assistantId)
	if err != nil {
		log.Printf("ExistsAssistant: Error fetching assistant with ID %s: %v", assistantId, err)
		return false
	}

	exists := assistant != nil && assistant.Id != ""
	log.Printf("ExistsAssistant: Assistant with ID %s exists: %v", assistantId, exists)
	return exists
}
