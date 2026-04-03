-- name: GetGamerActivity :many
SELECT id, student_number, pc_number, game, started_at, ended_at, exec_name
FROM gamer_activity
WHERE student_number = $1
ORDER BY started_at DESC;

-- name: GetTodayActivitiesByStudent :many
SELECT ga.id, ga.student_number, ga.pc_number, ga.game, ga.started_at, ga.ended_at, ga.exec_name
FROM gamer_activity ga
JOIN gamer_profile gp ON ga.student_number = gp.student_number
WHERE ga.student_number = $1
AND gp.membership_tier = 1
AND DATE(ga.started_at AT TIME ZONE 'America/Los_Angeles') = DATE(NOW() AT TIME ZONE 'America/Los_Angeles');

-- name: GetRecentActivities :many
SELECT ga.id, ga.student_number, ga.pc_number, ga.game, ga.started_at, ga.ended_at, ga.exec_name,
       gp.first_name, gp.last_name
FROM gamer_activity ga
JOIN gamer_profile gp ON ga.student_number = gp.student_number
ORDER BY ga.started_at DESC NULLS LAST
LIMIT $1 OFFSET $2;

-- name: GetRecentActivitiesWithSearch :many
SELECT ga.id, ga.student_number, ga.pc_number, ga.game, ga.started_at, ga.ended_at, ga.exec_name,
       gp.first_name, gp.last_name
FROM gamer_activity ga
JOIN gamer_profile gp ON ga.student_number = gp.student_number
WHERE ga.student_number::TEXT ILIKE '%' || $3 || '%'
   OR gp.first_name ILIKE '%' || $3 || '%'
   OR gp.last_name ILIKE '%' || $3 || '%'
   OR ga.game ILIKE '%' || $3 || '%'
   OR ga.exec_name ILIKE '%' || $3 || '%'
ORDER BY ga.started_at DESC NULLS LAST
LIMIT $1 OFFSET $2;

-- name: CreateGamerActivity :one
INSERT INTO gamer_activity (id, student_number, pc_number, game, started_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, student_number, pc_number, game, started_at, ended_at, exec_name;

-- name: UpdateActivityEndTime :one
UPDATE gamer_activity
SET ended_at = $1, exec_name = $2
WHERE student_number = $3
AND pc_number = $4
AND ended_at IS NULL
RETURNING id, student_number, pc_number, game, started_at, ended_at, exec_name;

-- name: GetActiveSessions :many
SELECT ga.id, ga.student_number, ga.pc_number, ga.game, ga.started_at, ga.ended_at, ga.exec_name,
       gp.first_name, gp.last_name, gp.membership_tier
FROM gamer_activity ga
JOIN gamer_profile gp ON ga.student_number = gp.student_number
WHERE ga.ended_at IS NULL;

-- name: GetExecLeaderboard :many
SELECT exec_name, COUNT(*)::BIGINT AS signout_count
FROM gamer_activity
WHERE ended_at IS NOT NULL
AND exec_name IS NOT NULL
AND ended_at >= $1
AND ended_at < $2
GROUP BY exec_name
ORDER BY signout_count DESC, exec_name ASC;
