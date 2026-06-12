-- name: CreateEncounter :one
INSERT INTO combat_encounters (name)
VALUES ($1)
RETURNING *;

-- name: GetEncounter :one
SELECT * FROM combat_encounters
WHERE id = $1;

-- name: ListEncounters :many
SELECT * FROM combat_encounters
ORDER BY created_at DESC;

-- name: GetActiveEncounter :one
SELECT * FROM combat_encounters
WHERE is_active = TRUE
LIMIT 1;

-- name: StartEncounter :one
UPDATE combat_encounters
SET is_active = TRUE,
    round     = 1
WHERE id = $1
RETURNING *;

-- name: EndEncounter :one
UPDATE combat_encounters
SET is_active = FALSE
WHERE id = $1
RETURNING *;

-- name: NextRound :one
UPDATE combat_encounters
SET round = round + 1
WHERE id = $1
RETURNING *;

-- name: DeleteEncounter :exec
DELETE FROM combat_encounters
WHERE id = $1;

-- name: GetParticipant :one
SELECT * FROM combat_participants
WHERE id = $1;

-- name: ListParticipants :many
SELECT * FROM combat_participants
WHERE encounter_id = $1
ORDER BY initiative DESC, name;

-- name: ListActiveParticipants :many
SELECT * FROM combat_participants
WHERE encounter_id = $1
AND is_active = TRUE
ORDER BY initiative DESC, name;

-- name: AddParticipant :one
INSERT INTO combat_participants (
    encounter_id,
    character_id,
    name,
    initiative,
    current_hp,
    max_hp,
    temp_hp,
    armor_class,
    speed
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: UpdateParticipantInitiative :one
UPDATE combat_participants
SET initiative = $2
WHERE id = $1
RETURNING *;

-- name: UpdateParticipantHP :one
UPDATE combat_participants
SET current_hp = $2
WHERE id = $1
RETURNING *;

-- name: UpdateParticipantTempHP :one
UPDATE combat_participants
SET temp_hp = $2
WHERE id = $1
RETURNING *;

-- name: UpdateParticipantConditions :one
UPDATE combat_participants
SET conditions = $2
WHERE id = $1
RETURNING *;

-- name: ToggleParticipantConcentration :one
UPDATE combat_participants
SET concentration = NOT concentration
WHERE id = $1
RETURNING *;

-- name: DeactivateParticipant :one
UPDATE combat_participants
SET is_active = FALSE
WHERE id = $1
RETURNING *;

-- name: RemoveParticipant :exec
DELETE FROM combat_participants
WHERE id = $1;
