package model

import (
	"time"
)

type Url struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Link      string     `json:"link"`
	Hash      string     `json:"hash"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type UrlCreate struct {
	Name      string     `json:"name"`
	Link      string     `json:"link"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type UrlGet struct {
	Name      string     `json:"name"`
	Link      string     `json:"link"`
	Hash      string     `json:"hash"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

func NewURL(
	name string,
	link string,
	expiresAt *time.Time,
) Url {
	return Url{
		Name:      name,
		Link:      link,
		ExpiresAt: expiresAt,
	}
}
