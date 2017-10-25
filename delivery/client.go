package delivery

import (
	"fmt"
	"net/http"
	"time"

	rate "github.com/beefsack/go-rate"
	"github.com/ingaged/sling"
)

const baseURL = "https://cdn.contentful.com"

// PaginationSizeLimit is the sizel limit for pages
var PaginationSizeLimit = 1000

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
	type Params struct {
		AccessToken string `url:"access_token,omitempty"`
		Locale      string `url:"locale,omitempty"`
	}

	params := &Params{AccessToken: accessToken, Locale: "*"}

	client := &Client{
		AccessToken: accessToken,
		sling: sling.New().Client(httpClient).Base(baseURL).
			Set("Content-Type", contentTypeHeader(version)).
			QueryStruct(params),
	}

	client.rl = rate.New(10, time.Second*1)

	return client
}

func contentTypeHeader(version string) string {
	return fmt.Sprintf("application/vnd.contentful.delivery.%v+json", version)
}

func authorizationHeader(accessToken string) string {
	return fmt.Sprintf("Bearer %v", accessToken)
}

func handleError(reqErr error, err *ContentfulError) error {
	if reqErr != nil {
		return reqErr
	}

	if err.RequestID == "" && err.Message == "" {
		return nil
	}

	return err
}

// ContentError is the error object for errors that get returned as party of entry queries
type ContentError struct {
	Details struct {
		ID       string `json:"id"`
		Type     string `json:"type"`
		LinkType string `json:"linkType"`
	} `json:"details"`
	Sys struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"sys"`
}

func (e *ContentError) Error() string {
	return fmt.Sprintf("Error: %v", e.Sys.ID)
}

// ContentfulError represnts the error object that is returned when something
// goes wrong with a Contentful API request. This struct conforms to the `error`
// interface.
type ContentfulError struct {
	RequestID string `json:"requestId"`
	Message   string `json:"message"`
	Sys       struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"sys"`
}

func (e *ContentfulError) Error() string {
	return fmt.Sprintf("%v, %v, %v", e.Message, e.RequestID, e.Sys)
}

// Doer executes http requests.  It is implemented by *http.Client.  You can
// wrap *http.Client with layers of Doers to form a stack of client-side
// middleware.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}
