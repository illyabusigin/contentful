package management

import (
	_ "fmt"
	"io/ioutil"
	"net/http"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestFetchAllSpacesRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	_, _, err := client.FetchAllSpaces()
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces", req.URL.String())
	assert.Equal(t, req.Method, http.MethodGet)
}

func TestFetchAllSpacesResponseSuccess(t *testing.T) {

}

func TestFetchSpaceRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	spaceID := "space123"
	_, err := client.FetchSpace(spaceID)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123", req.URL.String())
	assert.Equal(t, req.Method, http.MethodGet)
}

func TestFetchSpaceResponseSuccess(t *testing.T) {

}

func TestCreateSpaceRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	_, err := client.CreateAsset(&goodFile)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/abc123/assets", req.URL.String())
	assert.Equal(t, req.Method, http.MethodPost)

	expectedJSON := string(`{
    "fields": {
        "title": {
            "en-US": "Cat pictures"
        },
        "file": {
            "en-US": {
                "contentType": "image/png",
                "fileName": "cats.png",
                "upload": "//images.contentful.com/haian05f1d28/6rDHXkKllCOwoIiKMqgUQu/ea8f4bbba7581f21a32a3e68f56850e3/med106330_1210_bacon_pancakes_horiz.jpg"
            }
        }
    }
}`)

	requestJSON, _ := ioutil.ReadAll(req.Body)
	assert.JSONEq(t, expectedJSON, string(requestJSON))
}

func TestCreateSpaceResponseSuccess(t *testing.T) {

}

func TestUpdateSpaceRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	updated := &Space{
		Name: "Test Space, updated name",
		System: System{
			Type:      "Space",
			ID:        "space123",
			Version:   1,
			CreatedAt: &createdDate,
			UpdatedAt: &updatedDate,
		},
	}

	_, err := client.UpdateSpace(updated)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123", req.URL.String())
	assert.Equal(t, req.Method, http.MethodPut)
	assert.Equal(t, "1", req.Header.Get("X-Contentful-Version"))

	expectedJSON := string(`{
  "name":"Test Space, updated name",
  "sys":{
    "type":"Space",
    "id":"space123",
    "version":1,
    "createdAt":"2016-07-25T16:00:00Z",
    "updatedAt":"2016-07-26T16:00:00Z"
  }
}
`)

	requestJSON, _ := ioutil.ReadAll(req.Body)
	assert.JSONEq(t, expectedJSON, string(requestJSON))
}

func TestUpdateSpaceResponseSuccess(t *testing.T) {

}

func TestDeleteSpaceRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	err := client.DeleteSpace("space123")
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123", req.URL.String())
	assert.Equal(t, req.Method, http.MethodDelete)
}

func TestDeleteSpaceResponseSuccess(t *testing.T) {

}
