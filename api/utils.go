package api

import (
	"encoding/json"
	"net/http"
	"strings"
)

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
