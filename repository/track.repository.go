package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/DiansSopandi/media_stream/db"
	"github.com/DiansSopandi/media_stream/models"
)

type TrackRepository struct {
	DB *sql.DB
	TX *sql.Tx
}

func NewTrackRepository(tx *sql.Tx) (*TrackRepository, error) {
	return &TrackRepository{
		DB: db.InitDatabase(),
		TX: tx,
	}, nil
}

func (r *TrackRepository) CreateTrack(tx *sql.Tx, track *models.Track) (models.Track, error) {
	fmt.Println("Creating track metadata in repository:", track.Metadata)
	var jsonMetadata []byte
	var err error
	if jsonMetadata, err = json.Marshal(track.Metadata); err != nil {
		return models.Track{}, err
	}
	fmt.Println("Serialized track metadata to JSON:", string(jsonMetadata))

	query := `INSERT INTO tracks (user_id, filename, metadata, title, artist, duration, is_public) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7) 
			  RETURNING id, user_id, filename, metadata, title, artist, duration, is_public, created_at, updated_at`
	err = tx.QueryRow(query, track.UserID, track.Filename, jsonMetadata, track.Title, track.Artist, track.Duration, track.IsPublic).Scan(
		&track.ID,
		&track.UserID,
		&track.Filename,
		&track.Metadata,
		&track.Title,
		&track.Artist,
		&track.Duration,
		&track.IsPublic,
		&track.CreatedAt,
		&track.UpdatedAt)

	if err != nil {
		return models.Track{}, err
	}

	return *track, nil
}
