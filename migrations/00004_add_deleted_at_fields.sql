-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE streams ADD `deleted_at` timestamp NULL DEFAULT NULL;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE streams DROP `deleted_at`;
