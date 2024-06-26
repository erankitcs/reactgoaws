package dbrepo

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

type PostgresDBRepo struct {
	DB *sql.DB
}

const dbTimeout = time.Second * 3

func (m *PostgresDBRepo) Connection() *sql.DB {
	return m.DB
}

func (m *PostgresDBRepo) AllMovies(genre ...int) ([]*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	filter := ""
	if len(genre) > 0 {
		filter = fmt.Sprintf(" where id in ( select movie_id from movies_genres where genre_id=%d)", genre[0])
	}

	query := fmt.Sprintf(`
		SELECT 
		  id, title, release_date, runtime, 
		  mpaa_rating, description,coalesce(image, ''),
		  created_at, updated_at
		FROM 
		  movies %s
		order by 
		  title	
	`, filter)
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*models.Movie

	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.RunTime,
			&movie.MPAARating,
			&movie.Description,
			&movie.Image,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}

	return movies, nil
}

func (m *PostgresDBRepo) OneMovie(id int) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT 
		  id, title, release_date, runtime, 
		  mpaa_rating, description,coalesce(image, ''),
		  created_at, updated_at
		FROM 
		  movies
		where id = $1
	`
	row := m.DB.QueryRowContext(ctx, query, id)

	var movie models.Movie
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.ReleaseDate,
		&movie.RunTime,
		&movie.MPAARating,
		&movie.Description,
		&movie.Image,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// get genres, if any
	query = `
		select
		  g.id, g.genre from movies_genres mg left join genres g on ( mg.genre_id = g.id )
		where mg.movie_id =$1
		order by g.genre
	`
	rows, err := m.DB.QueryContext(ctx, query, id)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	defer rows.Close()

	var genres []*models.Genre

	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID,
			&g.Genre,
		)

		if err != nil {
			return nil, err
		}

		genres = append(genres, &g)
	}
	movie.Genres = genres
	return &movie, nil
}

func (m *PostgresDBRepo) OneMovieForEdit(id int) (*models.Movie, []*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT 
		  id, title, release_date, runtime, 
		  mpaa_rating, description,coalesce(image, ''),
		  created_at, updated_at
		FROM 
		  movies
		where id = $1
	`
	row := m.DB.QueryRowContext(ctx, query, id)

	var movie models.Movie
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.ReleaseDate,
		&movie.RunTime,
		&movie.MPAARating,
		&movie.Description,
		&movie.Image,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)

	if err != nil {
		return nil, nil, err
	}

	// get genres, if any
	query = `
		select
		  g.id, g.genre from movies_genres mg left join genres g on ( mg.genre_id = g.id )
		where mg.movie_id =$1
		order by g.genre
	`
	rows, err := m.DB.QueryContext(ctx, query, id)

	if err != nil && err != sql.ErrNoRows {
		return nil, nil, err
	}

	defer rows.Close()

	var genres []*models.Genre
	var genresArray []int

	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID,
			&g.Genre,
		)

		if err != nil {
			return nil, nil, err
		}

		genres = append(genres, &g)
		genresArray = append(genresArray, g.ID)
	}
	movie.Genres = genres
	movie.GenresArray = genresArray

	var allGenres []*models.Genre

	query = `
		select
		  id, genre
		from genres
		order by genre
	`
	gRows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, nil, err
	}

	defer gRows.Close()

	for gRows.Next() {
		var g models.Genre
		gRows.Scan(
			&g.ID,
			&g.Genre,
		)
		allGenres = append(allGenres, &g)
	}

	return &movie, allGenres, nil
}

