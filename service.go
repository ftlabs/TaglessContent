package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kennygrant/sanitize"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var UUIDRegexp = regexp.MustCompile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
var openPullQuote string = "<pull-quote>"
var closePullQuote string = "</pull-quote>"

type stripTagsService interface {
	stripTagsFromContent(thing interface{}) (string, error, int)
}

type stripTagsServiceImpl struct {
	apiKey      string
	contentAddr url.URL
	token       string
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
	if strings.Contains(contentWithTags, openPullQuote) {
		contentWithTags = stripPullQuotes(contentWithTags)
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

func stripPullQuotes(content string) string {
	count := strings.Count(content, openPullQuote)
	for strings.Contains(content, openPullQuote) == true {
		left := strings.SplitN(content, openPullQuote, 2)
		right := strings.SplitN(left[1], closePullQuote, 2)

		//strips first <pull-quote> text </pull-quote> it finds
		content = left[0] + right[1]
		count--
	}
	return content
}
