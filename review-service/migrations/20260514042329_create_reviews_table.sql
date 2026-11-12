-- +goose Up
CREATE TABLE reviews
(
    id         SERIAL PRIMARY KEY,
    user_id    INT NOT NULL,
    film_id    INT NOT NULL,
    review     SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT reviews_user_id_film_id_unique UNIQUE (user_id, film_id)
);
-- +goose Down
DROP TABLE IF EXISTS reviews;