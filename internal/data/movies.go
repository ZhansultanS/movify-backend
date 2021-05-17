package data

type Movie struct {
	ID          int64   `json:"id"`
	Title       string  `json:"title"`
	Overview    string  `json:"overview"`
	ReleaseDate string  `json:"release_date"`
	Runtime     Runtime `json:"runtime,omitempty"`
	Popularity  float32 `json:"popularity"`
	PosterPath  string  `json:"poster_path"`
	Genres      []Genre `json:"genres"`
}

type Genre struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
