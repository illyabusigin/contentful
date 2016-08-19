package management

import (
	"fmt"
	"net/http"
	"time"

	rate "github.com/beefsack/go-rate"
	"github.com/ingaged/sling"
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

func handleError(reqErr error, err *ContentfulError) error {
	if reqErr != nil {
		return reqErr
	}

	return err
}

////////////////
// Base Types //
////////////////

// System contains system managed metadata. The exact metadata available depends
// on the type of the resource but at minimum System.Type property is defined.
// Note that none of the sys fields are editable and only the sys.id field can
// be specified in the creation of an item (as long as it is not a Space).
type System struct {
	ID        string     `json:"id"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`

	Type    string `json:"type,omitempty"`
	Version int    `json:"version,omitempty"`

	Space       *Link `json:"space,omitempty"`
	ContentType *Link `json:"contentType,omitempty"`

	FirstPublished   *time.Time `json:"firstPublishedAt,omitempty"`
	PublishedAt      *time.Time `json:"publishedAt,omitempty"`
	PublishedVersion int        `json:"publishedVersion,omitempty"`
	ArchivedAt       *time.Time `json:"archivedAt,omitempty"`
}

// Link represents a link to another Contentful object
type Link struct {
	*LinkData `json:"sys"`
}

// LinkData contains the link information
type LinkData struct {
	Type     string `json:"type"`
	LinkType string `json:"linkType"`
	ID       string `json:"id"`
}

// Pagination represents all paginated data and is returned for collection
// related client methods.
type Pagination struct {
	Total int
	Skip  int
	Limit int
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
