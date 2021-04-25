package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"log"
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

	for _, link := range links {
		// todo print errors from docker compose
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
		dockerComposeFile := fmt.Sprintf("%s%cdocker-compose.yml", archiveBoxDir, os.PathSeparator)
		_, err := url.Parse(link)
		if err != nil {
			log.Printf("Bad URL: %s\n", link)
		} else {
			// yes, there's a potential security hole here but I'm not considering myself a worthy enough victim to go
			// beyond checking the URL parses. I trust ArchiveBox will do some proper validations itself when executing
			// code that actually uses the URLs
			cmd := exec.CommandContext(timeoutCtx, "docker-compose", "-f", dockerComposeFile, "run", "archivebox", "add", link)

			// Because this will be called very infrequently pooling the runs is not a priority
			if err := cmd.Run(); err != nil {
				log.Println(err)
				cancel()
				return errors.New(fmt.Sprintf("Failed to archive %s", link))
			} else {
				cancel()
				log.Printf("Archiving %s", link)
			}
		}
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(r.Body)
		var l payload
		if err := decoder.Decode(&l); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Links from request: %v", l.Links)

		if err := archiveLinks(l.Links, r.Context()); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			if _, writeErr := w.Write([]byte(err.Error())); writeErr != nil {
				log.Fatal(err)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func injectContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		archiveBoxDir, isSet := os.LookupEnv("ARCHIVE_BOX_DIR")
		if !isSet {
			log.Fatal("ARCHIVE_BOX_DIR must be set to the directory of ArchiveBox's docker-compose.yml file")
		}
		ctx := context.WithValue(r.Context(), archiveBoxDirContextKey, archiveBoxDir)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func main() {
	http.Handle("/", injectContext(http.HandlerFunc(handler)))
	port := os.Getenv("ARCHIVE_BRIDGE_PORT")
	if len(port) == 0 {
		port = "3344"
	}
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		panic(err)
	}
}
