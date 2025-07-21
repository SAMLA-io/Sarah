package sarah

import (
	"context"
	"log"
	"sarah/mongodb"
	mongodbTypes "sarah/types/mongodb"

	"os"

	vapiApi "github.com/VapiAI/server-sdk-go"

	vapiclient "github.com/VapiAI/server-sdk-go/client"
	"github.com/joho/godotenv"
)

var VapiClient *vapiclient.Client

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using system environment variables")
	}

	VapiClient = createClient(os.Getenv("VAPI_API_KEY"))
}

func CreateCall(assistantId string, assistantNumberId string, customers []mongodbTypes.Customer) *vapiApi.CallsCreateResponse {
	customerList := []*vapiApi.CreateCustomerDto{}
	for _, customer := range customers {
		customerList = append(customerList, &vapiApi.CreateCustomerDto{
			Number: vapiApi.String(customer.PhoneNumber),
		})
	}

	if len(customerList) == 0 {
		log.Printf("No customers/phone numbers provided")
		return nil
	}

	resp, err := VapiClient.Calls.Create(context.Background(), &vapiApi.CreateCallDto{
		AssistantId:   vapiApi.String(assistantId),
		PhoneNumberId: vapiApi.String(assistantNumberId),
		Customers:     customerList,
	})

	if err != nil {
		log.Printf("Error creating call: %v", err)
		return nil
	}

	log.Printf("Call created successfully: %+v\n", resp)

	return resp
}

func GetCall(callId string) *vapiApi.Call {
	resp, err := VapiClient.Calls.Get(context.Background(), callId)
	if err != nil {
		log.Printf("Error getting call: %v", err)
		return nil
	}

	return resp
}

func ListCalls(callListRequest *vapiApi.CallsListRequest) []*vapiApi.Call {
	resp, err := VapiClient.Calls.List(context.Background(), callListRequest)
	if err != nil {
		log.Printf("Error listing calls: %v", err)
		return nil
	}

	return resp
}

func GetOrganizationCalls(orgId string) []*vapiApi.Call {
	assistants := mongodb.GetOrganizationAssistants(orgId)

	calls := []*vapiApi.Call{}
	for _, assistant := range assistants {
		assistantCalls, err := VapiClient.Calls.List(context.Background(), &vapiApi.CallsListRequest{
			AssistantId: vapiApi.String(assistant.VapiAssistantId),
		})
		if err != nil {
			log.Printf("Error listing calls: %v", err)
			return nil
		}
		calls = append(calls, assistantCalls...)
	}

	return calls
}
