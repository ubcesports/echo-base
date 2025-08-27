-- +migrate Up
CREATE TABLE application
(
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  app_name TEXT NOT NULL,
  key_id VARCHAR(10) NOT NULL UNIQUE,
  hashed_key VARCHAR(75) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  last_used_at TIMESTAMP
);
-- +migrate Down
drop TABLE application;
