-- +goose Up
ALTER TABLE inventory_items
    ALTER COLUMN weight TYPE INTEGER USING weight::INTEGER,
    ADD COLUMN value INTEGER NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE inventory_items
    DROP COLUMN value,
    ALTER COLUMN weight TYPE NUMERIC(8,2) USING weight::NUMERIC(8,2);
