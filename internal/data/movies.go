package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/DARKestMODE/movify/internal/validator"
	"github.com/lib/pq"
)

type Movie struct {
	Id          int64          `gorm:"primaryKey"`
	IdTMDB      int64          `json:"IdTMDB"`
	Title       string         `json:"title"`
	Overview    string         `json:"overview"`
	ReleaseDate string         `json:"release_date"`
	Runtime     Runtime        `json:"runtime"`
	Popularity  float32        `json:"popularity"`
	PosterPath  string         `json:"poster_path"`
	Genres      pq.StringArray `json:"genres" gorm:"type:text[]"`
}

func (m Movie) MarshalJSON() ([]byte, error) {
	var runtime string
	if m.Runtime != 0 {
		runtime = fmt.Sprintf("%d mins", m.Runtime)
	}

	type MovieAlias Movie

	aux := struct {
		MovieAlias
		Runtime string `json:"runtime,omitempty"`
	}{
		MovieAlias: MovieAlias(m),
		Runtime:    runtime,
	}

	return json.Marshal(aux)
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.IdTMDB != 0, "id_tmdb", "must be provided")
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

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(mv *Movie) error {
	query := `INSERT INTO movies (id_tmdb, title, overview, release_date, runtime, genres, popularity, poster_path)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			  RETURNING id`
	args := []interface{}{mv.IdTMDB, mv.Title, mv.Overview, mv.ReleaseDate, mv.Runtime, pq.Array(mv.Genres), mv.Popularity, mv.PosterPath}
	return m.DB.QueryRow(query, args...).Scan(&mv.Id)
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	return nil, nil
}

func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, error) {
	var movies []*Movie
	return movies, nil
}

func (m MovieModel) Update(movie *Movie) error {
	return nil
}

func (m MovieModel) Delete(id int64) error {
	return nil
}
