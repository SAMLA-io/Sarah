package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	api "github.com/VapiAI/server-sdk-go"
	vapiclient "github.com/VapiAI/server-sdk-go/client"
	"github.com/VapiAI/server-sdk-go/option"
	"github.com/joho/godotenv"
)

var client *vapiclient.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client = vapiclient.NewClient(
		option.WithToken(os.Getenv("VAPI_API_KEY")),
		option.WithHTTPClient(
			&http.Client{
				Timeout: 30 * time.Second,
			}),
	)
}

func CreateCall(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
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
		w.WriteHeader(http.StatusBadRequest)
	}

	resp, err := client.Calls.Create(context.Background(), &api.CreateCallDto{
		AssistantId:   api.String(assistantId),
		PhoneNumberId: api.String(assistantNumberId),
		Customers:     customerList,
	})

	if err != nil {
		log.Fatalf("Error creating call: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Printf("Call created successfully: %+v\n", resp)
	w.WriteHeader(http.StatusCreated)
}
