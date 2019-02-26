-- +migrate Up
ALTER TABLE snap_rules ADD COLUMN volume_organization varchar(255) NOT NULL DEFAULT '';

-- +migrate Down
ALTER TABLE snap_rules DROP COLUMN volume_organization;
