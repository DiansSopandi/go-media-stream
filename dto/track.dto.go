package dto

import "encoding/json"

type TrackCreateRequest struct {
	Filename string `json:"filename" validate:"required,max=50" example:"track.mp3"`
	// Metadata map[string]interface{} `json:"metadata" example:"{\"key\":\"value\"}"`
	Metadata    json.RawMessage `json:"metadata" example:"{\"key\":\"value\"}"`
	Title       string          `json:"title" validate:"required,max=100" example:"song title"`
	Artist      string          `json:"artist" validate:"required,max=100" example:"artist"`
	Album       string          `json:"album" validate:"required,max=100" example:"Tik Tok"`
	Genre       string          `json:"genre" validate:"required,max=100" example:"Tik Tok"`
	Duration    int             `json:"duration" validate:"required,min=0" example:"240"`
	Description string          `json:"description" validate:"max=500" example:"This is a sample track description."`
	IsPublic    bool            `json:"is_public" example:"true"`
}

type TrackCreateResponse struct {
	ID       string `json:"id" example:"1"`
	Filename string `json:"filename" example:"track.mp3"`
	Title    string `json:"title" example:"song title"`
	Artist   string `json:"artist" example:"artist"`
	Album    string `json:"album" example:"Tik Tok"`
	Genre    string `json:"genre" example:"Tik Tok"`
	Duration int    `json:"duration" example:"240"`
	Created  string `json:"created" example:"2023-01-01T00:00:00Z"`
	Updated  string `json:"updated" example:"2023-01-01T00:00:00Z"`
}
