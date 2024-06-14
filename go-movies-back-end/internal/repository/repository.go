package repository

import (
	"backend/internal/models"
	"database/sql"
)

type DatabaseRepo interface {
	Connection() *sql.DB
	//AllMovies() ([]*models.Movie, error)
	AllMovies(genre ...int) ([]*models.Movie, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	OneMovie(id int) (*models.Movie, error)
	OneMovieForEdit(id int) (*models.Movie, []*models.Genre, error)
	AllGenres() ([]*models.Genre, error)
	InsertMovie(movie models.Movie) (int, error)
	UpdateMovie(movie models.Movie) error
	UpdateMovieGenres(id int, genreIDs []int) error
	DeleteMovie(id int) error
	// Video Management
	InsertMovieVideo(movieVideo models.MovieVideo) (*models.MovieVideo, error)
	GetMovieVideo(id int, vid int) (*models.MovieVideo, error)
	GetMovieVideos(id int) ([]models.MovieVideo, error)
	DeleteMovieVideo(id int, vid int) error
	UpdateMovieVideo(movieVideo models.MovieVideo) error
	// Chat management
	GetMovieChatsHistory(id int) ([]models.Event, error)
	InsertMovieChat(chat models.MovieChat) error
	// User Management
	InsertUser(user models.User) (int, error)
	UpdateUser(user models.User) error
	DeleteUser(id int) error
}
