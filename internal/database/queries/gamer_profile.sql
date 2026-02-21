-- name: GetStudent :one
SELECT id, student_number, first_name, last_name, membership_tier,
       banned, notes, created_at, membership_expiry_date
FROM gamer_profile
WHERE student_number = $1;
