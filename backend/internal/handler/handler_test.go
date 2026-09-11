package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/foolxdev/url-shortener/internal/model"
)

func TestValidateURLCreate(t *testing.T) {
	future := time.Now().Add(time.Hour)
	past := time.Now().Add(-time.Hour)

	tests := []struct {
		name  string
		input model.UrlCreate
		want  string
	}{
		{name: "valid URL", input: model.UrlCreate{Name: "Docs", Link: "https://example.com/docs", ExpiresAt: &future}},
		{name: "missing name", input: model.UrlCreate{Link: "https://example.com"}, want: "a link name is required"},
		{name: "name too long", input: model.UrlCreate{Name: strings.Repeat("a", 121), Link: "https://example.com"}, want: "link name must be 120 characters or fewer"},
		{name: "unsupported scheme", input: model.UrlCreate{Name: "File", Link: "ftp://example.com"}, want: "a valid http or https URL is required"},
		{name: "missing host", input: model.UrlCreate{Name: "Path", Link: "https:///path"}, want: "a valid http or https URL is required"},
		{name: "expired", input: model.UrlCreate{Name: "Old", Link: "https://example.com", ExpiresAt: &past}, want: "expiration time must be in the future"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUrlCreate(tt.input)
			if tt.want == "" && err != nil {
				t.Fatalf("validateUrlCreate returned %v", err)
			}
			if tt.want != "" && (err == nil || err.Error() != tt.want) {
				t.Fatalf("validateUrlCreate error = %v, want %q", err, tt.want)
			}
		})
	}
}

func TestCreateShortURLRejectsMalformedPayload(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/urls", strings.NewReader(`{"name":"Docs","link":"https://example.com","extra":true}`))
	rec := httptest.NewRecorder()

	CreateShortUrl(nil).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
