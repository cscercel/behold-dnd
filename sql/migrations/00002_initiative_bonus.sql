-- +goose Up
ALTER TABLE combat_participants 
ADD COLUMN initiative_bonus INTEGER NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE combat_participants
DROP COLUMN initiative_bonus;
