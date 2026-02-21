-- name: GetGamerProfile :one
SELECT *
FROM gamer_profile
WHERE student_number = $1;

-- name: UpsertGamerProfile :one
INSERT INTO gamer_profile (first_name, last_name, student_number, membership_tier, banned, notes, created_at, membership_expiry_date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (student_number)
DO UPDATE SET
    first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name,
    membership_tier = EXCLUDED.membership_tier,
    banned = EXCLUDED.banned,
    notes = EXCLUDED.notes,
    created_at = EXCLUDED.created_at,
    membership_expiry_date = EXCLUDED.membership_expiry_date
RETURNING *;

-- name: DeleteGamerProfile :execrows
DELETE FROM gamer_profile WHERE student_number = $1;

-- name: CheckMembershipValidity :one
SELECT membership_tier, membership_expiry_date
FROM gamer_profile
WHERE student_number = $1;
