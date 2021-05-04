package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"
)

type payload struct {
	Links []string `json:"links"`
}

func contextValue(ctx context.Context, key string) (string, error) {
	if cv := ctx.Value(key); cv != nil {
		if value, ok := cv.(string); ok {
			return value, nil
		}
	}

	return "", errors.New(fmt.Sprintf("%s not set on the context", key))
}

const archiveBoxDirContextKey = "archiveBoxDir"

func archiveLinks(links []string, ctx context.Context) error {
	archiveBoxDir, err := contextValue(ctx, archiveBoxDirContextKey)
	if err != nil {
		return err
	}

	timeout := 300 * time.Second

	logger = GetLogger()

	linkArgString := ""
	for _, link := range links {
		_, err := url.Parse(link)
		if err != nil {
			logger.Printf("Bad URL: %s\n", link)
			return err
		} else {
			linkArgString = fmt.Sprintf("%s \"%s\"", linkArgString, link)
		}
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	dockerComposeFile := fmt.Sprintf("%s%cdocker-compose.yml", archiveBoxDir, os.PathSeparator)
	// yes, there's a potential security hole here but I'm not considering myself a worthy enough victim to go
	// beyond checking the URL parses. I trust ArchiveBox will do some proper validations itself when executing
	// code that actually uses the URLs
	cmd := exec.CommandContext(timeoutCtx, "docker-compose", "-f", dockerComposeFile, "run", "archivebox", "add", linkArgString)
	cmd.Stdout = logger.Writer()
	cmd.Stderr = logger.Writer()

	// Because this will be called very infrequently pooling the runs is not a priority
	if err := cmd.Run(); err != nil {
		logger.Println("An error occurred when running archivebox via docker-compose")
		logger.Println(err)
		return errors.New(fmt.Sprintf("Failed to archive one of the following links: %s", linkArgString))
	}
	logger.Printf("Archiving links: %s", linkArgString)
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	logger := GetLogger()
	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var l payload
		if err := decoder.Decode(&l); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.Printf("Links from request: %v", l.Links)

		if err := archiveLinks(l.Links, r.Context()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			if _, writeErr := w.Write([]byte(err.Error())); writeErr != nil {
				logger.Fatal(err)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	default:
		logger.Printf("%s method request attempt", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func injectContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestStart := time.Now()
		logger := GetLogger()
		archiveBoxDir, isSet := os.LookupEnv("ARCHIVE_BOX_DIR")
		if !isSet {
			logger.Fatal("ARCHIVE_BOX_DIR must be set to the directory of ArchiveBox's docker-compose.yml file")
		}
		ctx := context.WithValue(r.Context(), archiveBoxDirContextKey, archiveBoxDir)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
		duration := time.Since(requestStart).Milliseconds()
		logger.Printf("Request duration: %dms", duration)
	})
}

func main() {
	logger := GetLogger()
	http.Handle("/", injectContext(http.HandlerFunc(handler)))
	port := os.Getenv("ARCHIVE_BRIDGE_PORT")
	if len(port) == 0 {
		port = "3344"
	}
	logger.Printf("Listening on port %s\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		panic(err)
	}
}
