package repositories

// Builds query for getting gamer activity by student number
func BuildGamerActivityByStudentQuery() string {
	query := `
        SELECT *
        FROM gamer_activity
        WHERE student_number = $1`
	return query
}

// Builds query for getting today's Tier One member activity by student number
func BuildGamerActivityByTierOneStudentTodayQuery() string {
	query := `
		SELECT ga.*
		FROM gamer_activity ga
        JOIN gamer_profile gp 
        ON ga.student_number = gp.student_number
        WHERE ga.student_number = $1
        AND gp.membership_tier = 1
		AND DATE(ga.started_at) = DATE($2)`
	return query
}

// Builds query for recent gamer activity with optional search
func BuildGamerActivityRecentQuery(limit, offset int, search string) (string, []interface{}) {
	base := `
        SELECT ga.*, gp.first_name, gp.last_name
        FROM gamer_activity ga
        JOIN gamer_profile gp 
        ON ga.student_number = gp.student_number
    `
	args := []interface{}{limit, offset}

	if search != "" {
		base += `
            WHERE (ga.student_number ILIKE $3
            OR gp.first_name ILIKE $3
            OR gp.last_name ILIKE $3
            OR ga.game ILIKE $3
            OR ga.exec_name ILIKE $3
            OR TO_CHAR(ga.started_at, 'YYYY-MM-DD') ILIKE $3)
        `
		args = append(args, "%"+search+"%")
	}

	base += ` ORDER BY ga.started_at DESC NULLS LAST LIMIT $1 OFFSET $2`
	return base, args
}

// Builds query for inserting a new gamer activity
func BuildInsertGamerActivityQuery() string {
	query := `
        INSERT INTO gamer_activity 
        (student_number, pc_number, game, started_at)
        VALUES ($1, $2, $3, $4)
        RETURNING *
    `
	return query
}

// Builds query for updating gamer activity end time
func BuildUpdateGamerActivityQuery() string {
	query := `
        UPDATE gamer_activity
        SET ended_at = $1, exec_name = $4
        WHERE student_number = $2 AND pc_number = $3 AND ended_at IS NULL
        RETURNING *
    `
	return query
}

// Builds query for getting all active PCs
func BuildGetAllActivePCsQuery() string {
	query := `
        SELECT ga.*, gp.first_name, gp.last_name, gp.membership_tier,
            gp.banned, gp.notes, gp.created_at
        FROM gamer_activity ga
        JOIN gamer_profile gp 
        ON ga.student_number = gp.student_number
        WHERE ga.ended_at IS NULL
    `
	return query
}

func BuildCheckMemberQuery() string {
	query := `
		SELECT membership_expiry_date, membership_tier
		FROM gamer_profile
		WHERE student_number = $1
	`
	return query
}
