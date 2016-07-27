package management

import (
	"io/ioutil"
	"net/http"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestLocaleValidationFailures(t *testing.T) {

	var validationTests = []struct {
		locale   Locale
		expected string
	}{
		{Locale{}, "Empty locale should return an error"},
		{Locale{Name: ""}, "Locale with empty name should return error"},
		{Locale{Name: "Test", Code: ""}, "Locale with empty code should return error"},
		{Locale{Name: "English", Code: "en-US", System: System{Space: &SpaceField{Link: &Link{ID: ""}}}}, "Locale with empty space ID should return error"},
	}

	for _, test := range validationTests {
		err := test.locale.Validate()
		assert.NotNil(t, err, test.expected)
	}
}

func TestLocaleValidationSucces(t *testing.T) {
	file := goodFile
	err := file.Validate()
	assert.Nil(t, err, "Error should be nil since we have a well-formed file!")
}

func TestFetchAllLocalesRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	_, _, err := client.FetchAllLocales("space123")
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/locales", req.URL.String())
	assert.Equal(t, req.Method, http.MethodGet)
}

func TestFetchAllLocalesResponseSuccess(t *testing.T) {

}

func TestFetchLocaleRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	spaceID := "space123"
	localeID := "locale456"
	_, err := client.FetchLocale(spaceID, localeID)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/locales/locale456", req.URL.String())
	assert.Equal(t, req.Method, http.MethodGet)
}

func TestFetchLocaleResponseSuccess(t *testing.T) {

}

func TestCreateLocaleRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	newLocale := Locale{
		Name: "German (Germany)",
		Code: "de-DE",
		System: System{
			ID: "123",
			Space: &SpaceField{
				Link: &Link{
					ID:       "space123",
					LinkType: "Space",
					Type:     "Link",
				},
			},
		},
	}

	_, err := client.CreateLocale(&newLocale)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/locales", req.URL.String())
	assert.Equal(t, req.Method, http.MethodPost)

	expectedJSON := string(`{
    "sys": {
        "id": "123",
        "space": {
            "sys": {
                "type": "Link",
                "linkType": "Space",
                "id": "space123"
            }
        }
    },
    "name": "German (Germany)",
    "code": "de-DE",
    "default": false,
    "optional": false
}`)

	requestJSON, _ := ioutil.ReadAll(req.Body)
	assert.JSONEq(t, expectedJSON, string(requestJSON))
}

func TestCreateLocaleResponseSuccess(t *testing.T) {

}

func TestUpdateLocaleRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	updated := &Locale{
		Name:     "German (Germany), updated",
		Code:     "de-DE",
		Fallback: "en-US",
		Default:  false,
		Optional: true,

		System: System{
			Type:      "Locale",
			ID:        "locale123",
			Version:   1,
			CreatedAt: &createdDate,
			UpdatedAt: &updatedDate,
			Space: &SpaceField{
				Link: &Link{
					ID:       "space123",
					LinkType: "Space",
					Type:     "Link",
				},
			},
		},
	}

	_, err := client.UpdateLocale(updated)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/locales/locale123", req.URL.String())
	assert.Equal(t, req.Method, http.MethodPut)
	assert.Equal(t, "1", req.Header.Get("X-Contentful-Version"))

	expectedJSON := string(`{
    "name": "German (Germany), updated",
    "code": "de-DE",
    "fallbackCode": "en-US",
    "default": false,
    "optional": true,
    "sys": {
        "type": "Locale",
        "id": "locale123",
        "version": 1,
        "space": {
            "sys": {
                "type": "Link",
                "linkType": "Space",
                "id": "space123"
            }
        },
        "createdAt": "2016-07-25T16:00:00Z",
        "updatedAt": "2016-07-26T16:00:00Z"
    }
}`)

	requestJSON, _ := ioutil.ReadAll(req.Body)
	assert.JSONEq(t, expectedJSON, string(requestJSON))
}

func TestUpdateLocaleResponseSuccess(t *testing.T) {

}

func TestDeleteLocaleRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	err := client.DeleteLocale("space123", "locale123")
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/locales/locale123", req.URL.String())
	assert.Equal(t, req.Method, http.MethodDelete)
}
