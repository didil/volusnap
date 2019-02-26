-- +migrate Up
CREATE TABLE snapshots (
    id                     serial PRIMARY KEY,
    created_at             timestamp with time zone,
    updated_at             timestamp with time zone,
    provider_snapshot_id   varchar(255) NOT NULL,
    snap_rule_id           integer NOT NULL
);

ALTER TABLE snapshots ADD CONSTRAINT snapshots_snap_rule_id_fkey FOREIGN KEY (snap_rule_id) REFERENCES snap_rules(id) ON UPDATE CASCADE ON DELETE CASCADE;

-- +migrate Down
DROP TABLE snapshots;