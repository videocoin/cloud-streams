-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE streams ADD `drm_xml` TEXT;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE streams DROP `drm_xml`;
