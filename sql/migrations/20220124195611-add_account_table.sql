
-- +migrate Up
CREATE TABLE account (
  id VARCHAR NOT NULL PRIMARY KEY,
  telegram_id BIGINT UNIQUE NOT NULL,
  version BIGINT NOT NULL CHECK (version >= 0)
);

-- +migrate Down
TRUNCATE TABLE account;
DROP TABLE account;

