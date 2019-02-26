-- +migrate Up
CREATE TABLE snap_rules (
    id             serial PRIMARY KEY,
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    frequency      integer NOT NULL,
    volume_id      varchar(255) NOT NULL,
    volume_name    varchar(255) NOT NULL,
    volume_region  varchar(50) NOT NULL,
    account_id     integer NOT NULL
);

ALTER TABLE snap_rules ADD CONSTRAINT snap_rules_account_id_fkey FOREIGN KEY (account_id) REFERENCES accounts(id) ON UPDATE CASCADE ON DELETE CASCADE;

-- +migrate Down
DROP TABLE snap_rules;