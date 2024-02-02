-- TODO: Edit this file to create your own schema

-- +migrate Up
CREATE TABLE book (
    id          TEXT PRIMARY KEY,
    title       TEXT NOT NULL,
    author      TEXT NOT NULL,
    update_time BIGINT NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS my_table;
