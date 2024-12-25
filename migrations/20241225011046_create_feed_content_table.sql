-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS feed_content (
    item_id SERIAL PRIMARY KEY,
    feed_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    published_at TIMESTAMP NOT NULL,
    pub_link TEXT NOT NULL,
    is_read BOOL DEFAULT FALSE NOT NULL,
    FOREIGN KEY (feed_id) REFERENCES feed(id) ON DELETE CASCADE
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS feed_content
-- +goose StatementEnd
