package management

import (
	"fmt"
	"net/http"
	"time"

	rate "github.com/beefsack/go-rate"
	"github.com/ingaged/sling"
)

const baseURL = "https://api.contentful.com"

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

////////////////
// Base Types //
////////////////

type System struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time

	Type    string
	Version int

	Space *struct {
		*Link `json:"sys"`
	} `json:"space, omitempty"`
}

type Pagination struct {
	Total int
	Skip  int
	Limit int
}

type ContentfulError struct {
	RequestID string `json:"requestId"`
	Message   string `json:"message"`
	Sys       struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"sys"`
}

func (e *ContentfulError) Error() string {
	return e.Message
}
