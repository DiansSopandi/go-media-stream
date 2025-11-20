package service

import (
	"database/sql"

	"github.com/DiansSopandi/media_stream/models"
	"github.com/DiansSopandi/media_stream/repository"
)

type TrackService struct {
	Repo *repository.TrackRepository
}

func NewTrackService(trackRepo *repository.TrackRepository) *TrackService {
	return &TrackService{
		Repo: trackRepo,
	}
}

// func (s *TrackService) GetAllTracks() ([]models.Track, error) {
// 	return s.Repo.GetAllTracks()
// }

func (s *TrackService) CreateTrack(tx *sql.Tx, track *models.Track) (models.Track, error) {
	return s.Repo.CreateTrack(tx, track)
}
