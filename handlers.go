package main

import (
	"encoding/json"
	"net/http"
	//"github.com/gorilla/mux"
	//"net/url"
	"io"
)

type stripTagsHandler struct {
	service stripTagsService
}

func newStripTagsHandler(service stripTagsService) stripTagsHandler {
	return stripTagsHandler{service: service}
}

func (h *stripTagsHandler) putHandler(w http.ResponseWriter, req *http.Request) {
	var body io.Reader = req.Body

	dec := json.NewDecoder(body)
	params, _ := decodeJSON(dec)

	resp, err, status := h.service.stripTagsFromContent(params)
	if err != nil {
		http.Error(w, err.Error(), status)
	}
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(resp))
}

func decodeJSON(dec *json.Decoder) (interface{}, error) {
	c := request{}
	err := dec.Decode(&c)
	return c, err
}
