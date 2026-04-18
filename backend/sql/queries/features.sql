-- name: ListFeatures :many
SELECT * FROM features
WHERE character_id = $1
ORDER BY action_type, name;

-- name: ListFeaturesByActionType :many
SELECT * FROM features
WHERE character_id = $1
AND action_type = $2
ORDER BY name;

-- name: CreateFeature :one
INSERT INTO features (
    character_id,
    name,
    source,
    description,
    action_type
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: UpdateFeature :one
UPDATE features
SET
    name        = COALESCE(sqlc.narg('name'), name),
    source      = COALESCE(sqlc.narg('source'), source),
    description = COALESCE(sqlc.narg('description'), description),
    action_type = COALESCE(sqlc.narg('action_type'), action_type)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteFeature :exec
DELETE FROM features
WHERE id = $1;
