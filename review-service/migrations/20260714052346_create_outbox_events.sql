-- +goose Up
CREATE TABLE outbox_events
(
    event_id   SERIAL PRIMARY KEY,
    user_id    INT NOT NULL,
    film_id    INT NOT NULL,
    review     SMALLINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose Down
DROP TABLE IF EXISTS outbox_events;