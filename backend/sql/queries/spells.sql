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
    name            = $2,
    level           = $3,
    school          = $4,
    casting_time    = $5,
    range           = $6,
    components      = $7,
    duration        = $8,
    description     = $9,
    is_prepared     = $10
WHERE id = $1
RETURNING *;

-- name: DeleteSpell :exec
DELETE FROM spells
WHERE id = $1;

-- name: ToggleSpellPrepared :one
UPDATE spells
SET is_prepared = NOT is_prepared
WHERE id = $1
RETURNING *;
