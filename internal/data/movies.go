package data

import (
	"database/sql"
	"errors"
	"github.com/DARKestMODE/movify/internal/validator"
	"github.com/lib/pq"
)

type Movie struct {
	Id          int64          `gorm:"primaryKey"`
	IdTMDB      int64          `json:"IdTMDB"`
	Title       string         `json:"title"`
	Overview    string         `json:"overview"`
	ReleaseDate string         `json:"release_date"`
	Runtime     int16          `json:"runtime"`
	Popularity  float32        `json:"popularity"`
	PosterPath  string         `json:"poster_path"`
	Genres      pq.StringArray `json:"genres"`
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
	q := `INSERT INTO movies (id_tmdb, title, overview, release_date, runtime, genres, popularity, poster_path)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		  RETURNING id`
	args := []interface{}{mv.IdTMDB, mv.Title, mv.Overview, mv.ReleaseDate, mv.Runtime, pq.Array(mv.Genres), mv.Popularity, mv.PosterPath}
	return m.DB.QueryRow(q, args...).Scan(&mv.Id)
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	q := `SELECT *
		  FROM movies
		  WHERE id = $1`

	var mv Movie
	var genres []sql.NullString
	err := m.DB.QueryRow(q, id).Scan(
		&mv.Id,
		&mv.IdTMDB,
		&mv.Title,
		&mv.Overview,
		&mv.ReleaseDate,
		&mv.Runtime,
		pq.Array(&genres),
		&mv.Popularity,
		&mv.PosterPath,
	)
	for _, g := range genres {
		if !g.Valid {
			continue
		}
		mv.Genres = append(mv.Genres, g.String)
	}
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &mv, nil
}

func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, error) {
	var movies []*Movie
	return movies, nil
}

func (m MovieModel) Update(mv *Movie) error {
	q := `UPDATE movies
		  SET id_tmdb = $2, title = $3, overview = $4, release_date = $5, runtime = $6, popularity = $7, poster_path = $8, genres = $9
		  WHERE id = $1`
	args := []interface{}{
		mv.Id, mv.IdTMDB, mv.Title,
		mv.Overview, mv.ReleaseDate, mv.Runtime,
		mv.Popularity, mv.PosterPath, pq.Array(mv.Genres),
	}

	return m.DB.QueryRow(q, args...).Err()
}

func (m MovieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	q := `DELETE FROM movies
		  WHERE id = $1`
	result, err := m.DB.Exec(q, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
