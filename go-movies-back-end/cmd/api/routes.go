package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)
	mux.Get("/", app.Home)
	mux.Post("/authenticate", app.authenticate)
	mux.Get("/refresh", app.refreshToken)
	mux.Get("/logout", app.logout)
	mux.Get("/movies", app.AllMovies)
	mux.Get("/movie/{id}", app.GetMovie)
	mux.Get("/movie/{id}/video", app.MovieVideoDownload)
	mux.Get("/genres", app.AllGenres)
	mux.Get("/movies/genres/{id}", app.AllMoviesByGenre)

	mux.Post("/graph", app.moviesGraphQL)
	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.authRequired)
		mux.Get("/movies", app.MoviesCatalog)
		mux.Get("/movies/{id}", app.MovieForEdit)
		mux.Put("/movies/0", app.InsertMovie)
		mux.Patch("/movies/{id}", app.UpdateMovie)
		mux.Delete("/movies/{id}", app.DeleteMovie)
		mux.Post("/movies/{id}/video", app.MovieVidoeUpload)
		mux.Patch("/movies/{id}/videos/{vid}", app.UpdateMovieVideo)
		mux.Get("/movies/{id}/videos", app.MovieVideos)
		mux.Delete("/movies/{id}/videos/{vid}", app.DeleteMovieVideo)
	})
	return mux
}
