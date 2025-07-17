package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

func ExtractPhoneNumbers(r *http.Request) []string {
	var phoneNumbers []string

	// Parse the JSON body
	type requestBody struct {
		PhoneNumbers []string `json:"phoneNumbers"`
	}

	var body requestBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&body)
	if err != nil {
		return []string{}
	}

	// Trim whitespace and filter out empty strings
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