func (m *PostgresDBRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
		SELECT 
		  id, email, first_name, last_name, password, created_at, updated_at
		FROM 
		  users
		where
		  approved = 't' and
		  email = $1
	`
	var user models.User
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// All the roles for the user
	// get all the roles from the roles table filter by user id
	// and add to the user struct
	query = `
		select
		  role
		from
		  users_roles
		where
		  user_id = $1
	`
	rolesRows, err := m.DB.QueryContext(ctx, query, user.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	defer rolesRows.Close()

	var roles []*models.Role

	for rolesRows.Next() {
		var role models.Role
		rolesRows.Scan(
			&role.Role,
		)
		roles = append(roles, &role)
	}

	user.Roles = roles
	return &user, nil
}

func (m *PostgresDBRepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
		SELECT 
		  id, email, first_name, last_name, password, created_at, updated_at
		FROM 
		  users
		where
		  approved = 't' and
		  id = $1
	`
	var user models.User
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// All the roles for the user
	// get all the roles from the roles table filter by user id
	// and add to the user struct
	query = `
		select
		  role
		from
		  users_roles
		where
		  user_id = $1
	`
	rolesRows, err := m.DB.QueryContext(ctx, query, user.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	defer rolesRows.Close()

	var roles []*models.Role

	for rolesRows.Next() {
		var role models.Role
		rolesRows.Scan(
			&role.Role,
		)
		roles = append(roles, &role)
	}

	user.Roles = roles

	return &user, nil
}

func (m *PostgresDBRepo) AllGenres() ([]*models.Genre, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
		SELECT 
		  id, genre, created_at, updated_at
		FROM 
		  genres
		order by 
		  genre
	`
	var genres []*models.Genre
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var g models.Genre
		rows.Scan(
			&g.ID,
			&g.Genre,
			&g.CreatedAt,
			&g.UpdatedAt,
		)
		genres = append(genres, &g)
	}
	return genres, nil
}

func (m *PostgresDBRepo) InsertMovie(movie models.Movie) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `insert into movies (title, description, release_date, runtime, mpaa_rating, created_at, updated_at, image)
			values ($1, $2, $3, $4, $5, $6, $7, $8) returning id`
	var newID int
	err := m.DB.QueryRowContext(ctx, stmt,
		movie.Title,
		movie.Description,
		movie.ReleaseDate,
		movie.RunTime,
		movie.MPAARating,
		movie.CreatedAt,
		movie.UpdatedAt,
		movie.Image,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *PostgresDBRepo) UpdateMovie(movie models.Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `update movies set title = $1, description = $2, release_date = $3,
			runtime = $4, mpaa_rating = $5, updated_at = $6, image = $7
			where id = $8`
	_, err := m.DB.ExecContext(ctx, stmt,
		movie.Title,
		movie.Description,
		movie.ReleaseDate,
		movie.RunTime,
		movie.MPAARating,
		movie.UpdatedAt,
		movie.Image,
		movie.ID,
	)

	if err != nil {
		return err
	}
	return nil
}

func (m *PostgresDBRepo) UpdateMovieGenres(id int, genreIDs []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `delete from movies_genres where movie_id = $1`

	_, err := m.DB.ExecContext(ctx, stmt, id)

	if err != nil {
		return err
	}

	for _, n := range genreIDs {
		stmt = `insert into movies_genres (movie_id, genre_id) values ($1, $2)`
		_, err := m.DB.ExecContext(ctx, stmt, id, n)
		if err != nil {
			return err
		}
	}

	return nil

}

func (m *PostgresDBRepo) DeleteMovie(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `delete from movies where id = $1`
	_, err := m.DB.ExecContext(ctx, stmt, id)

	if err != nil {
		return err
	}

	stmt = `delete from movies_genres where movie_id = $1`

	_, err = m.DB.ExecContext(ctx, stmt, id)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostgresDBRepo) InsertMovieVideo(movieVideo models.MovieVideo) (*models.MovieVideo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	// Update all other movie videos for this movie to no longer be the latest
	stmt := `update movies_videos set is_latest = 'f'
			where movie_id = $1`
	_, err := m.DB.ExecContext(ctx, stmt,
		movieVideo.MovieID,
	)

	if err != nil {
		return nil, err
	}

	stmt = `insert into movies_videos (movie_id, video_path, is_latest, created_at)
			values ($1, $2, $3, $4) returning id`
	//var vid int
	err = m.DB.QueryRowContext(ctx, stmt,
		movieVideo.MovieID,
		movieVideo.VideoPath,
		movieVideo.IsLatest,
		movieVideo.CreatedAt,
	).Scan(
		&movieVideo.ID,
	)

	if err != nil {
		return nil, err
	}

	return &movieVideo, nil
}

// Get Video from Movies_Videos table
func (m *PostgresDBRepo) GetMovieVideo(id int, vid int) (*models.MovieVideo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	filter := ""
	// if vid is positive then return that else return latest
	if vid > 0 {
		filter = fmt.Sprintf(" and id = %d", vid)
	} else {
		filter = " and is_latest = 't'"
	}
	defer cancel()
	query := `
		SELECT
		 id, movie_id, video_path, is_latest, created_at
		FROM
		  movies_videos
		where
		  movie_id = $1 ` + filter + `
		limit 1
	`
	var video models.MovieVideo
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&video.ID,
		&video.MovieID,
		&video.VideoPath,
		&video.IsLatest,
		&video.CreatedAt,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &video, nil
}

// Get All Videos history from Movies_Videos table for given movie
func (m *PostgresDBRepo) GetMovieVideos(id int) ([]models.MovieVideo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
		SELECT
		  id, movie_id, video_path, is_latest, created_at
		FROM
		  movies_videos
		where
		  movie_id = $1
		order by created_at desc
	`
	rows, err := m.DB.QueryContext(ctx, query, id)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	defer rows.Close()

	var videos []models.MovieVideo

	for rows.Next() {
		var video models.MovieVideo
		err := rows.Scan(
			&video.ID,
			&video.MovieID,
			&video.VideoPath,
			&video.IsLatest,
			&video.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		videos = append(videos, video)
	}
	return videos, nil
}

// Delete a movie video from the Movies_Videos table
func (m *PostgresDBRepo) DeleteMovieVideo(id int, vid int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `delete from movies_videos where id = $1 and movie_id=$2`
	_, err := m.DB.ExecContext(ctx, stmt, vid, id)

	if err != nil {
		return err
	}

	return nil
}

// Update a movie video from the Movies_Videos table with movievideo object
func (m *PostgresDBRepo) UpdateMovieVideo(movieVideo models.MovieVideo) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `update movies_videos set is_latest = $1
			where id = $2 and movie_id = $3`
	_, err := m.DB.ExecContext(ctx, stmt,
		movieVideo.IsLatest,
		movieVideo.ID,
		movieVideo.MovieID,
	)

	if err != nil {
		return err
	}

	return nil
}

// Function to read all the chat history for given movieID from movies_chats table
// Join users table as well to get username
func (m *PostgresDBRepo) GetMovieChatsHistory(id int) ([]models.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `
		SELECT
		  m.chattext, CONCAT (u.first_name,' ', u.last_name ) as full_name , m.created_at
		FROM
		  movies_chats m
		  left join users u on (m.user_id = u.id)
		where
		  movie_id = $1
	`
	rows, err := m.DB.QueryContext(ctx, query, id)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	defer rows.Close()

	var chatEvents []models.Event

	for rows.Next() {
		var chatEvent models.Event
		chatEvent.Type = "chat_history"
		var chatPayload models.NewMessageEvent
		var chat models.SendMessageEvent
		err := rows.Scan(
			&chat.Message,
			&chatPayload.From,
			&chatPayload.Sent,
		)

		if err != nil {
			return nil, err
		}
		chatPayload.SendMessageEvent = chat
		//Marshal the chat payload
		chatJSONPayload, err := json.Marshal(chatPayload)
		if err != nil {
			return nil, err
		}
		chatEvent.Payload = chatJSONPayload

		chatEvents = append(chatEvents, chatEvent)
	}
	return chatEvents, nil
}

// Function to insert a new chat message into movies_chats table
func (m *PostgresDBRepo) InsertMovieChat(chat models.MovieChat) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `insert into movies_chats (movie_id, user_id, chattext, created_at)
			values ($1, $2, $3, $4)`
	_, err := m.DB.ExecContext(ctx, stmt,
		chat.MovieID,
		chat.UserID,
		chat.ChatText,
		chat.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

// Function to insert a user into users table
func (m *PostgresDBRepo) InsertUser(user models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `insert into users (email, first_name, last_name, password, approved, created_at, updated_at)
			values ($1, $2, $3, $4, $5, $6, $7) returning id`
	var newID int
	// Begin a transaction
	tx, err := m.DB.BeginTx(ctx, nil)
	// Handle error
	if err != nil {
		return 0, err
	}
	// Close the transaction at the end
	defer tx.Rollback()

	// Execute the query
	err = tx.QueryRowContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Password,
		user.Approved,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}
	// Inserting Roles into users_roles table
	// Loop through all the roles and insert into users_roles table
	for _, role := range user.Roles {
		stmt = `insert into users_roles (user_id, role, created_at, updated_at) values ($1, $2, $3, $4)`
		_, err := tx.ExecContext(ctx, stmt, newID, role.Role, user.CreatedAt, user.UpdatedAt)
		if err != nil {
			return 0, err
		}
	}

	// Committing the transaction
	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return newID, nil
}

// Function to delete a user in users table
func (m *PostgresDBRepo) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `delete from users where id = $1`
	_, err := m.DB.ExecContext(ctx, stmt, id)

	if err != nil {
		return err
	}

	return nil
}

// Function to update a user in users table
func (m *PostgresDBRepo) UpdateUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	stmt := `update users set email = $1, first_name = $2, last_name = $3, password = $4, updated_at = $5
			where id = $6`
	_, err := m.DB.ExecContext(ctx, stmt,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Password,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return err
	}

	return nil
}
