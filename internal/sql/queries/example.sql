-- TODO: Edit this file to create your own queries.
-- Just some contrived examples to get you started.

-- name: List :many
SELECT * FROM book;

-- name: Get :one
SELECT * FROM book
WHERE id = $1;

-- name: Create :exec
INSERT INTO book (id, title, author, update_time)
VALUES ($1, $2, $3, $4);

-- name: BatchCreate :execrows
INSERT INTO book (id, title, author, update_time)
SELECT
    UNNEST(@id::TEXT[]),
    UNNEST(@title::TEXT[]),
    UNNEST(@author::TEXT[]),
    UNNEST(@update_time::BIGINT[])
ON CONFLICT (id) DO UPDATE
    SET
        title = EXCLUDED.title,
        author = EXCLUDED.author,
        update_time = EXCLUDED.update_time
    WHERE update_time < EXCLUDED.update_time; -- only update if the new value is newer

-- name: Upsert :exec
INSERT INTO book (id, title, author, update_time)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE
SET
    title = $2,
    author = $3,
    update_time = $4;

-- name: Update :exec
UPDATE book
SET
    title = COALESCE(NULLIF(@title::text, ''), title), -- should not be empty
    author = CASE @set_author::boolean WHEN true THEN @author ELSE author END, -- can be set to empty,
    update_time = @update_time
WHERE id = @id;

-- name: Delete :exec
DELETE FROM book
WHERE id = $1;
