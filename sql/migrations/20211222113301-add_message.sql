
-- +migrate Up
CREATE TABLE message (
  id VARCHAR NOT NULL,
  hashtags VARCHAR ARRAY NOT NULL,
  channel_messages JSON NOT NULL,
  PRIMARY KEY (id)
);

-- +migrate Down
TRUNCATE TABLE message;
DROP TABLE message;
