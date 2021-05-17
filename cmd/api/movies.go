package main

import (
	"fmt"
	"github.com/DARKestMODE/movify/internal/data"
	"github.com/DARKestMODE/movify/internal/validator"
	"net/http"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string       `json:"title"`
		Overview    string       `json:"overview"`
		ReleaseDate string       `json:"release_date"`
		Runtime     int32        `json:"runtime"`
		Popularity  float32      `json:"popularity"`
		PosterPath  string       `json:"poster_path"`
		Genres      []data.Genre `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:       input.Title,
		Overview:    input.Overview,
		ReleaseDate: input.ReleaseDate,
		Runtime:     input.Runtime,
		Popularity:  input.Popularity,
		PosterPath:  input.PosterPath,
		Genres:      input.Genres,
	}

	v := validator.New()

	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))


	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:          id,
		Title:       "Casablanca",
		Overview:    "BlaBla",
		ReleaseDate: "01-12-2020",
		Runtime:     102,
		Popularity:  10.5,
		PosterPath:  "/path.jpg",
		Genres: []data.Genre{{
			Id:   1,
			Name: "Drama",
		}},
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r, err)
	}
}
