-- name: CreateCharacter :one
INSERT INTO characters (
    owner_id, is_npc,
    name, race, class, level, background, alignment, xp,
    strength, dexterity, constitution, intelligence, wisdom, charisma,
    save_prof_strength, save_prof_dexterity, save_prof_constitution,
    save_prof_intelligence, save_prof_wisdom, save_prof_charisma,
    skill_acrobatics, skill_animal_handling, skill_arcana, skill_athletics,
    skill_deception, skill_history, skill_insight, skill_intimidation,
    skill_investigation, skill_medicine, skill_nature, skill_perception,
    skill_performance, skill_persuasion, skill_religion, skill_sleight_of_hand,
    skill_stealth, skill_survival,
    max_hp, current_hp, temp_hp, armor_class, speed,
    hit_dice_type, hit_dice_remaining,
    inspiration, attunement_slots,
    training_armor, training_weapons, training_tools, training_languages,
    copper, silver, electrum, gold, platinum,
    conditions, resistances, vulnerabilities, immunities,
    personality_traits, ideals, bonds, flaws, notes,
    spellcasting_ability
) VALUES (
    $1,  $2,  $3,  $4,  $5,  $6,  $7,  $8,  $9,
    $10, $11, $12, $13, $14, $15,
    $16, $17, $18, $19, $20, $21,
    $22, $23, $24, $25, $26, $27, $28, $29, $30,
    $31, $32, $33, $34, $35, $36, $37, $38, $39,
    $40, $41, $42, $43, $44, $45, $46,
    $47, $48,
    $49, $50, $51, $52,
    $53, $54, $55, $56, $57,
    $58, $59, $60, $61,
    $62, $63, $64, $65, $66,
    $67
)
RETURNING *;

-- name: GetCharacter :one
SELECT 
    *,
    proficiency_bonus = CEILING(1 + level / 4),

    strength_modifier = FLOOR((strength - 10 / 2)),
    dexterity_modifier = FLOOR((dexterity - 10 / 2)),
    constitution_modifier = FLOOR((constitution - 10 / 2)),
    intelligence_modifier = FLOOR((intelligence - 10 / 2)),
    wisdom_modifier = FLOOR((wisdom - 10 / 2)),
    charisma_modifier = FLOOR((charisma - 10 / 2)),

    strength_saving_throw_modifier = strength_modifier + save_prof_strength * proficiency_bonus,
    dexterity_saving_throw_modifier = dexterity_modifier + save_prof_dexterity * proficiency_bonus,
    constitution_saving_throw_modifier = constitution_modifier + save_prof_constitution * proficiency_bonus,
    intelligence_saving_throw_modifier = intelligence_modifier + save_prof_intelligence * proficiency_bonus,
    wisdom_saving_throw_modifier = wisdom_modifier + save_prof_wisdom * proficiency_bonus,
    charisma_saving_throw_modifier = charisma_modifier + save_prof_charisma * proficiency_bonus,

    skill_acrobatics_modifier = dexterity_modifier + skill_acrobatics * proficiency_bonus,
    skill_animal_handling_modifier = wisdom_modifier + skill_animal_handling * proficiency_bonus,
    skill_arcana_modifier = intelligence_modifier + skill_arcana * proficiency_bonus,
    skill_athletics_modifier = strength_modifier + skill_athletics * proficiency_bonus,
    skill_deception_modifier = charisma_modifier + skill_deception * proficiency_bonus,
    skill_history_modifier = intelligence_modifier + skill_history * proficiency_bonus,
    skill_insight_modifier = wisdom_modifier + skill_insight * proficiency_bonus,
    skill_intimidation_modifier = charisma_modifier + skill_intimidation * proficiency_bonus,
    skill_investigation_modifier = intelligence_modifier + skill_investigation * proficiency_bonus,
    skill_medicine_modifier = wisdom_modifier + skill_medicine * proficiency_bonus,
    skill_nature_modifier = intelligence_modifier + skill_nature * proficiency_bonus,
    skill_perception_modifier = wisdom_modifier + skill_perception * proficiency_bonus,
    skill_performance_modifier = charisma_modifier + skill_performance * proficiency_bonus,
    skill_persuasion_modifier = charisma_modifier + skill_persuasion * proficiency_bonus,
    skill_religion_modifier = intelligence_modifier + skill_religion * proficiency_bonus,
    skill_sleight_of_hand_modifier = dexterity_modifier + skill_sleight_of_hand * proficiency_bonus,
    skill_stealth_modifier = dexterity_modifier + skill_stealth * proficiency_bonus,
    skill_survival_modifier = wisdom_modifier + skill_survival * proficiency_bonus,

    passive_perception = 10 + skill_perception_modifier,
    passive_investigation = 10 + skill_investigation_modifier,
    passive_insight = 10 + skill_insight_modifier,

    initiative = dexterity_modifier

