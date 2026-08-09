-- name: GetFeature :one
SELECT * FROM features
WHERE id = $1;

-- name: ListFeatures :many
SELECT * FROM features
WHERE character_id = $1
ORDER BY source, name;

-- name: CreateFeature :one
INSERT INTO features (
    character_id,
    name,
    action_type,
    source,
    description
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateFeature :one
UPDATE features
SET
    name         = COALESCE(sqlc.narg('name'), name),
    action_type  = COALESCE(sqlc.narg('action_type'), action_type),
    source       = COALESCE(sqlc.narg('source'), source),
    description  = COALESCE(sqlc.narg('description'), description)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteFeature :exec
DELETE FROM features
WHERE id = $1;
