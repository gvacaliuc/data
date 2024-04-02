package main

import (
	"log"
	"os"

	"net/http"

	"github.com/go-playground/webhooks/v6/github"
	"go.uber.org/zap"
)

const (
	EnvGithubSecret = "GITHUB_TOKEN"
	path            = "/webhooks"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	secret, ok := os.LookupEnv(EnvGithubSecret)
	if !ok {
		logger.Fatal("must set environment variable", zap.String("variable", EnvGithubSecret))
	}

	hook, _ := github.New(github.Options.Secret(secret))

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				return
			}
		}

		switch payload.(type) {
		case github.PullRequestPayload:
			pullRequest := payload.(github.PullRequestPayload)
			logger.Info("pr event", zap.Any("pullRequest", pullRequest))
		}
	})
	http.ListenAndServe(":8080", nil)
}
