-- +migrate Up
CREATE TABLE message (
    id          INTEGER PRIMARY KEY,
    author      TEXT NOT NULL,
    message     TEXT NOT NULL,
    create_time BIGINT NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS message;