FROM characters
WHERE id = $1;

-- name: ListCharacters :many
SELECT * FROM characters
ORDER BY name;

-- name: ListUserCharacters :many
SELECT * FROM characters
WHERE owner_id = $1
ORDER BY name;

-- name: ListPlayerCharacters :many
SELECT * FROM characters
WHERE is_npc = FALSE
ORDER BY name;

-- name: ListNPCs :many
SELECT * FROM characters
WHERE is_npc = TRUE
ORDER BY name;

-- name: UpdateCharacterInfo :one
UPDATE characters
SET
    name                   = COALESCE(sqlc.narg('name'), name),
    race                   = COALESCE(sqlc.narg('race'), race),
    background             = COALESCE(sqlc.narg('background'), background),
    alignment              = COALESCE(sqlc.narg('alignment'), alignment),
    inspiration            = COALESCE(sqlc.narg('inspiration'), inspiration),
    speed                  = COALESCE(sqlc.narg('speed'), speed),
    ideals                 = COALESCE(sqlc.narg('ideals'), ideals),
    bonds                  = COALESCE(sqlc.narg('bonds'), bonds),
    flaws                  = COALESCE(sqlc.narg('flaws'), flaws),
    notes                  = COALESCE(sqlc.narg('notes'), notes),
    updated_at             = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateCharacterAbilityScores :one
UPDATE characters
SET
    strength               = COALESCE(sqlc.narg('strength'), strength),
    dexterity              = COALESCE(sqlc.narg('dexterity'), dexterity),
    constitution           = COALESCE(sqlc.narg('constitution'), constitution),
    intelligence           = COALESCE(sqlc.narg('intelligence'), intelligence),
    wisdom                 = COALESCE(sqlc.narg('wisdom'), wisdom),
    charisma               = COALESCE(sqlc.narg('charisma'), charisma),
    save_prof_strength     = COALESCE(sqlc.narg('save_prof_strength'), save_prof_strength),
    save_prof_dexterity    = COALESCE(sqlc.narg('save_prof_dexterity'), save_prof_dexterity),
    save_prof_constitution = COALESCE(sqlc.narg('save_prof_constitution'), save_prof_constitution),
    save_prof_intelligence = COALESCE(sqlc.narg('save_prof_intelligence'), save_prof_intelligence),
    save_prof_wisdom       = COALESCE(sqlc.narg('save_prof_wisdom'), save_prof_wisdom),
    save_prof_charisma     = COALESCE(sqlc.narg('save_prof_charisma'), save_prof_charisma),
    updated_at             = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateCharacterSkills :one
UPDATE characters
SET
    skill_acrobatics       = COALESCE(sqlc.narg('skill_acrobatics'), skill_acrobatics),
    skill_animal_handling  = COALESCE(sqlc.narg('skill_animal_handling'), skill_animal_handling),
    skill_arcana           = COALESCE(sqlc.narg('skill_arcana'), skill_arcana),
    skill_athletics        = COALESCE(sqlc.narg('skill_athletics'), skill_athletics),
    skill_deception        = COALESCE(sqlc.narg('skill_deception'), skill_deception),
    skill_history          = COALESCE(sqlc.narg('skill_history'), skill_history),
    skill_insight          = COALESCE(sqlc.narg('skill_insight'), skill_insight),
    skill_intimidation     = COALESCE(sqlc.narg('skill_intimidation'), skill_intimidation),
    skill_investigation    = COALESCE(sqlc.narg('skill_investigation'), skill_investigation),
    skill_medicine         = COALESCE(sqlc.narg('skill_medicine'), skill_medicine),
    skill_nature           = COALESCE(sqlc.narg('skill_nature'), skill_nature),
    skill_perception       = COALESCE(sqlc.narg('skill_perception'), skill_perception),
    skill_performance      = COALESCE(sqlc.narg('skill_performance'), skill_performance),
    skill_persuasion       = COALESCE(sqlc.narg('skill_persuasion'), skill_persuasion),
    skill_religion         = COALESCE(sqlc.narg('skill_religion'), skill_religion),
    skill_sleight_of_hand  = COALESCE(sqlc.narg('skill_sleight_of_hand'), skill_sleight_of_hand),
    skill_stealth          = COALESCE(sqlc.narg('skill_stealth'), skill_stealth),
    skill_survival         = COALESCE(sqlc.narg('skill_survival'), skill_survival),
    updated_at             = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateCharacterLevel :one
