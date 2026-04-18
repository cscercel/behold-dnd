-- +goose Up
ALTER TABLE characters
    ADD COLUMN spellcasting_ability TEXT NOT NULL DEFAULT '';

AlTER TABLE features
    ADD COLUMN action_type TEXT NOT NULL DEFAULT 'none'
    CHECK (action_type IN ('none', 'action', 'bonus_action', 'reaction', 'free'));

-- +goose Down
ALTER TABLE features
    DROP COLUMN action_type;

ALTER TABLE characters
    DROP COLUMN spellcasting_ability;
