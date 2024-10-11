-- name: ListMessages :many
SELECT * FROM message
ORDER BY create_time DESC;

-- name: GetMessage :one
SELECT * FROM message
WHERE id = ? LIMIT 1;

-- name: CreateMessage :exec
INSERT INTO message (author, message, create_time)
VALUES (?, ?, ?);

