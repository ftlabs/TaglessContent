package main

import (
	"regexp"
	"net/http"
	"net/url"
	"errors"
	"encoding/json"
	"github.com/kennygrant/sanitize"
	"fmt"
)

var UUIDRegexp = regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")

type stripTagsService interface {
	stripTagsFromContent(thing interface{}) (string, error, int)
}

type stripTagsServiceImpl struct {
	apiKey string
	contentAddr url.URL
	token string
}

func newStripTagsService(apiKey string, contentAddr url.URL, token string) stripTagsService {
	return &stripTagsServiceImpl{apiKey: apiKey, contentAddr: contentAddr, token: token}
}

func (sc *stripTagsServiceImpl) stripTagsFromContent(thing interface{}) (string, error, int) {
	requestParams := thing.(request)

	if sc.token != requestParams.Token {
		return "", errors.New("Token " + requestParams.Token + " is invalid!"), http.StatusUnauthorized
	}
	if !UUIDRegexp.MatchString(requestParams.Id) {
		return "", errors.New("Id:" + requestParams.Id + ". Is not a valid uuid."), http.StatusBadRequest
	}
	contentWithTags, err, status := getContent(sc.contentAddr, sc.apiKey, requestParams.Id)
	if err != nil {
		return "", err, status
	}
	contentNoTags := sanitize.HTML(contentWithTags)
	return contentNoTags, nil, status
}

func getContent(contentAddr url.URL, apiKey string, uuid string) (string, error, int) {
	contentUrl := contentAddr.String() + uuid + "?apiKey=" + apiKey
	resp, err := http.Get(contentUrl)
	if err != nil {
		return "", errors.New("Could not retrieve content from " + contentUrl), resp.StatusCode
	}
	fmt.Printf("Status is %v", resp.StatusCode)
	defer resp.Body.Close()
	contentStruct := new(content)
	json.NewDecoder(resp.Body).Decode(contentStruct)
	if contentStruct.BodyXml == "" {
		return contentStruct.BodyXml, errors.New("Content with id " + uuid + " not found."), resp.StatusCode
	}
	return contentStruct.BodyXml, nil, resp.StatusCode
}


