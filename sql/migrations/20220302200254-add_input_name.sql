
-- +migrate Up
CREATE TABLE input (
    id BIGINT UNIQUE NOT NULL PRIMARY KEY,
    name VARCHAR UNIQUE NOT NULL,
    owner_id BIGINT NOT NULL,
    inputType VARCHAR NOT NULL,
    description VARCHAR NOT NULL,
    version BIGINT NOT NULL
);

-- +migrate Down
TRUNCATE TABLE input;
DROP TABLE input;

