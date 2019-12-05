-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS `streams` (
  `id` char(36) NOT NULL,
  `user_id` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `profile_id` varchar(255) DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  `input_status` int(11) DEFAULT NULL,
  `stream_contract_id` bigint(20) unsigned DEFAULT NULL,
  `stream_contract_address` varchar(255) DEFAULT NULL,
  `input_url` varchar(255) DEFAULT NULL,
  `output_url` varchar(255) DEFAULT NULL,
  `refunded` tinyint(1) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `ready_at` timestamp NULL DEFAULT NULL,
  `completed_at` timestamp NULL DEFAULT NULL,
  `rtmp_url` varchar(255) DEFAULT NULL,
  `xxx_unrecognized` varbinary(255) DEFAULT NULL,
  `xxx_sizecache` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE streams;