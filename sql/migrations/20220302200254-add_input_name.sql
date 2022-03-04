
-- +migrate Up
CREATE TABLE input (
    ref VARCHAR UNIQUE NOT NULL PRIMARY KEY,
    name VARCHAR UNIQUE NOT NULL,
    owner VARCHAR NOT NULL,
    inputType VARCHAR NOT NULL,
    description VARCHAR NOT NULL
);

-- +migrate Down
TRUNCATE TABLE input;
DROP TABLE input;

