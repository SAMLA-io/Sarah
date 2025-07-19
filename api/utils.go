package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	vapiclient "github.com/VapiAI/server-sdk-go/client"
	"github.com/VapiAI/server-sdk-go/option"
)

// ClientFactory creates a new VAPI client with the given API key
func createClient(apiKey string) *vapiclient.Client {
	return vapiclient.NewClient(
		option.WithToken(apiKey),
		option.WithHTTPClient(
			&http.Client{
				Timeout: 30 * time.Second,
			}),
	)
}

// getClientFromRequest creates a client using the API key from the request
func GetClientFromRequest(r *http.Request) (*vapiclient.Client, error) {
	apiKey := ExtractAuthHeader(r)
	if apiKey == "" {
		// Fallback to environment variable for backward compatibility
		apiKey = os.Getenv("VAPI_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("no API key provided in Authorization header or VAPI_API_KEY environment variable")
		}
	}
	return createClient(apiKey), nil
}

func ExtractAuthHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	return strings.TrimPrefix(authHeader, "Bearer ")
}

func ExtractPhoneNumbers(r *http.Request) []string {
	var phoneNumbers []string

	type requestBody struct {
		PhoneNumbers []string `json:"phoneNumbers"`
	}

	var body requestBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		return []string{}
	}

	for _, phone := range body.PhoneNumbers {
		phoneNumbers = append(phoneNumbers, strings.TrimSpace(phone))
	}

	return phoneNumbers
}

func ExtractAssistantId(r *http.Request) string {
	assistantId := r.URL.Query().Get("assistantId")
	return strings.TrimSpace(assistantId)
}

func ExtractAssistantNumberId(r *http.Request) string {
	assistantNumberId := r.URL.Query().Get("assistantNumberId")
	return strings.TrimSpace(assistantNumberId)
}

func VerifyMethod(r *http.Request, allowedMethods []string) bool {
	for _, method := range allowedMethods {
		if r.Method == strings.ToUpper(method) {
			return true
		}
	}
	return false
}

func ExtractCallId(r *http.Request) string {
	path := strings.TrimPrefix(r.URL.Path, "/calls/")
	return strings.TrimSpace(path)
}
