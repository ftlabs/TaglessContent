package main

import (
	"github.com/jawher/mow.cli"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"net/url"
	"net/http"
	"os"
)

func main() {
	app := cli.App("text-to-speech", "A RESTful API for retreiving Capi V2 content with no xml tags")
	apiKey := app.String(cli.StringOpt{
		Name:   "apiKey",
		Value:  "",
		Desc:   "Api Key for Capi V2 auth",
		EnvVar: "API_KEY",
	})
	contentAddr := app.String(cli.StringOpt{
		Name:   "contentAddr",
		Value:  "",
		Desc:   "Address to get content from Capi V2",
		EnvVar: "CONTENT_ADDR",
	})
	userToken := app.String(cli.StringOpt{
		Name:   "userToken",
		Value:  "",
		Desc:   "Token for accessing app",
		EnvVar: "TOKEN",
	})
	port := app.String(cli.StringOpt{
		Name:   "port",
		Value:  "8080",
		Desc:   "Port to listen on",
		EnvVar: "PORT",
	})

	app.Action = func() {
		contentUrl, err := url.Parse(*contentAddr)
		if err != nil {
			log.Fatalf("Invalid content URL: %v (%v)", *contentAddr, err)
		}

		s := newStripTagsService(*apiKey, *contentUrl, *userToken)
		h := newStripTagsHandler(s)

		m := mux.NewRouter()
		http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, m))
		m.HandleFunc("/strip", h.putHandler).Methods("PUT")

		log.Infof("Listening on [%v]", *port)
		http.ListenAndServe(":" + *port, nil)
	}
	app.Run(os.Args)

}