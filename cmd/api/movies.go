package main

import (
	"fmt"
	"github.com/DARKestMODE/movify/internal/data"
	"net/http"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string       `json:"title"`
		ReleaseDate int32        `json:"release_date"`
		Runtime     int32        `json:"runtime"`
		Genres      []data.Genre `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
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
