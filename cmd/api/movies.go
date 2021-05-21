package main

import (
	"errors"
	"fmt"
	"github.com/DARKestMODE/movify/internal/data"
	"github.com/DARKestMODE/movify/internal/validator"
	"net/http"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		IdTMDB      int64    `json:"id_tmdb"`
		Title       string   `json:"title"`
		Overview    string   `json:"overview"`
		ReleaseDate string   `json:"release_date"`
		Runtime     int16    `json:"runtime"`
		Genres      []string `json:"genres"`
		Popularity  float32  `json:"popularity"`
		PosterPath  string   `json:"poster_path"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	mv := &data.Movie{
		IdTMDB:      input.IdTMDB,
		Title:       input.Title,
		Overview:    input.Overview,
		ReleaseDate: input.ReleaseDate,
		Runtime:     input.Runtime,
		Popularity:  input.Popularity,
		PosterPath:  input.PosterPath,
		Genres:      input.Genres,
	}

	v := validator.New()

	if data.ValidateMovie(v, mv); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Insert(mv)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", mv.Id))

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": mv}, headers)
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

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		IdTMDB      *int64   `json:"id_tmdb"`
		Title       *string  `json:"title"`
		Overview    *string  `json:"overview"`
		ReleaseDate *string  `json:"release_date"`
		Runtime     *int16   `json:"runtime"`
		Genres      []string `json:"genres"`
		Popularity  *float32 `json:"popularity"`
		PosterPath  *string  `json:"poster_path"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.IdTMDB != nil {
		movie.IdTMDB = *input.IdTMDB
	}
	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Overview != nil {
		movie.Overview = *input.Overview
	}
	if input.ReleaseDate != nil {
		movie.ReleaseDate = *input.ReleaseDate
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres
	}
	if input.Popularity != nil {
		movie.Popularity = *input.Popularity
	}
	if input.PosterPath != nil {
		movie.PosterPath = *input.PosterPath
	}

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		data.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genres = app.readCSV(qs, "genres", []string{})
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "title", "release_date", "runtime", "popularity", "-id", "-title", "-release_date", "-runtime", "-popularity"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	movies, metadata, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movies": movies, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
