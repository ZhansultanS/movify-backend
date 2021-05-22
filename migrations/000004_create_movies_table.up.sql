CREATE TABLE IF NOT EXISTS movies
(
    id           bigserial PRIMARY KEY,
    id_tmdb      bigint UNIQUE               NOT NULL,
    title        text                        NOT NULL,
    overview     text                        NOT NULL,
    release_date text                        NOT NULL,
    runtime      integer                     NOT NULL,
    genres       text[]                      NOT NULL,
    popularity   numeric                     NOT NULL,
    poster_path  text                        NOT NULL,
    created_at   timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version      integer                     NOT NULL DEFAULT 1
);