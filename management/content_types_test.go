package management

import (
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/illyabusigin/contentful/models"
	assert "github.com/stretchr/testify/require"
)

func TestContentTypeValidationFailures(t *testing.T) {
	var validationTests = []struct {
		locale   ContentType
		expected string
	}{
		{ContentType{}, "Empty content type should return an error"},
		{ContentType{Name: ""}, "Content type with empty name should return error"},
		{ContentType{Name: "Test Type", System: System{Space: &Link{LinkData: &LinkData{ID: "space123"}}}}, "Content type with empty ID should return error"},
		{ContentType{Name: "Test Type", System: System{ID: "TestType", Space: &Link{LinkData: &LinkData{ID: ""}}}}, "Content type with empty space ID should return error"},
	}

	for _, test := range validationTests {
		err := test.locale.Validate()
		assert.NotNil(t, err, test.expected)
	}
}

func TestContentTypeValidationSucces(t *testing.T) {
	file := goodFile
	err := file.Validate()
	assert.Nil(t, err, "Error should be nil since we have a well-formed file!")
}

func TestFetchAllContentTypesRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	// Unpublished
	spaceID := "space123"
	_, _, err := client.FetchContentTypes(spaceID, false, 500, 0)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/content_types?limit=100&skip=0", req.URL.String())
	assert.Equal(t, req.Method, http.MethodGet)

	// Published
	_, _, err = client.FetchContentTypes(spaceID, true, 100, 0)
	req = doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/public/content_types?limit=100&skip=0", req.URL.String())
	assert.Equal(t, req.Method, http.MethodGet)

	// Invalid Space identifier
	_, _, err = client.FetchContentTypes("", true, 100, 0)
	req = doer.request

	assert.Equal(t, err.Error(), "FetchContentTypes failed. Space identifier is not valid!")

	// Invalid limit
	_, _, err = client.FetchContentTypes(spaceID, true, -100, 0)
	req = doer.request

	assert.Equal(t, err.Error(), "FetchContentTypes failed. Limit must be greater than 0")
}

func TestFetchAllContentTypesResponseSuccess(t *testing.T) {

}

func TestFetchContentTypeRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	// Happy path
	spaceID := "space123"
	contentTypeID := "type456"
	_, err := client.FetchContentType(spaceID, contentTypeID)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/content_types/type456", req.URL.String())
	assert.Equal(t, req.Method, http.MethodGet)

	// Invalid space/content type ID
	_, err = client.FetchContentType("", contentTypeID)
	req = doer.request

	assert.Equal(t, err.Error(), "FetchContentType failed. Invalid spaceID or contentTypeID.")
	assert.NotNil(t, req)
}

func TestFetchContentTypeResponseSuccess(t *testing.T) {

}

func TestCreateContentTypeRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	newContentType := &ContentType{
		Name:         "Test Type",
		Description:  "Test type description",
		DisplayField: "title",
		Fields: []Field{
			Field{},
		},
		System: System{
			ID: "TestType",
			Space: &Link{
				LinkData: &LinkData{
					ID:       "space123",
					LinkType: "Space",
					Type:     "Link",
				},
			},
		},
	}

	_, err := client.CreateContentType(newContentType)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/content_types/TestType", req.URL.String())
	assert.Equal(t, req.Method, http.MethodPut)

	expectedJSON := string(`{
    "sys": {
        "id": "TestType",
        "space": {
            "sys": {
                "type": "Link",
                "linkType": "Space",
                "id": "space123"
            }
        }
    },
    "name": "Test Type",
    "description": "Test type description",
    "displayField": "title",
    "fields": [
        {
            "ID": "",
            "Name": "",
            "Type": "",
            "Localized": false,
            "Required": false,
            "Validations": null,
            "Disabled": false,
            "Omitted": false
        }
    ]
}`)

	requestJSON, _ := ioutil.ReadAll(req.Body)
	assert.JSONEq(t, expectedJSON, string(requestJSON))

	// Try passing nil as content type
	_, err = client.CreateContentType(nil)
	assert.Equal(t, err.Error(), "CreateContentType failed. Type argument was nil!")

	// Pass invalid content type
	_, err = client.CreateContentType(&ContentType{})
	assert.NotNil(t, err, "Error should not be nil since content type will fail validation!")

}

