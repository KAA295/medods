
-- +migrate Up
CREATE TABLE tokens(
  id SERIAL PRIMARY KEY,
  token VARCHAR NOT NULL,
  user_id VARCHAR NOT NULL, 
  expiration_time TIMESTAMP 
);

-- +migrate Down
DROP TABLE IF EXISTS tokens;
