-- +migrate Up
CREATE TABLE channel (
  id VARCHAR NOT NULL,
  hashtags VARCHAR ARRAY NOT NULL,
  PRIMARY KEY (id)
);
-- +migrate Down
TRUNCATE TABLE channel;
DROP TABLE channel;
