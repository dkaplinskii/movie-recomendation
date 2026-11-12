-- +goose Up
CREATE TABLE movies
(
    id           SERIAL PRIMARY KEY,
    kinopoisk_id BIGINT NOT NULL,
    title        TEXT NOT NULL,
    description  TEXT NULL,
    genres       TEXT[] NOT NULL,
    poster_url   TEXT NULL,
    year         SMALLINT NOT NULL,
    countries    TEXT[] NOT NULL,
    trailer_url  TEXT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_movies_kinopoisk_id UNIQUE (kinopoisk_id)
);
-- +goose Down
DROP TABLE IF EXISTS movies;