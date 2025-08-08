package sarah

import (
	"context"
	"errors"
	"log"
	"sarah/mongodb"
	mongodbTypes "sarah/types/mongodb"

	vapiApi "github.com/VapiAI/server-sdk-go"

	vapiclient "github.com/VapiAI/server-sdk-go/client"
	"github.com/joho/godotenv"
)

var VapiClient *vapiclient.Client

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

}

func CreateCall(assistantId string, assistantNumberId string, customers []mongodbTypes.Customer) (*vapiApi.CallsCreateResponse, error) {
	if err := ensureVapiClient(); err != nil {
		return nil, err
	}

	customerList := []*vapiApi.CreateCustomerDto{}
	for _, customer := range customers {
		customerList = append(customerList, &vapiApi.CreateCustomerDto{
			Number: vapiApi.String(customer.PhoneNumber),
		})
	}

	if len(customerList) == 0 {
		log.Printf("No customers/phone numbers provided")
		return nil, errors.New("no customers/phone numbers provided")
	}

	resp, err := VapiClient.Calls.Create(context.Background(), &vapiApi.CreateCallDto{
		AssistantId:   vapiApi.String(assistantId),
		PhoneNumberId: vapiApi.String(assistantNumberId),
		Customers:     customerList,
	})

	if err != nil {
		log.Printf("Error creating call: %v", err)
		return nil, err
	}

	log.Printf("Call created successfully: %+v\n", resp)

	return resp, nil
}

func GetCall(callId string) (*vapiApi.Call, error) {
	if err := ensureVapiClient(); err != nil {
		return nil, err
	}

	resp, err := VapiClient.Calls.Get(context.Background(), callId)
	if err != nil {
		log.Printf("Error getting call: %v", err)
		return nil, err
	}

	return resp, nil
}

func ListCalls(callListRequest *vapiApi.CallsListRequest) ([]*vapiApi.Call, error) {
	if err := ensureVapiClient(); err != nil {
		return nil, err
	}

	resp, err := VapiClient.Calls.List(context.Background(), callListRequest)
	if err != nil {
		log.Printf("Error listing calls: %v", err)
		return nil, err
	}

	return resp, nil
}

func GetOrganizationCalls(orgId string) ([]*vapiApi.Call, error) {
	if err := ensureVapiClient(); err != nil {
		return nil, err
	}

	assistants, err := mongodb.GetOrganizationAssistants(orgId)
	if err != nil {
		log.Printf("Error getting organization assistants: %v", err)
		return nil, err
	}

	calls := []*vapiApi.Call{}
	for _, assistant := range assistants {
		assistantCalls, err := VapiClient.Calls.List(context.Background(), &vapiApi.CallsListRequest{
			AssistantId: vapiApi.String(assistant.VapiAssistantId),
		})
		if err != nil {
			log.Printf("Error listing calls: %v", err)
			return nil, err
		}
		calls = append(calls, assistantCalls...)
	}

	return calls, nil
}
