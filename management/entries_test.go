package management

import (
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/illyabusigin/contentful/models"
	assert "github.com/stretchr/testify/require"
)

var (
	goodEntry = &Entry{
		Fields: EntryFields{
			"en-US": map[string]interface{}{"title": "blah"},
		},
		System: System{
			ID: "123",
			Space: &Link{
				LinkData: &LinkData{ID: "space123"},
			},
		},
	}
)

func TestNewEntryValidation(t *testing.T) {
	var validationTests = []struct {
		entry    NewEntry
		expected string
	}{
		{NewEntry{}, "Empty NewEntry should return an error"},
	}

	for _, test := range validationTests {
		err := test.entry.Validate()
		assert.NotNil(t, err, test.expected)
	}
}

func TestEntryValidation(t *testing.T) {
	var validationTests = []struct {
		entry    Entry
		expected string
	}{
		{Entry{}, "Empty Entry should return an error"},
		{Entry{System: System{Space: &Link{LinkData: &LinkData{ID: "space123"}}}}, "Empty Entry should return an error"},
		{Entry{System: System{ID: "123", Space: &Link{LinkData: &LinkData{ID: "space123"}}}}, "Entry with empty fields should return an error"},
	}

	for _, test := range validationTests {
		err := test.entry.Validate()
		assert.NotNil(t, err, test.expected)
	}

	goodEntry := Entry{
		Fields: EntryFields{
			"en-US": map[string]interface{}{"title": "blah"},
		},
		System: System{
			ID: "123",
			Space: &Link{
				LinkData: &LinkData{ID: "space123"},
			},
		},
	}
	err := goodEntry.Validate()
	assert.Nil(t, err, "Error should be nil since the entry is valid!")
}

func TestQueryEntriesRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	params := map[string]string{
		"include": "all",
	}
	_, _, err := client.QueryEntries("space123", params, 100, 0)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/entries?include=all&limit=100&skip=0", req.URL.String())
	assert.Equal(t, req.Method, http.MethodGet)

	// Invalid spaceID
	_, _, err = client.QueryEntries("", nil, 100, 0)
	assert.NotNil(t, err)

	// Invalid limit
	_, _, err = client.QueryEntries("space123", nil, -100, 0)
	assert.NotNil(t, err)

	// Restricting limit to 100
	_, _, err = client.QueryEntries("space123", params, 500, 0)
	req = doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/entries?include=all&limit=100&skip=0", req.URL.String())
}

func TestQueryEntriesResponseSuccess(t *testing.T) {
}

func TestFetchEntryRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	spaceID := "space123"
	entryID := "entry456"
	_, err := client.FetchEntry(spaceID, entryID)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/entries/entry456", req.URL.String())
	assert.Equal(t, req.Method, http.MethodGet)

	// Fetch with invalid spaceID and entryID
	_, err = client.FetchEntry("", "")
	assert.NotNil(t, err)
}

func TestFetchEntryResponseSuccess(t *testing.T) {

}

func TestCreateEntryRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	newEntry := &NewEntry{
		Fields: EntryFields{
			"title": map[string]interface{}{
				"en-US": "Cat-pictures",
			},
			"file": map[string]interface{}{
				"en-US": "Cat-pictures.png",
			},
		},
	}

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
					ID:       "abc123",
					LinkType: "Space",
					Type:     "Link",
				},
			},
		},
	}

	_, err := client.CreateEntry(newEntry, newContentType)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/abc123/entries", req.URL.String())
	assert.Equal(t, "TestType", req.Header.Get("X-Contentful-Content-Type"))
	assert.Equal(t, http.MethodPost, req.Method)

	expectedJSON := string(`{
    "fields": {
        "file": {
            "en-US": "Cat-pictures.png"
        },
        "title": {
            "en-US": "Cat-pictures"
        }
    }
}`)

	requestJSON, _ := ioutil.ReadAll(req.Body)
	assert.JSONEq(t, expectedJSON, string(requestJSON))

	// nil entry
	_, err = client.CreateEntry(nil, nil)
	assert.NotNil(t, err)

	// invalid content type
	_, err = client.CreateEntry(newEntry, &ContentType{})
	assert.NotNil(t, err)

	// invalid content type
	_, err = client.CreateEntry(&NewEntry{}, newContentType)
	assert.NotNil(t, err)
}

