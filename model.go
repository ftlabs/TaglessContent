package main

type request struct {
	Id    string `json:"uuid,omitempty"`
	Token string `json:"token,omitempty`
}

type content struct {
	BodyXml string `json:"bodyXML,omitempty"`
}
