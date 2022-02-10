package main

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

type ctxKey string

const wrapKey ctxKey = "params"

func (app *application) wrap(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), wrapKey, ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (app *application) routes() http.Handler {
	router := httprouter.New()

	secure := alice.New(app.checkToken)

	// Get API info
	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodPost, "/v1/graphql/list", app.moviesGQL)

	// Sign in method
	router.HandlerFunc(http.MethodPost, "/v1/signin", app.signIn)

	// Get a movie data
	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getOneMovie)

	// Get all movie data
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllMovies)

	// Get a list of movie by its genre
	router.HandlerFunc(http.MethodGet, "/v1/movies/:genre_id", app.getAllMOviesByGenre)

	// Get all genre data
	router.HandlerFunc(http.MethodGet, "/v1/genres", app.getAllGenres)

	// Edit existing data
	router.POST("/v1/admin/editmovie", app.wrap(secure.ThenFunc(app.editMovie)))

	// Delete existing data
	router.GET("/v1/admin/deletemovie/:id", app.wrap(secure.ThenFunc(app.deleteMovie)))

	return app.enableCORS(router)
}
