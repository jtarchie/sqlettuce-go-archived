-- name: Set :exec
INSERT INTO keys (name, value)
VALUES (@name, @value) ON CONFLICT(name) DO
UPDATE
SET value = excluded.value;
-- name: AppendValue :one
INSERT INTO keys (name, value)
VALUES (@name, @value) ON CONFLICT(name) DO
UPDATE
SET value = value || excluded.value
RETURNING length(value);
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
-- name: ListSet :one
UPDATE keys
SET value = json_replace(
    value,
    '$[' || IIF(@index >= 0, @index, '#' || @index) || ']',
    @value
  )
WHERE name = @name
RETURNING json_valid(value);
-- name: ListRightPush :one
INSERT INTO keys (name, value)
VALUES (@name, json_insert('[]', '$[#]', @value)) ON CONFLICT(name) DO
UPDATE
SET value = json_insert(
    value,
    '$[#]',
    json_extract(excluded.value, '$[0]')
  )
RETURNING CAST(json_valid(value) AS boolean) AS valid,
  CAST(json_array_length(value) AS INTEGER) AS length;