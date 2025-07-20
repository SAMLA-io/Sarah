package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sarah/supabase"
	"sort"

	api "github.com/VapiAI/server-sdk-go"
	vapiclient "github.com/VapiAI/server-sdk-go/client"
	"github.com/joho/godotenv"
)

var VapiClient *vapiclient.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	VapiClient = createClient(os.Getenv("VAPI_API_KEY"))
}

func CreateCall(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"POST"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	assistantId := ExtractAssistantId(r)
	phoneNumbers := ExtractPhoneNumbers(r)
	assistantNumberId := ExtractAssistantNumberId(r)

	customerList := []*api.CreateCustomerDto{}
	for _, phoneNumber := range phoneNumbers {
		customerList = append(customerList, &api.CreateCustomerDto{
			Number: api.String(phoneNumber),
		})
	}

	if len(customerList) == 0 {
		http.Error(w, "No phone numbers provided", http.StatusBadRequest)
		return
	}

	resp, err := VapiClient.Calls.Create(context.Background(), &api.CreateCallDto{
		AssistantId:   api.String(assistantId),
		PhoneNumberId: api.String(assistantNumberId),
		Customers:     customerList,
	})

	if err != nil {
		log.Printf("Error creating call: %v", err)
		http.Error(w, "Failed to create call", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Call created successfully: %+v\n", resp)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func GetCall(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"GET"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	callId := ExtractCallId(r)

	resp, err := VapiClient.Calls.Get(context.Background(), callId)

	if err != nil {
		log.Printf("Error getting call: %v", err)
		http.Error(w, "Failed to get call", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Call retrieved successfully: %+v\n", resp)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func ListCalls(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"GET"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	callListRequest := ExtractCallListRequest(r)

	calls, err := VapiClient.Calls.List(context.Background(), callListRequest)

	if err != nil {
		log.Printf("Error listing calls: %v", err)
		http.Error(w, "Failed to list calls", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(calls)
}

func GetCallListByOrgId(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"GET"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	orgID := ExtractOrgId(r)

	assistants := supabase.GetAssistantByOrgId(orgID)

	calls := []*api.Call{}
	for _, assistant := range assistants {
		assistantCalls, err := VapiClient.Calls.List(context.Background(), &api.CallsListRequest{
			AssistantId: api.String(assistant.VapiAssistantId),
		})
		if err != nil {
			log.Printf("Error listing calls: %v", err)
			http.Error(w, "Failed to list calls", http.StatusInternalServerError)
			return
		}
		calls = append(calls, assistantCalls...)
	}

	sort.Slice(calls, func(i, j int) bool {
		return calls[i].CreatedAt.After(calls[j].CreatedAt)
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(calls)
}
