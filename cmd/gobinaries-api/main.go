package main

import (
	"context"
	"net/http"
	"os"

	googlestorage "cloud.google.com/go/storage"
	"github.com/apex/httplog"
	"github.com/apex/log"
	"github.com/apex/log/handlers/apexlogs"
	"github.com/apex/log/handlers/logfmt"
	"github.com/apex/log/handlers/multi"
	"github.com/google/go-github/v28/github"
	"github.com/tj/go/env"
	"golang.org/x/oauth2"

	"github.com/tj/gobinaries/resolver"
	"github.com/tj/gobinaries/server"
	"github.com/tj/gobinaries/storage"
)

// main
func main() {
	// logs
	handler := &apexlogs.Handler{
		URL:       env.Get("APEX_LOGS_URL"),
		ProjectID: env.Get("APEX_LOGS_PROJECT_ID"),
	}

	if os.Getenv("APEX_LOGS_DISABLE") == "" {
		log.SetHandler(multi.New(handler, logfmt.Default))
	}

	// context
	ctx := context.Background()

	// github client
	gh := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: env.Get("GITHUB_TOKEN"),
		},
	)

	// storage client
	gs, err := googlestorage.NewClient(ctx)
	if err != nil {
		log.Fatalf("error creating storage client: %s", err)
	}

	// server
	addr := ":" + env.GetDefault("PORT", "3000")
	s := &server.Server{
		Static: "static",
		URL:    env.GetDefault("URL", "http://127.0.0.1"+addr),
		Resolver: &resolver.GitHub{
			Client: github.NewClient(oauth2.NewClient(ctx, gh)),
		},
		Storage: &storage.Google{
			Client: gs,
			Bucket: "gobinaries",
			Prefix: "production",
		},
	}

	// add request level logging
	h := flusher(httplog.New(s), handler)

	// listen
	log.WithField("addr", addr).Info("starting server")
	err = http.ListenAndServe(addr, h)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
}

// Flusher interface.
type Flusher interface {
	Flush() error
}

// flusher returns an HTTP handler which flushes after each request.
func flusher(h http.Handler, f Flusher) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)

		err := f.Flush()
		if err != nil {
			log.WithError(err).Error("error flushing logs")
			return
		}
	})
}
