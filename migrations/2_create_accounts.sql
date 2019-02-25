-- +migrate Up
CREATE TABLE accounts (
    id serial PRIMARY KEY,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    name      varchar(255) NOT NULL,
    provider      varchar(30) NOT NULL,
    user_id   integer  NOT NULL
);

ALTER TABLE accounts ADD CONSTRAINT accounts_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE;

-- +migrate Down
DROP TABLE accounts;