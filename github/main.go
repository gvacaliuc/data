package main

import (
	"fmt"
	"log"
	"os"

	"net/http"

	"github.com/go-playground/webhooks/v6/github"
)

const (
	EnvGithubSecret = "GITHUB_TOKEN"
	path            = "/webhooks"
)

func main() {
	secret, ok := os.LookupEnv(EnvGithubSecret)
	if !ok {
		log.Fatalf("must set environment variable '%s'", EnvGithubSecret)
	}

	hook, _ := github.New(github.Options.Secret(secret))

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.ReleaseEvent, github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				return
			}
		}

		switch payload.(type) {
		case github.ReleasePayload:
			release := payload.(github.ReleasePayload)
			// Do whatever you want from here...
			fmt.Printf("%+v", release)

		case github.PullRequestPayload:
			pullRequest := payload.(github.PullRequestPayload)
			// Do whatever you want from here...
			fmt.Printf("%+v", pullRequest)
		}
	})
	http.ListenAndServe(":8080", nil)
}
