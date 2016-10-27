package management

import (
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/illyabusigin/contentful/models"
	assert "github.com/stretchr/testify/require"
)

func TestSpaceValidationFailures(t *testing.T) {
	var validationTests = []struct {
		space    Space
		expected string
	}{
		{Space{}, "Empty space should return an error"},
		{Space{Name: ""}, "Space with empty name should return error"},
	}

	for _, test := range validationTests {
		err := test.space.Validate()
		assert.NotNil(t, err, test.expected)
	}

	goodSpace := Space{Name: "test"}
	err := goodSpace.Validate()
	assert.Nil(t, err, "Error should be nil since the space is valid!")
}

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

	newSpace := &Space{
		Name: "Test Space",
		System: System{
			Type:      "Space",
			ID:        "space123",
			Version:   1,
			CreatedAt: &createdDate,
			UpdatedAt: &updatedDate,
		},
	}

	_, err := client.CreateSpace(newSpace)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces", req.URL.String())
	assert.Equal(t, req.Method, http.MethodPost)

	expectedJSON := string(`{
    "sys": {
        "id": "space123",
        "createdAt": "2016-07-25T16:00:00Z",
        "updatedAt": "2016-07-26T16:00:00Z",
        "type": "Space",
        "version": 1
    },
    "name": "Test Space"
}`)

	requestJSON, _ := ioutil.ReadAll(req.Body)
	assert.JSONEq(t, expectedJSON, string(requestJSON))

	// nil space
	_, err = client.CreateSpace(nil)
	assert.NotNil(t, err)

	// invalid space
	_, err = client.CreateSpace(&Space{})
	assert.NotNil(t, err)
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

	_, err = client.UpdateSpace(nil)
	assert.NotNil(t, err)
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
