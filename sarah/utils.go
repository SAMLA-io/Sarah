package sarah

import (
	"net/http"
	"time"

	vapiclient "github.com/VapiAI/server-sdk-go/client"
	"github.com/VapiAI/server-sdk-go/option"
)

// createClient initializes and returns a new VapiAI client with the provided API key.
// The client is configured with a 30-second timeout for HTTP requests.
//
// Parameters:
//   - apiKey: The VapiAI API key for authentication
//
// Returns:
//   - *vapiclient.Client: Configured VapiAI client instance
func createClient(apiKey string) *vapiclient.Client {
	return vapiclient.NewClient(
		option.WithToken(apiKey),
		option.WithHTTPClient(
			&http.Client{
				Timeout: 30 * time.Second,
			}),
	)
}
