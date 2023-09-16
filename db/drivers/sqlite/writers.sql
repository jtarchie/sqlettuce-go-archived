-- name: Set :exec
INSERT INTO keys (name, value)
VALUES (@name, @value) ON CONFLICT(name) DO
UPDATE
SET value = excluded.value;

-- name: Append :one
INSERT INTO keys (name, value)
VALUES (@name, @value) ON CONFLICT(name) DO
UPDATE
SET value = value || excluded.value RETURNING length(value);

-- name: AddFloat :one
INSERT INTO keys (name, value)
VALUES (@name, @value) ON CONFLICT(name) DO
UPDATE
SET value = CAST(value AS REAL) + CAST(excluded.value AS REAL)
WHERE printf("%.17f", value) GLOB SUBSTRING(value, 1, 1) || '*'
RETURNING CAST(value AS REAL);

-- name: AddInt :one
INSERT INTO keys (name, value)
VALUES (@name, @value) ON CONFLICT(name) DO
UPDATE
SET value = CAST(value AS INTEGER) + CAST(excluded.value AS INTEGER)
WHERE printf("%d", value) = value
RETURNING CAST(value AS INTEGER);

-- name: FlushAll :exec
DELETE FROM keys;