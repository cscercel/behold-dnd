-- name: ListCharacters :many
SELECT *
FROM characters;

-- name: FindCharacterByID :one
SELECT *
FROM characters
WHERE id = $1;
