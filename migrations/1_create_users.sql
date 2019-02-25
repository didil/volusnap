-- +migrate Up
CREATE TABLE users (
    id serial PRIMARY KEY,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    email      varchar(255) NOT NULL,
    password   varchar(255) NOT NULL
);

CREATE UNIQUE INDEX users_email_idx ON users (email);

-- +migrate Down
DROP TABLE users;