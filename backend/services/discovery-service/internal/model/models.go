package model

import "time"

type Profile struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	Bio           string    `json:"bio"`
	Location      string    `json:"location"`
	ProfileType   string    `json:"profile_type"`
	AvatarURL     *string   `json:"avatar_url,omitempty"`
	Position      string    `json:"position,omitempty"`
	PreferredFoot string    `json:"preferred_foot,omitempty"`
	Height        int       `json:"height,omitempty"`
	Weight        int       `json:"weight,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type Video struct {
	ID          string    `json:"id"`
	ProfileID   string    `json:"profile_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	FileName    string    `json:"file_name"`
	BlobURL     string    `json:"blob_url"`
	Status      string    `json:"status"`
	ViewCount   int       `json:"view_count"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProfileFilters struct {
	ProfileType string `form:"profile_type"`
	Position    string `form:"position"`
	Location    string `form:"location"`
}

type SearchResponse struct {
	Results interface{} `json:"results"`
	Total   int         `json:"total"`
	Limit   int         `json:"limit"`
	Offset  int         `json:"offset"`
}