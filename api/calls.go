package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	api "github.com/VapiAI/server-sdk-go"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}
}

func CreateCall(w http.ResponseWriter, r *http.Request) {
	if !VerifyMethod(r, []string{"POST"}) {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	client, err := GetClientFromRequest(r)
	if err != nil {
		log.Printf("Error creating client: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

	resp, err := client.Calls.Create(context.Background(), &api.CreateCallDto{
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

	client, err := GetClientFromRequest(r)
	if err != nil {
		log.Printf("Error creating client: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	callId := ExtractCallId(r)

	resp, err := client.Calls.Get(context.Background(), callId)

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

	client, err := GetClientFromRequest(r)
	if err != nil {
		log.Printf("Error creating client: %v", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	calls, err := client.Calls.List(context.Background(), &api.CallsListRequest{
		AssistantId: api.String(ExtractAssistantId(r)),
	})

	if err != nil {
		log.Printf("Error listing calls: %v", err)
		http.Error(w, "Failed to list calls", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(calls)
}
