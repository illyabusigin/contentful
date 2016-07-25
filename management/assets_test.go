package management

import (
	"bytes"
	_ "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"
)

type interceptor struct {
	request  *http.Request
	response *http.Response
	err      error
}

func (i *interceptor) Do(req *http.Request) (*http.Response, error) {
	i.request = req

	if i.response != nil {
		i.response.Request = req
	}

	return i.response, i.err
}

var (
	accessToken  = "access_token"
	version      = "v1"
	errIntercept = fmt.Errorf("Intercept error")

	goodURL  = "//images.contentful.com/haian05f1d28/6rDHXkKllCOwoIiKMqgUQu/ea8f4bbba7581f21a32a3e68f56850e3/med106330_1210_bacon_pancakes_horiz.jpg"
	goodFile = File{
		SpaceID: "abc123",
		Fields: FileFields{
			Title: map[string]string{"en-US": "Cat pictures"},
			File: map[string]FileData{"en-US": FileData{
				Name:     "cats.png",
				MIMEType: "image/png",
				URL:      &goodURL,
			}},
		},
	}
	goodAsset = Asset{
		System: System{
			ID:        "6rDHXkKllCOwoIiKMqgUQu",
			CreatedAt: time.Date(2016, 07, 25, 16, 00, 00, 0, time.UTC),
			UpdatedAt: time.Date(2016, 07, 25, 16, 00, 00, 0, time.UTC),

			Version: 1,
			Type:    "Asset",
			Space: &SpaceField{
				Link: &Link{
					Type:     LinkType,
					LinkType: "Space",
					ID:       "haian05f1d28",
				},
			},
		},
		Fields: AssetFields{
			File: map[string]AssetData{"en-US": {
				Name:     "med106330_1210_bacon_pancakes_horiz.jpg",
				URL:      &goodURL,
				MIMEType: "image/png",
			}},
			Title: map[string]string{"en-US": "Bacon Pancakes"},
		},
		Published: true,
		Processed: true,
	}
)

func TestFileValidationFailures(t *testing.T) {
	emptyURL := ""

	var validationTests = []struct {
		file     File
		expected string
	}{
		{File{}, "Empty file should return an error"},
		{File{SpaceID: "test_ID", Fields: FileFields{}}, "File with empty fields should return error"},
		{File{SpaceID: "test_ID", Fields: FileFields{Title: map[string]string{"en-US": "Cat pictures"}}}, "File with file fields should return error"},
		{File{SpaceID: "test_ID", Fields: FileFields{Title: map[string]string{"en-US": "Cat pictures"}, File: map[string]FileData{"en-US": FileData{Name: ""}}}}, "File with empty file name should return error"},
		{File{SpaceID: "test_ID", Fields: FileFields{Title: map[string]string{"en-US": "Cat pictures"}, File: map[string]FileData{"en-US": FileData{Name: "TBD", MIMEType: ""}}}}, "File with empty MIMEType should return error"},
		{File{SpaceID: "test_ID", Fields: FileFields{Title: map[string]string{"en-US": "Cat pictures"}, File: map[string]FileData{"en-US": FileData{Name: "TBD", MIMEType: "image/png", URL: &emptyURL}}}}, "File with empty URL should return error"},
	}

	for _, test := range validationTests {
		err := test.file.Validate()
		assert.NotNil(t, err, test.expected)
	}
}

func TestFileValidationSuccess(t *testing.T) {
	file := goodFile
	err := file.Validate()
	assert.Nil(t, err, "Error should be nil since we have a well-formed file!")
}

func TestAssetValidationFailures(t *testing.T) {
	//emptyURL := ""

	var validationTests = []struct {
		asset    Asset
		expected string
	}{
		{Asset{System: System{}}, "Empty file should return an error"},
	}

	for _, test := range validationTests {
		err := test.asset.Validate()
		assert.NotNil(t, err, test.expected)
	}
}

