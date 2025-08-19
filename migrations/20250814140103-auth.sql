-- +migrate Up
CREATE TABLE auth
(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  app_name VARCHAR(100) NOT NULL,
  key_id VARCHAR(10) NOT NULL UNIQUE,
  hashed_key VARCHAR(75) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  last_used_at TIMESTAMP
);
-- +migrate Down
drop TABLE auth;
