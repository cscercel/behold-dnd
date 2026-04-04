-- name: ListSpellSlots :many
SELECT * FROM spell_slots
WHERE character_id = $1
ORDER BY spell_level;

-- name: UpsertSpellSlot :one
INSERT INTO spell_slots (character_id, spell_level, total, used)
VALUES ($1, $2, $3, $4)
ON CONFLICT (character_id, spell_level)
DO UPDATE SET
    total = EXCLUDED.total,
    used = EXCLUDED.used
RETURNING *;

-- name: UseSpellSlot :one
UPDATE spell_slots
SET used = used + 1
WHERE character_id = $1
AND spell_level = $2
AND used < total
RETURNING *;

-- name: ResetSpellSlots :exec
UPDATE spell_slots
SET used = 0
WHERE character_id = $1;
