package management

import (
	"testing"

	assert "github.com/stretchr/testify/require"
	//expect "gopkg.in/gavv/httpexpect.v1"
)

func TestNewClient(t *testing.T) {
	accessToken := "access_token"
	version := "v1"
	client := NewClient(accessToken, version, nil)

	assert.NotNil(t, client, "Client should not be nil")
	assert.NotNil(t, client.rl)
	assert.Equal(t, client.AccessToken, accessToken)
}

func TestContentTypeHeader(t *testing.T) {
	header := contentTypeHeader("v1")
	assert.Equal(t, "application/vnd.contentful.management.v1+json", header)
}

func TestAuthorizationHeader(t *testing.T) {
	header := authorizationHeader("access_token")
	assert.Equal(t, "Bearer access_token", header)
}
