-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS feed (
    id SERIAL PRIMARY KEY,
    feed_link VARCHAR(255) UNIQUE
);

INSERT INTO feed (feed_link)
VALUES
    ('https://habr.com/ru/rss/all/all/'),
    ('https://dtf.ru/rss/'),
    ('https://www.it-world.ru/tech/products/rss/'),
    ('https://www.theverge.com/rss/index.xml'),
    ('https://www.engadget.com/rss.xml'),
    ('https://www.cnet.com/rss/all/'),
    ('https://www.zdnet.com/news/rss.xml'),
    ('https://www.wired.com/feed/rss')
ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS feed;
-- +goose StatementEnd
