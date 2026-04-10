-- name: GetSpell :one
SELECT * FROM spells
WHERE id = $1;

-- name: ListSpells :many
SELECT * FROM spells
WHERE character_id = $1
ORDER BY level, name;

-- name: ListPreparedSpells :many
SELECT * FROM spells
WHERE character_id = $1
AND is_prepared = TRUE
ORDER BY level, name;

-- name: CreateSpell :one
INSERT INTO spells (
    character_id,
    name,
    level,
    school,
    casting_time,
    range,
    components,
    duration,
    description,
    is_prepared
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: UpdateSpell :one
UPDATE spells
SET
    name         = COALESCE(sqlc.narg('name'), name),
    level        = COALESCE(sqlc.narg('level'), level),
    school       = COALESCE(sqlc.narg('school'), school),
    casting_time = COALESCE(sqlc.narg('casting_time'), casting_time),
    range        = COALESCE(sqlc.narg('range'), range),
    components   = COALESCE(sqlc.narg('components'), components),
    duration     = COALESCE(sqlc.narg('duration'), duration),
    description  = COALESCE(sqlc.narg('description'), description),
    is_prepared  = COALESCE(sqlc.narg('is_prepared'), is_prepared)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: DeleteSpell :exec
DELETE FROM spells
WHERE id = $1;

-- name: ToggleSpellPrepared :one
UPDATE spells
SET is_prepared = NOT is_prepared
WHERE id = $1
RETURNING *;
