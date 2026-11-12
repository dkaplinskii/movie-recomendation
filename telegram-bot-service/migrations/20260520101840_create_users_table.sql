-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    telegram_id BIGINT UNIQUE,
    telegram_chat_id BIGINT
);

-- +goose Down
DROP TABLE IF EXISTS users;
