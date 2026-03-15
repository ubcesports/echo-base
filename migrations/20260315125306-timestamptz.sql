-- +migrate Up
ALTER TABLE gamer_profile
    ALTER COLUMN created_at TYPE TIMESTAMPTZ
    USING (created_at::timestamp AT TIME ZONE 'UTC'),
    ALTER COLUMN membership_expiry_date TYPE TIMESTAMPTZ
    USING (membership_expiry_date::timestamp AT TIME ZONE 'UTC');

ALTER TABLE gamer_activity
    ALTER COLUMN started_at TYPE TIMESTAMPTZ
    USING (started_at AT TIME ZONE 'UTC'),
    ALTER COLUMN ended_at TYPE TIMESTAMPTZ
    USING (ended_at AT TIME ZONE 'UTC');

ALTER TABLE application
    ALTER COLUMN created_at TYPE TIMESTAMPTZ
    USING (created_at AT TIME ZONE 'UTC'),
    ALTER COLUMN last_used_at TYPE TIMESTAMPTZ
    USING (last_used_at AT TIME ZONE 'UTC');

-- +migrate Down
ALTER TABLE application
    ALTER COLUMN created_at TYPE TIMESTAMP
    USING (created_at AT TIME ZONE 'UTC'),
    ALTER COLUMN last_used_at TYPE TIMESTAMP
    USING (last_used_at AT TIME ZONE 'UTC');

ALTER TABLE gamer_activity
    ALTER COLUMN started_at TYPE TIMESTAMP
    USING (started_at AT TIME ZONE 'UTC'),
    ALTER COLUMN ended_at TYPE TIMESTAMP
    USING (ended_at AT TIME ZONE 'UTC');

ALTER TABLE gamer_profile
    ALTER COLUMN created_at TYPE DATE
    USING ((created_at AT TIME ZONE 'UTC')::date),
    ALTER COLUMN membership_expiry_date TYPE DATE
    USING ((membership_expiry_date AT TIME ZONE 'UTC')::date);
