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

	query := `INSERT INTO tracks (user_id, filename, title, artist, duration) 
			  VALUES ($1, $2, $3, $4, $5) 
			  RETURNING id, user_id, filename, title, artist, duration, created_at, updated_at`
	err := tx.QueryRow(query, track.UserID, track.Filename, track.Title, track.Artist, track.Duration).Scan(
		&track.ID,
		&track.UserID,
		&track.Filename,
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