func TestCreateContentTypeResponseSuccess(t *testing.T) {

}

func TestUpdateContentTypeRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	updated := &ContentType{
		Name:         "Test Type, updated",
		Description:  "Test type description",
		DisplayField: "title",
		Fields: []Field{
			Field{},
		},
		System: System{
			ID:      "TestType",
			Version: 1,
			Space: &Link{
				LinkData: &LinkData{
					ID:       "space123",
					LinkType: "Space",
					Type:     "Link",
				},
			},
		},
	}

	_, err := client.UpdateContentType(updated)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/content_types/TestType", req.URL.String())
	assert.Equal(t, req.Method, http.MethodPut)
	assert.Equal(t, "1", req.Header.Get("X-Contentful-Version"))

	expectedJSON := string(`{
    "sys": {
        "id": "TestType",
        "version": 1,
        "space": {
            "sys": {
                "type": "Link",
                "linkType": "Space",
                "id": "space123"
            }
        }
    },
    "name": "Test Type, updated",
    "description": "Test type description",
    "displayField": "title",
    "fields": [
        {
            "ID": "",
            "Name": "",
            "Type": "",
            "Localized": false,
            "Required": false,
            "Validations": null,
            "Disabled": false,
            "Omitted": false
        }
    ]
}`)

	requestJSON, _ := ioutil.ReadAll(req.Body)
	assert.JSONEq(t, expectedJSON, string(requestJSON))
}

func TestUpdateContentTypeResponseSuccess(t *testing.T) {

}

func TestDeleteContentTypeRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	spaceID := "space123"
	contentTypeID := "type456"
	err := client.DeleteContentType(spaceID, contentTypeID)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/content_types/type456", req.URL.String())
	assert.Equal(t, req.Method, http.MethodDelete)
}

func TestActivateContentTypeRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	contentType := &ContentType{
		Name:         "Test Type",
		Description:  "Test type description",
		DisplayField: "title",
		Fields: []Field{
			Field{},
		},
		System: System{
			ID:      "TestType",
			Version: 1,
			Space: &Link{
				LinkData: &LinkData{
					ID:       "space123",
					LinkType: "Space",
					Type:     "Link",
				},
			},
		},
	}

	_, err := client.ActivateContentType(contentType)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/content_types/TestType/published", req.URL.String())
	assert.Equal(t, req.Method, http.MethodPut)

	// Pass nil content type
	_, err = client.ActivateContentType(nil)
	assert.NotNil(t, err, "Error should not be nil since content type will fail validation!")

	// Pass invalid content type
	_, err = client.ActivateContentType(&ContentType{})
	assert.NotNil(t, err, "Error should not be nil since content type will fail validation!")
}

func TestActivateContentTypeResponseSuccess(t *testing.T) {

}

func TestDeactivateContentTypeRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	contentType := &ContentType{
		Name:         "Test Type",
		Description:  "Test type description",
		DisplayField: "title",
		Fields: []Field{
			Field{},
		},
		System: System{
			ID:      "TestType",
			Version: 1,
			Space: &Link{
				LinkData: &LinkData{
					ID:       "space123",
					LinkType: "Space",
					Type:     "Link",
				},
			},
		},
	}

	_, err := client.DeactivateContentType(contentType)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/content_types/TestType/published", req.URL.String())
	assert.Equal(t, req.Method, http.MethodDelete)

	// Pass nil content type
	_, err = client.DeactivateContentType(nil)
	assert.NotNil(t, err, "Error should not be nil since content type will fail validation!")

	// Pass invalid content type
	_, err = client.DeactivateContentType(&ContentType{})
	assert.NotNil(t, err, "Error should not be nil since content type will fail validation!")
}

func TestDeactivateContentTypeResponseSuccess(t *testing.T) {

}
