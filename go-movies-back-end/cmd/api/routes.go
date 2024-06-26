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
	mux.Get("/movies/{id}", app.GetMovie)
	mux.Get("/movies/{id}/video", app.MovieVideoDownload)
	mux.Get("/genres", app.AllGenres)
	mux.Get("/movies/genres/{id}", app.AllMoviesByGenre)
	mux.Post("/graph", app.moviesGraphQL)
	mux.Post("/signup", app.CreateUser)
	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(app.authRequiredAdmin)
		mux.Get("/movies", app.MoviesCatalog)
		mux.Get("/movies/{id}", app.MovieForEdit)
		mux.Put("/movies/0", app.InsertMovie)
		mux.Patch("/movies/{id}", app.UpdateMovie)
		mux.Delete("/movies/{id}", app.DeleteMovie)
		mux.Post("/movies/{id}/video", app.MovieVidoeUpload)
		mux.Patch("/movies/{id}/videos/{vid}", app.UpdateMovieVideo)
		mux.Get("/movies/{id}/videos", app.MovieVideos)
		mux.Delete("/movies/{id}/videos/{vid}", app.DeleteMovieVideo)
		mux.Get("/movies/{id}/chats", app.MovieChatsHistory)
		mux.Get("/movies/{id}/chatws", app.MovieChatsWS)
	})

	mux.Route("/protected", func(mux chi.Router) {
		mux.Use(app.authRequiredProtected)
		mux.Get("/movies/{id}/chats", app.MovieChatsHistory)
		mux.Get("/movies/{id}/chatws", app.MovieChatsWS)
	})
	return mux
}
