package management

import (
	"fmt"
	"net/http"
	"time"

	rate "github.com/beefsack/go-rate"
	"github.com/ingaged/sling"

	"github.com/illyabusigin/contentful/models"
)

const baseURL = "https://api.contentful.com"

// A Client manages communication with the Contentful Management API.
type Client struct {
	AccessToken string

	sling *sling.Sling
	rl    *rate.RateLimiter
}

////////////////////
// Initialization //
////////////////////

// NewClient creates a new Contentful API client
func NewClient(accessToken string, version string, httpClient *http.Client) *Client {
	client := &Client{
		AccessToken: accessToken,
		sling: sling.New().Client(httpClient).Base(baseURL).
			Set("Content-Type", contentTypeHeader(version)).
			Set("Authorization", authorizationHeader(accessToken)),
	}

	client.rl = rate.New(10, time.Second*1)

	return client
}

func contentTypeHeader(version string) string {
	return fmt.Sprintf("application/vnd.contentful.management.%v+json", version)
}

func authorizationHeader(accessToken string) string {
	return fmt.Sprintf("Bearer %v", accessToken)
}

func handleError(reqErr error, err *models.Error) error {
	if reqErr != nil {
		return reqErr
	}

	if err.RequestID == "" && err.Message == "" {
		return nil
	}

	return err
}

// Doer executes http requests.  It is implemented by *http.Client.  You can
// wrap *http.Client with layers of Doers to form a stack of client-side
// middleware.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}