UPDATE characters
SET
    class                  = COALESCE(sqlc.narg('class'), class),
    level                  = COALESCE(sqlc.narg('level'), level),
    xp                     = COALESCE(sqlc.narg('xp'), xp),
    max_hp                 = COALESCE(sqlc.narg('max_hp'), max_hp),
    hit_dice_type          = COALESCE(sqlc.narg('hit_dice_type'), hit_dice_type),
    hit_dice_remaining     = COALESCE(sqlc.narg('hit_dice_remaining'), hit_dice_remaining),
    updated_at             = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateCharacterTraining :one
UPDATE characters
SET
    armor_class            = COALESCE(sqlc.narg('armor_class'), armor_class),
    attunement_slots       = COALESCE(sqlc.narg('attunement_slots'), attunement_slots),
    training_armor         = COALESCE(sqlc.narg('training_armor'), training_armor),
    training_weapons       = COALESCE(sqlc.narg('training_weapons'), training_weapons),
    training_tools         = COALESCE(sqlc.narg('training_tools'), training_tools),
    training_languages     = COALESCE(sqlc.narg('training_languages'), training_languages),
    spellcasting_ability   = COALESCE(sqlc.narg('spellcasting_ability'), spellcasting_ability),
    resistances            = COALESCE(sqlc.narg('resistances'), resistances),
    vulnerabilities        = COALESCE(sqlc.narg('vulnerabilities'), vulnerabilities),
    immunities             = COALESCE(sqlc.narg('immunities'), immunities),
    updated_at             = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateCharacterCurrency :one
UPDATE characters
SET
    copper                 = COALESCE(sqlc.narg('copper'), copper),
    silver                 = COALESCE(sqlc.narg('silver'), silver),
    electrum               = COALESCE(sqlc.narg('electrum'), electrum),
    gold                   = COALESCE(sqlc.narg('gold'), gold),
    platinum               = COALESCE(sqlc.narg('platinum'), platinum),
    conditions             = COALESCE(sqlc.narg('conditions'), conditions),
    updated_at             = NOW()
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: UpdateCharacterHP :one
UPDATE characters
SET
    current_hp = $2,
    temp_hp    = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateDeathSaves :one
UPDATE characters
SET
    death_save_successes    = $2,
    death_save_failures     = $3,
    updated_at              = NOW()
WHERE id = $1
RETURNING *;

-- name: ResetDeathSaves :one
UPDATE characters
SET
    death_save_successes    = 0,
    death_save_failures     = 0,
    updated_at              = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateConditions :one
UPDATE characters
SET
    conditions = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: LongRest :one
UPDATE characters
SET
    current_hp              = max_hp,
    temp_hp                 = 0,
    hit_dice_remaining      = GREATEST(hit_dice_remaining + (level / 2), level),
    death_save_successes    = 0,
    death_save_failures     = 0,
    conditions              = '{}',
    updated_at              = NOW()
WHERE id = $1
RETURNING *;

-- name: ShortRest :one
UPDATE characters
SET
    hit_dice_remaining      = GREATEST(hit_dice_remaining - $2, 0),
    current_hp              = LEAST(current_hp + $3, max_hp),
    temp_hp                 = 0,
    updated_at              = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCharacter :exec
DELETE FROM characters
WHERE id = $1;