func TestCreateEntryResponseSuccess(t *testing.T) {

}

func TestUpdateEntryRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	newEntry := &Entry{
		Fields: EntryFields{
			"title": map[string]interface{}{
				"en-US": "Cat-pictures",
			},
			"file": map[string]interface{}{
				"en-US": "Cat-pictures.png",
			},
		},
		System: System{
			ID:   "entry456",
			Type: "Entry",
			Space: &Link{
				LinkData: &LinkData{
					ID:       "abc123",
					LinkType: "Space",
					Type:     "Link",
				},
			},
			ContentType: &Link{
				LinkData: &LinkData{
					ID:       "TestType",
					LinkType: "Link",
					Type:     "Link",
				},
			},
		},
	}

	_, err := client.UpdateEntry(newEntry)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/abc123/entries/entry456", req.URL.String())
	assert.Equal(t, http.MethodPut, req.Method)

	expectedJSON := string(`{
    "sys": {
        "id": "entry456",
        "type": "Entry",
        "space": {
            "sys": {
                "type": "Link",
                "linkType": "Space",
                "id": "abc123"
            }
        },
        "contentType": {
            "sys": {
                "type": "Link",
                "linkType": "Link",
                "id": "TestType"
            }
        }
    },
    "fields": {
        "file": {
            "en-US": "Cat-pictures.png"
        },
        "title": {
            "en-US": "Cat-pictures"
        }
    }
}`)

	requestJSON, _ := ioutil.ReadAll(req.Body)
	assert.JSONEq(t, expectedJSON, string(requestJSON))

	// Nil entry
	_, err = client.UpdateEntry(nil)
	assert.NotNil(t, err)

	// Invalid entry
	_, err = client.UpdateEntry(&Entry{})
	assert.NotNil(t, err)
}

func TestUpdateEntryResponseSuccess(t *testing.T) {

}

func TestDeleteEntryRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	spaceID := "space123"
	entryID := "entry456"
	err := client.DeleteEntry(entryID, spaceID)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/entries/entry456", req.URL.String())
	assert.Equal(t, req.Method, http.MethodDelete)
}

func TestDeleteEntryResponseSuccess(t *testing.T) {

}

func TestPublishEntryRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	_, err := client.PublishEntry(goodEntry)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/entries/123/published", req.URL.String())
	assert.Equal(t, req.Method, http.MethodPut)

	// Nil entry
	_, err = client.PublishEntry(nil)
	assert.NotNil(t, err)

	// Invalid entry
	_, err = client.PublishEntry(&Entry{})
	assert.NotNil(t, err)
}

func TestPublishEntryResponseSuccess(t *testing.T) {

}

func TestUnpublishEntryRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	_, err := client.UnpublishEntry(goodEntry)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/entries/123/published", req.URL.String())
	assert.Equal(t, req.Method, http.MethodDelete)

	// Nil entry
	_, err = client.UnpublishEntry(nil)
	assert.NotNil(t, err)

	// Invalid entry
	_, err = client.UnpublishEntry(&Entry{})
	assert.NotNil(t, err)
}

func TestUnpublishEntryResponseSuccess(t *testing.T) {

}

func TestArchiveEntryRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	_, err := client.ArchiveEntry(goodEntry)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/entries/123/archived", req.URL.String())
	assert.Equal(t, req.Method, http.MethodPut)

	// Nil entry
	_, err = client.ArchiveEntry(nil)
	assert.NotNil(t, err)

	// Invalid entry
	_, err = client.ArchiveEntry(&Entry{})
	assert.NotNil(t, err)
}

func TestArchiveEntryResponseSuccess(t *testing.T) {

}

func TestUnarchiveEntryRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	_, err := client.UnarchiveEntry(goodEntry)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space123/entries/123/archived", req.URL.String())
	assert.Equal(t, req.Method, http.MethodDelete)

	// Nil entry
	_, err = client.UnarchiveEntry(nil)
	assert.NotNil(t, err)

	// Invalid entry
	_, err = client.UnarchiveEntry(&Entry{})
	assert.NotNil(t, err)
}

func TestUnarchiveEntryResponseSuccess(t *testing.T) {

}
