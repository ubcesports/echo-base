-- +migrate Up
ALTER TABLE gamer_profile ADD COLUMN membership_expiry_date DATE;
CREATE INDEX idx_gamer_profile_membership_expiry ON gamer_profile(membership_expiry_date);

-- +migrate Down
DROP INDEX idx_gamer_profile_membership_expiry;
ALTER TABLE gamer_profile DROP COLUMN membership_expiry_date;
