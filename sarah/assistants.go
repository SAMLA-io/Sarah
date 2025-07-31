package sarah

import (
	"context"
	"log"
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
}

func CreateAsisstant(orgId string, assistantCreateDto vapiApi.CreateAssistantDto) (*mongo.InsertOneResult, error) {
	assistant, err := VapiClient.Assistants.Create(context.Background(), &assistantCreateDto)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result, err := mongodb.CreateAssistant(orgId, mongodbTypes.Assistant{
		Name:            *assistant.Name,
		VapiAssistantId: assistant.Id,
		Type:            "assistantCreateDto.Type",
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func UpdateAssistant(assistantId string, assistantUpdateDto vapiApi.UpdateAssistantDto) (*vapiApi.Assistant, error) {
	result, err := VapiClient.Assistants.Update(context.Background(), assistantId, &assistantUpdateDto)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return result, nil
}

func DeleteAssistant(orgId string, assistantId string) (*mongo.DeleteResult, error) {
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
	assistant, err := VapiClient.Assistants.Get(context.Background(), assistantId)
	if err != nil {
		return false
	}

	return assistant != nil && assistant.Id != ""
}
