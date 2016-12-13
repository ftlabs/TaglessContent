package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/url"
)

const (
	validUuid = "07d60fd8-62fe-4949-860d-3c7026c4f8c0"
	invalidUuid = "Not a uuid"
)

func TestErrorThrownWithInvalidTokenInRequest(t *testing.T) {
	assert := assert.New(t)
	Url, _ := url.Parse("http://localhost:8080")
	sts := newStripTagsService("myValidKey", *Url, "myValidToken")
	_, err, status := sts.stripTagsFromContent(request{Id:validUuid, Token:"invalidToken"})
	assert.Error(err, "Token invalidToken is invalid!")
	assert.Equal(401, status, "Unexpected status")
}

func TestErrorThrownWithEmptyTokenInRequest(t *testing.T) {
	assert := assert.New(t)
	Url, _ := url.Parse("http://localhost:8080")
	sts := newStripTagsService("myValidKey", *Url, "myValidToken")
	_, err, status := sts.stripTagsFromContent(request{Id:validUuid})
	assert.Error(err, "Token  is invalid!")
	assert.Equal(401, status, "Unexpected status")
}

func TestErrorThrownWithInvalidUuidInRequest(t *testing.T) {
	assert := assert.New(t)
	Url, _ := url.Parse("http://localhost:8080")
	sts := newStripTagsService("myValidKey", *Url, "myValidToken")
	_, err, status := sts.stripTagsFromContent(request{Id:invalidUuid, Token:"myValidToken"})
	assert.Error(err, "Not a uuid. Is not a valid uuid.")
	assert.Equal(400, status, "Unexpected status")
}

func TestErrorThrownWithEmptyUuidInRequest(t *testing.T) {
	assert := assert.New(t)
	Url, _ := url.Parse("http://localhost:8080")
	sts := newStripTagsService("myValidKey", *Url, "myValidToken")
	_, err, status := sts.stripTagsFromContent(request{Token:"myValidToken"})
	assert.Error(err, ". Is not a valid uuid.")
	assert.Equal(400, status, "Unexpected status")
}


