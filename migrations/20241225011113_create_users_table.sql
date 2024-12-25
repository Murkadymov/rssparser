-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users
    (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(60) NOT NULL,
    createdAt TIMESTAMP NOT NULL
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS users
-- +goose StatementEnd
