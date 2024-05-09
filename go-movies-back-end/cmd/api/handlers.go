package main

import (
	"backend/internal/graph"
	"backend/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go Movies up and running",
		Version: "1.0.0",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) AllMovies(w http.ResponseWriter, r *http.Request) {
	// var movies []models.Movie
	// rd, _ := time.Parse("2006-01-02", "1986-03-07")
	// highlander := models.Movie{
	// 	ID:          1,
	// 	Title:       "Highlander",
	// 	ReleaseDate: rd,
	// 	MPAARating:  "R",
	// 	RunTime:     116,
	// 	Description: "A very nice highlander movie",
	// 	CreatedAt:   time.Now(),
	// 	UpdatedAt:   time.Now(),
	// }
	// movies = append(movies, highlander)

	// rd, _ = time.Parse("2006-01-02", "1981-06-12")
	// rotla := models.Movie{
	// 	ID:          2,
	// 	Title:       "Raiders of the Lost Ark",
	// 	ReleaseDate: rd,
	// 	MPAARating:  "PG-13",
	// 	RunTime:     115,
	// 	Description: "A very nice Raiders movie",
	// 	CreatedAt:   time.Now(),
	// 	UpdatedAt:   time.Now(),
	// }
	//movies = append(movies, rotla)

	movies, err := app.DB.AllMovies()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movies)
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	//read json payload
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}
	// validate user against database
	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		fmt.Println(err)
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}
	// check password
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}
	// create a jwt user
	u := jwtUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
	//generate tokens
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	//log.Println(tokens.Token)
	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)
	app.writeJSON(w, http.StatusAccepted, tokens)
}

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value
			// parse the token to get claims
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}
			// get the user id from token claim

			userID, err := strconv.Atoi(claims.Subject)

			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}
			//fmt.Println(userID)
			user, err := app.DB.GetUserByID(userID)
			if err != nil {
				fmt.Println(err)
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			u := jwtUser{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			}

			tokenPairs, err := app.auth.GenerateTokenPair(&u)
			if err != nil {
				app.errorJSON(w, errors.New("error generating token"), http.StatusUnauthorized)
				return
			}
			fmt.Println(tokenPairs)
			http.SetCookie(w, app.auth.GetRefreshCookie(tokenPairs.Token))
			app.writeJSON(w, http.StatusOK, tokenPairs)

		}
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}

func (app *application) MoviesCatalog(w http.ResponseWriter, r *http.Request) {
	movies, err := app.DB.AllMovies()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movies)
}

func (app *application) GetMovie(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movieID, err := strconv.Atoi(id)

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	movie, err := app.DB.OneMovie(movieID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movie)

}

func (app *application) MovieForEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movieID, err := strconv.Atoi(id)

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	movie, genres, err := app.DB.OneMovieForEdit(movieID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload = struct {
		Movie  *models.Movie   `json:"movie"`
		Genres []*models.Genre `json:"genres"`
	}{
		movie,
		genres,
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) AllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := app.DB.AllGenres()
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	_ = app.writeJSON(w, http.StatusOK, genres)
}

func (app *application) InsertMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie

	err := app.readJSON(w, r, &movie)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	// try to get an image
	movie = app.getPoster(movie)

	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()

	newID, err := app.DB.InsertMovie(movie)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	// now handle genres
	err = app.DB.UpdateMovieGenres(newID, movie.GenresArray)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	resp := JSONReponse{
		Error:   false,
		Message: "movie updated",
	}
	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) getPoster(movie models.Movie) models.Movie {
	type TheMovieDB struct {
		Page    int `json:"page"`
		Results []struct {
			PosterPath string `json:"poster_path"`
		} `json:"results"`
		TotalPages int `json:"total_pages"`
	}

	client := &http.Client{}
	theUrl := fmt.Sprintf("https://api.themoviedb.org/3/search/movie?api_key=%s", app.APIKey)
	req, err := http.NewRequest("GET", theUrl+"&query="+url.QueryEscape(movie.Title), nil)

	if err != nil {
		log.Println(err)
		return movie
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return movie
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return movie
	}
	var responseObj TheMovieDB

	json.Unmarshal(bodyBytes, responseObj)

	if len(responseObj.Results) > 0 {
		movie.Image = responseObj.Results[0].PosterPath
	}
	return movie
}

func (app *application) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	var payload models.Movie

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movie, err := app.DB.OneMovie(payload.ID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	movie.Title = payload.Title
	movie.Description = payload.Description
	movie.ReleaseDate = payload.ReleaseDate
	movie.MPAARating = payload.MPAARating
	movie.RunTime = payload.RunTime
	movie.UpdatedAt = time.Now()

	err = app.DB.UpdateMovie(*movie)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	// now handle genres
	err = app.DB.UpdateMovieGenres(payload.ID, payload.GenresArray)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	resp := JSONReponse{
		Error:   false,
		Message: "movie updated",
	}
	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movieID, err := strconv.Atoi(id)

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	err = app.DB.DeleteMovie(movieID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	resp := JSONReponse{
		Error:   false,
		Message: "movie deleted",
	}
	_ = app.writeJSON(w, http.StatusOK, resp)

}

func (app *application) AllMoviesByGenre(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movieID, err := strconv.Atoi(id)

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	movies, err := app.DB.AllMovies(movieID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, movies)

}

func (app *application) moviesGraphQL(w http.ResponseWriter, r *http.Request) {
	// get all the movies
	movies, _ := app.DB.AllMovies()
	// extract query string from request body
	q, _ := io.ReadAll(r.Body)
	// decode the query string
	query := string(q)

	// create a new graphql variable
	g := graph.New(movies)

	// Set the query string
	g.QueryString = query

	// perform the query
	resp, err := g.Query()
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	// write the response
	j, _ := json.MarshalIndent(resp, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// Upload Video Functinality
func (app *application) MovieUpload(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movieID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	// lets upload the video of the movie
	fmt.Printf("upload recieved-%d", movieID)
	moviefile, moviefile_header, err := app.readMultiPartForm(r)
	if err != nil {
		fmt.Println(err)
		app.errorJSON(w, err)
		return
	}
	fmt.Println("Calling video upload")
	filefile_ext := strings.Split(moviefile_header.Filename, ".")[1]
	path, err := app.Storage.UploadVideo(moviefile, filefile_ext)
	if err != nil {
		fmt.Println(err)
		app.errorJSON(w, err)
		return
	}
	// insert movie video into database
	var movieVideo models.MovieVideo
	movieVideo.MovieID = movieID
	movieVideo.VideoPath = path
	movieVideo.CreatedAt = time.Now()
	movieVideo.IsLatest = true
	fmt.Println("Adding  video upload into database")
	err = app.DB.InsertMovieVideo(movieVideo)
	// return response back
	if err != nil {
		fmt.Println(err)
		app.errorJSON(w, err)
		return
	}
	resp := JSONReponse{
		Error:   false,
		Message: "movie video uploaded",
	}
	_ = app.writeJSON(w, http.StatusOK, resp)
}

// Movie Video Download
func (app *application) MovieDownload(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movieID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// get the movie video from database
	movieVideo, err := app.DB.GetMovieVideo(movieID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	// get the video file from s3
	video, videoInfo, err := app.Storage.GetVideo(movieVideo)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	// write the video file to the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", videoInfo.Size()))

	_, err = io.Copy(w, video)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer video.Close()

}
