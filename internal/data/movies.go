package data

import (
	"github.com/DARKestMODE/movify/internal/validator"
)

type Movie struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	ReleaseDate string  `json:"release_date"`
	Runtime     int32   `json:"runtime"`
	Popularity  float32 `json:"popularity"`
	PosterPath  string  `json:"poster_path"`
	Genres      []Genre `json:"genres"`
}

type Genre struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(movie.Overview != "", "overview", "must be provided")
	v.Check(movie.ReleaseDate != "", "release date", "must be provided")
	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(movie.Popularity != 0, "popularity", "must be provided")
	v.Check(movie.Popularity > 0, "popularity", "must be a positive number")
	v.Check(movie.PosterPath != "", "poster path", "must be provided")
	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}
