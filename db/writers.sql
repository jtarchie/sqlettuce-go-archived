-- name: Set :exec
INSERT INTO strings (name, value)
VALUES (@name, @value) ON CONFLICT(name) DO
UPDATE
SET value = excluded.value;

-- name: Append :one
INSERT INTO strings (name, value)
VALUES (@name, @value) ON CONFLICT(name) DO
UPDATE
SET value = value || excluded.value RETURNING length(value);

-- name: FlushAll :exec
DELETE FROM strings;

-- name: Delete :exec
DELETE FROM strings WHERE name = @name;