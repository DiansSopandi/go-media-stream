package repository

import (
	"database/sql"

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

	var metaStr string
	if track.Metadata == "" {
		metaStr = "{}"
	} else {
		metaStr = track.Metadata
	}

	query := `INSERT INTO tracks (user_id, filename, metadata, title, artist, duration) 
			  VALUES ($1, $2, $3, $4, $5, $6) 
			  RETURNING id, user_id, filename, metadata, title, artist, duration, created_at, updated_at`
	err := tx.QueryRow(query, track.UserID, track.Filename, metaStr, track.Title, track.Artist, track.Duration).Scan(
		&track.ID,
		&track.UserID,
		&track.Filename,
		&track.Metadata,
		&track.Title,
		&track.Artist,
		&track.Duration,
		&track.CreatedAt,
		&track.UpdatedAt)

	if err != nil {
		return models.Track{}, err
	}

	return *track, nil
}
