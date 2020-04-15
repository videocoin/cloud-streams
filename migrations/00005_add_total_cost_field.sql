-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE streams ADD `total_cost` DECIMAL(10,6) DEFAULT 0;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE streams DROP `total_cost`;