func TestCreateAssetRequest(t *testing.T) {
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

func TestCreateAssetResponse(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	client.sling = client.sling.New().Doer(doer)
	header := http.Header{"Content-Type": []string{"application/vnd.contentful.management.v1+json"}}

	doer.response = &http.Response{
		Status:        "HTTP/1.1 201 Created",
		StatusCode:    http.StatusCreated,
		Proto:         "HTTP/1.0",
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: -1,
		Header:        header,
		Close:         true,
		Body: ioutil.NopCloser(bytes.NewBuffer([]byte(`{
    "fields": {
        "title": {
            "en-US": "Bacon Pancakes"
        },
        "file": {
            "en-US": {
                "contentType": "image/jpg",
                "fileName": "med106330_1210_bacon_pancakes_horiz.jpg",
                "upload": "http://www.marthastewart.com/sites/files/marthastewart.com/styles/wmax-1500/public/d29/med106330_1210_bacon_pancakes/med106330_1210_bacon_pancakes_horiz.jpg"
            }
        }
    },
    "sys": {
        "id": "6rDHXkKllCOwoIiKMqgUQu",
        "type": "Asset",
        "version": 1,
        "createdAt": "2016-07-22T20:47:59.616Z",
        "createdBy": {
            "sys": {
                "type": "Link",
                "linkType": "User",
                "id": "7FuHFqjOeGmz91MmmGV5Vm"
            }
        },
        "space": {
            "sys": {
                "type": "Link",
                "linkType": "Space",
                "id": "haian05f1d28"
            }
        },
        "updatedAt": "2016-07-22T20:47:59.616Z",
        "updatedBy": {
            "sys": {
                "type": "Link",
                "linkType": "User",
                "id": "7FuHFqjOeGmz91MmmGV5Vm"
            }
        }
    }
}`))),
	}

	asset, err := client.CreateAsset(&goodFile)

	assert.Nil(t, err)
	assert.NotNil(t, asset)
	assert.Equal(t, "6rDHXkKllCOwoIiKMqgUQu", asset.System.ID)
	assert.Equal(t, "haian05f1d28", asset.Space.ID)
	assert.Equal(t, "Bacon Pancakes", asset.Fields.Title["en-US"])
	assert.Equal(t, "med106330_1210_bacon_pancakes_horiz.jpg", asset.Fields.File["en-US"].Name)
}

func TestUpdateAssetRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	updated := goodAsset
	updated.Fields.Title["en-US"] = "Swans on the lake!"

	_, err := client.UpdateAsset(&updated)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/haian05f1d28/assets/6rDHXkKllCOwoIiKMqgUQu", req.URL.String())
	assert.Equal(t, req.Method, http.MethodPut)

	expectedJSON := string(`{
    "sys": {
        "id": "6rDHXkKllCOwoIiKMqgUQu",
        "createdAt": "2016-07-25T16:00:00Z",
        "updatedAt": "2016-07-25T16:00:00Z",
        "type": "Asset",
        "version": 1,
        "space": {
            "sys": {
                "type": "Link",
                "linkType": "Space",
                "id": "haian05f1d28"
            }
        }
    },
    "fields": {
        "title": {
            "en-US": "Swans on the lake!"
        },
        "file": {
            "en-US": {
                "contentType": "image/png",
                "fileName": "med106330_1210_bacon_pancakes_horiz.jpg",
                "url": "//images.contentful.com/haian05f1d28/6rDHXkKllCOwoIiKMqgUQu/ea8f4bbba7581f21a32a3e68f56850e3/med106330_1210_bacon_pancakes_horiz.jpg"
            }
        }
    }
}`)

	requestJSON, _ := ioutil.ReadAll(req.Body)
	fmt.Println("JSON: ", string(requestJSON))
	assert.JSONEq(t, expectedJSON, string(requestJSON))

}

func TestUpdateAssetResponseSuccess(t *testing.T) {

}

func TestFetchAssetRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	spaceID := "space1234"
	assetID := "asset123"

	_, err := client.FetchAsset(spaceID, assetID)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space1234/assets/asset123", req.URL.String())
	assert.Equal(t, req.Method, http.MethodGet)
}

func TestFetchAssetResponseSuccess(t *testing.T) {

}

func TestFetchAssetsSuccessRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	spaceID := "space1234"

	// All assets
	_, _, err := client.FetchAssets(spaceID, false, 100, 0)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space1234/assets", req.URL.String())
	assert.Equal(t, req.Method, http.MethodGet)

	// Public assets
	doer.err = errIntercept
	_, _, err = client.FetchAssets(spaceID, true, 100, 0)
	req = doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/space1234/public/assets", req.URL.String())
	assert.Equal(t, req.Method, http.MethodGet)
}

func TestFetchAssetsSuccessResponse(t *testing.T) {
}

func TestProcessAssetRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	assetToProcess := goodAsset

	err := client.ProcessAsset(&assetToProcess, "en-US")
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/haian05f1d28/assets/6rDHXkKllCOwoIiKMqgUQu/files/en-US/process", req.URL.String())
	assert.Equal(t, req.Method, http.MethodPut)
}

func TestProcessAssetResponse(t *testing.T) {

}

func TestPublishAssetRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	assetToPublish := goodAsset

	_, err := client.PublishAsset(&assetToPublish)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/haian05f1d28/assets/6rDHXkKllCOwoIiKMqgUQu/published", req.URL.String())
	assert.Equal(t, req.Method, http.MethodPut)
}

func TestPublishAssetSuccessResponse(t *testing.T) {

}

func TestUnpublishAssetRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	assetToPublish := goodAsset

	_, err := client.UnpublishAsset(&assetToPublish)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/haian05f1d28/assets/6rDHXkKllCOwoIiKMqgUQu/published", req.URL.String())
	assert.Equal(t, req.Method, http.MethodDelete)
}

func TestUnpublishAssetSuccessResponse(t *testing.T) {
}

func TestArchiveAssetRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	assetToPublish := goodAsset

	_, err := client.ArchiveAsset(&assetToPublish)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/haian05f1d28/assets/6rDHXkKllCOwoIiKMqgUQu/archived", req.URL.String())
	assert.Equal(t, req.Method, http.MethodPut)
}

func TestArchiveAssetResponseSuccess(t *testing.T) {

}

func TestUnarchiveAssetRequest(t *testing.T) {
	client := NewClient(accessToken, version, nil)
	assert.NotNil(t, client, "Client should not be nil")

	// Inject request interceptor
	doer := &interceptor{}
	doer.err = errIntercept
	client.sling = client.sling.New().Doer(doer)

	assetToPublish := goodAsset

	_, err := client.UnarchiveAsset(&assetToPublish)
	req := doer.request

	assert.Equal(t, err, errIntercept)
	assert.NotNil(t, req)
	assert.Equal(t, "https://api.contentful.com/spaces/haian05f1d28/assets/6rDHXkKllCOwoIiKMqgUQu/archived", req.URL.String())
	assert.Equal(t, req.Method, http.MethodDelete)
}

func TestUnarchiveAssetResponseSuccess(t *testing.T) {

}
