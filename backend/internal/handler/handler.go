package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/foolxdev/url-shortener/internal/model"
	"github.com/go-chi/chi/v5"
)

func CreateShortURL(app *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var urlCreate model.UrlCreate

		r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&urlCreate); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if err := decoder.Decode(&struct{}{}); err != io.EOF {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		urlCreate.Name = strings.TrimSpace(urlCreate.Name)
		urlCreate.Link = strings.TrimSpace(urlCreate.Link)
		if err := validateURLCreate(urlCreate); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		url, err := app.Insert(urlCreate)
		if err != nil {
			log.Printf("create short URL: %v", err)
			http.Error(w, "Could not create the short URL", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(url); err != nil {
			log.Printf("encode short URL response: %v", err)
		}
	}
}

func RedirectURL(app *Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")

		link, err := app.GetHash(hash)
		if err != nil {
			if errors.Is(err, errExpired) {
				http.Error(w, "This short URL has expired", http.StatusGone)
				return
			}
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}

		log.Printf("HASH=%q LINK=%q", hash, link)
		http.Redirect(w, r, link, http.StatusFound)
	}
}

func validateURLCreate(input model.UrlCreate) error {
	if input.Name == "" {
		return errors.New("a link name is required")
	}
	if len(input.Name) > 120 {
		return errors.New("link name must be 120 characters or fewer")
	}
	parsed, err := url.ParseRequestURI(input.Link)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return errors.New("a valid http or https URL is required")
	}
	if input.ExpiresAt != nil && !input.ExpiresAt.After(time.Now()) {
		return errors.New("expiration time must be in the future")
	}
	return nil
}
