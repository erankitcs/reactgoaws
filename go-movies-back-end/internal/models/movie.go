package models

import (
	"time"
)

type Movie struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"release_date"`
	RunTime     int       `json:"runtime"`
	MPAARating  string    `json:"mpaa_rating"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Genres      []*Genre  `json:"genres,omitempty"`
	GenresArray []int     `json:"genres_array,omitempty"`
}

type Genre struct {
	ID        int       `json:"id"`
	Genre     string    `json:"genre"`
	Checked   bool      `json:"checked"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type MovieVideo struct {
	ID        int       `json:"id"`
	MovieID   int       `json:"movie_id"`
	VideoPath string    `json:"video_path"`
	CreatedAt time.Time `json:"created_at"`
	IsLatest  bool      `json:"is_latest"`
}

type MovieChat struct {
	ID      int `json:"id"`
	MovieID int `json:"movie_id"`
	UserID  int `json:"user_id"`
	//FirstName string    `json:"first_name"`
	//LastName  string    `json:"last_name"`
	ChatText  string    `json:"chattext"`
	CreatedAt time.Time `json:"created_at"`
}
