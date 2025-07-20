package api

import (
	"encoding/json"
	"io"
	"net/http"
	"sarah/types/mongodb"
	"strings"
	"time"

	api "github.com/VapiAI/server-sdk-go"
	vapiclient "github.com/VapiAI/server-sdk-go/client"
	"github.com/VapiAI/server-sdk-go/option"
)

func createClient(apiKey string) *vapiclient.Client {
	return vapiclient.NewClient(
		option.WithToken(apiKey),
		option.WithHTTPClient(
			&http.Client{
				Timeout: 30 * time.Second,
			}),
	)
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
	callId := r.URL.Query().Get("callId")
	return strings.TrimSpace(callId)
}

func ExtractCallListRequest(r *http.Request) *api.CallsListRequest {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil
	}

	var bodyMap map[string]json.RawMessage
	if err := json.Unmarshal(body, &bodyMap); err != nil {
		return nil
	}

	raw, ok := bodyMap["callListRequest"]
	if !ok || raw == nil {
		return nil
	}

	var callListRequest api.CallsListRequest
	if err := json.Unmarshal(raw, &callListRequest); err != nil {
		return nil
	}
	return &callListRequest
}

func ExtractCampaignId(r *http.Request) string {
	path := strings.TrimPrefix(r.URL.Path, "/campaigns/campaign/")
	campaignId := strings.TrimSpace(path)
	return strings.TrimSpace(campaignId)
}

func ExtractOrgId(r *http.Request) string {
	orgId := r.URL.Query().Get("orgId")
	return strings.TrimSpace(orgId)
}

func ExtractCampaignCreateDto(r *http.Request) *mongodb.CampaignCreateDto {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil
	}

	var requestBody struct {
		CampaignCreateRequest mongodb.CampaignCreateDto `json:"campaignCreateRequest"`
	}

	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		return nil
	}

	return &requestBody.CampaignCreateRequest
}
