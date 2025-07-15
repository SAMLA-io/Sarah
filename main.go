package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	vapiclient "github.com/VapiAI/server-sdk-go/client"
	"github.com/VapiAI/server-sdk-go/option"
	"github.com/joho/godotenv"
)

func main() {
	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create client with API key and custom timeout
	client := vapiclient.NewClient(
		option.WithToken(os.Getenv("VAPI_API_KEY")),
		option.WithHTTPClient(
			&http.Client{
				Timeout: 5 * time.Second,
			}),
	)

	fmt.Println(client)
}
