-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE streams 
    ADD `input_type`  varchar(255) DEFAULT 'INPUT_TYPE_RTMP', 
    ADD `output_type` varchar(255) DEFAULT 'OUTPUT_TYPE_HLS';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE streams DROP `input_type`, DROP `output_type`;
