package models

import (
	"time"
)

type Track struct {
	ID          int        `json:"id" db:"id"`
	UserID      int        `json:"user_id" db:"user_id"`
	Filename    string     `json:"filename" db:"filename"`
	Metadata    string     `json:"metadata" db:"metadata"`
	Title       string     `json:"title" db:"title"`
	Artist      string     `json:"artist" db:"artist"`
	Album       string     `json:"album" db:"album"`
	Genre       string     `json:"genre" db:"genre"`
	Description string     `json:"description" db:"description"`
	Duration    int        `json:"duration" db:"duration"`
	IsPublic    bool       `json:"is_public" db:"is_public"`
	PlayCount   int        `json:"play_count" db:"play_count"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type TrackResponse struct {
	ID        int        `json:"id" db:"id"`
	AlbumID   int        `json:"album_id" db:"album_id"`
	Title     string     `json:"title" db:"title"`
	Duration  int        `json:"duration" db:"duration"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}

type TrackRequest struct {
	AlbumID   int    `json:"album_id" db:"album_id"`
	Title     string `json:"title" db:"title"`
	Duration  int    `json:"duration" db:"duration"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `json:"deleted_at,omitempty" db:"deleted_at"`
}
