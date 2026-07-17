package db

import (
	"fmt"
	"time"
)

// --- Accountability Tracking Tables ---

// UserProfile stores user preferences for accountability check-ins
type UserProfile struct {
	ID                   int64  `json:"id"`
	UserID               string `json:"user_id"`
	Platform             string `json:"platform"` // telegram or discord
	Timezone             string `json:"timezone"`
	ScheduledCheckinTime string `json:"scheduled_checkin_time"` // e.g. "21:00"
	CronJobID            int64  `json:"cron_job_id"`            // reference to the cron scheduler entry
	CreateTime           int64  `json:"create_time"`
	UpdateTime           int64  `json:"update_time"`
	Addictions           string `json:"addictions"`
	LinkCode             string `json:"link_code"`
	LinkCodeExpires      int64  `json:"link_code_expires"`
}

// CheckIn logs daily metrics for a user
type CheckIn struct {
	ID            int64  `json:"id"`
	UserID        string `json:"user_id"`
	Timestamp     int64  `json:"timestamp"`
	CravingLevel  int    `json:"craving_level"` // 1-10
	RelapseStatus bool   `json:"relapse_status"`
	Notes         string `json:"notes"`
}

// Streak tracks current and highest sobriety streak per user
type Streak struct {
	ID            int64  `json:"id"`
	UserID        string `json:"user_id"`
	CurrentStreak int    `json:"current_streak"` // in days
	HighestStreak int    `json:"highest_streak"` // in days
	LastCheckIn   int64  `json:"last_check_in"`  // unix timestamp of last check-in
	UpdateTime    int64  `json:"update_time"`
}

// --- User Profiles ---

// UpsertUserProfile inserts or updates a user profile
func UpsertUserProfile(userID, platform, timezone, checkinTime string) error {
	now := time.Now().Unix()
	upsertSQL := `
		INSERT INTO user_profiles (user_id, platform, timezone, scheduled_checkin_time, cron_job_id, create_time, update_time, addictions, link_code, link_code_expires)
		VALUES (?, ?, ?, ?, 0, ?, ?, '', '', 0)
		ON CONFLICT(user_id) DO UPDATE SET
			platform = excluded.platform,
			timezone = excluded.timezone,
			scheduled_checkin_time = excluded.scheduled_checkin_time,
			update_time = excluded.update_time`

	_, err := DB.Exec(upsertSQL, userID, platform, timezone, checkinTime, now, now)
	return err
}

// GetUserProfile retrieves a user's profile
func GetUserProfile(userID string) (*UserProfile, error) {
	querySQL := `SELECT id, user_id, platform, timezone, scheduled_checkin_time, cron_job_id, create_time, update_time, addictions, link_code, link_code_expires
		FROM user_profiles WHERE user_id = ?`
	var p UserProfile
	err := DB.QueryRow(querySQL, userID).Scan(
		&p.ID, &p.UserID, &p.Platform, &p.Timezone, &p.ScheduledCheckinTime,
		&p.CronJobID, &p.CreateTime, &p.UpdateTime, &p.Addictions, &p.LinkCode, &p.LinkCodeExpires,
	)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("get user profile error: %w", err)
	}
	return &p, nil
}

// CreatePendingProfile creates a profile that is pending pairing via code
func CreatePendingProfile(timezone, scheduledCheckinTime, addictions, platform, googleUserID, googleName string) (string, error) {
	now := time.Now().Unix()
	code := fmt.Sprintf("%06d", (time.Now().UnixNano()/1000)%1000000)
	userID := fmt.Sprintf("pending:%s", code)
	expires := time.Now().Add(15 * time.Minute).Unix()

	insertSQL := `
		INSERT INTO user_profiles (user_id, platform, timezone, scheduled_checkin_time, cron_job_id, create_time, update_time, addictions, link_code, link_code_expires, google_user_id, google_name)
		VALUES (?, ?, ?, ?, 0, ?, ?, ?, ?, ?, ?, ?)`

	_, err := DB.Exec(insertSQL, userID, platform, timezone, scheduledCheckinTime, now, now, addictions, code, expires, googleUserID, googleName)
	if err != nil {
		return "", err
	}
	return code, nil
}

// LinkProfileByCode links a pending profile using a pairing code
func LinkProfileByCode(userID, platform, code string) (*UserProfile, error) {
	now := time.Now().Unix()
	querySQL := `SELECT id, user_id, platform, timezone, scheduled_checkin_time, cron_job_id, create_time, update_time, addictions, link_code, link_code_expires
		FROM user_profiles WHERE link_code = ? AND link_code_expires > ?`
	var p UserProfile
	err := DB.QueryRow(querySQL, code, now).Scan(
		&p.ID, &p.UserID, &p.Platform, &p.Timezone, &p.ScheduledCheckinTime,
		&p.CronJobID, &p.CreateTime, &p.UpdateTime, &p.Addictions, &p.LinkCode, &p.LinkCodeExpires,
	)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}

	// Delete any existing profile for this actual user to prevent UNIQUE conflict
	_, _ = DB.Exec("DELETE FROM user_profiles WHERE user_id = ?", userID)

	// Update the pending profile to link it to the actual user
	updateSQL := `UPDATE user_profiles SET user_id = ?, platform = ?, link_code = '', link_code_expires = 0, update_time = ? WHERE id = ?`
	_, err = DB.Exec(updateSQL, userID, platform, now, p.ID)
	if err != nil {
		return nil, err
	}

	p.UserID = userID
	p.Platform = platform
	p.LinkCode = ""
	p.LinkCodeExpires = 0
	p.UpdateTime = now

	return &p, nil
}

// UpdateUserProfileCronJobID updates the cron job ID tied to the user's scheduled check-in
func UpdateUserProfileCronJobID(userID string, cronJobID int64) error {
	updateSQL := `UPDATE user_profiles SET cron_job_id = ?, update_time = ? WHERE user_id = ?`
	_, err := DB.Exec(updateSQL, cronJobID, time.Now().Unix(), userID)
	return err
}

// --- Check-Ins ---

// InsertCheckIn logs a new check-in entry
func InsertCheckIn(userID string, cravingLevel int, relapseStatus bool, notes string) (int64, error) {
	insertSQL := `INSERT INTO check_ins (user_id, timestamp, craving_level, relapse_status, notes)
		VALUES (?, ?, ?, ?, ?)`
	result, err := DB.Exec(insertSQL, userID, time.Now().Unix(), cravingLevel, relapseStatus, notes)
	if err != nil {
		return 0, fmt.Errorf("insert check_in error: %w", err)
	}
	return result.LastInsertId()
}

// GetLatestCheckIn retrieves the most recent check-in for a user
func GetLatestCheckIn(userID string) (*CheckIn, error) {
	querySQL := `SELECT id, user_id, timestamp, craving_level, relapse_status, notes
		FROM check_ins WHERE user_id = ? ORDER BY timestamp DESC LIMIT 1`
	var c CheckIn
	err := DB.QueryRow(querySQL, userID).Scan(
		&c.ID, &c.UserID, &c.Timestamp, &c.CravingLevel, &c.RelapseStatus, &c.Notes,
	)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("get latest check_in error: %w", err)
	}
	return &c, nil
}

// --- Streaks ---

// GetOrCreateStreak retrieves the streak record for a user, creating one if it doesn't exist
func GetOrCreateStreak(userID string) (*Streak, error) {
	querySQL := `SELECT id, user_id, current_streak, highest_streak, last_check_in, update_time
		FROM streaks WHERE user_id = ?`
	var s Streak
	err := DB.QueryRow(querySQL, userID).Scan(
		&s.ID, &s.UserID, &s.CurrentStreak, &s.HighestStreak, &s.LastCheckIn, &s.UpdateTime,
	)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			// Create a new streak record
			now := time.Now().Unix()
			insertSQL := `INSERT INTO streaks (user_id, current_streak, highest_streak, last_check_in, update_time)
				VALUES (?, 0, 0, 0, ?)`
			result, insertErr := DB.Exec(insertSQL, userID, now)
			if insertErr != nil {
				return nil, fmt.Errorf("create streak error: %w", insertErr)
			}
			id, _ := result.LastInsertId()
			return &Streak{
				ID:            id,
				UserID:        userID,
				CurrentStreak: 0,
				HighestStreak: 0,
				LastCheckIn:   0,
				UpdateTime:    now,
			}, nil
		}
		return nil, fmt.Errorf("get streak error: %w", err)
	}
	return &s, nil
}

// UpdateStreak updates the streak record after a check-in
func UpdateStreak(userID string, relapseStatus bool) error {
	streak, err := GetOrCreateStreak(userID)
	if err != nil {
		return err
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	lastDay := time.Unix(streak.LastCheckIn, 0)
	lastDayStart := time.Date(lastDay.Year(), lastDay.Month(), lastDay.Day(), 0, 0, 0, 0, lastDay.Location()).Unix()

	newCurrent := streak.CurrentStreak
	newHighest := streak.HighestStreak

	if relapseStatus {
		// Relapse resets the streak
		newCurrent = 0
	} else {
		if streak.LastCheckIn == 0 {
			// First ever check-in
			newCurrent = 1
		} else if today-lastDayStart <= 86400 && today > lastDayStart {
			// Consecutive day
			newCurrent = streak.CurrentStreak + 1
		} else if today == lastDayStart {
			// Same day, don't increment
			newCurrent = streak.CurrentStreak
		} else {
			// Missed days, reset streak but count today
			newCurrent = 1
		}
	}

	if newCurrent > newHighest {
		newHighest = newCurrent
	}

	updateSQL := `UPDATE streaks SET current_streak = ?, highest_streak = ?, last_check_in = ?, update_time = ?
		WHERE user_id = ?`
	_, err = DB.Exec(updateSQL, newCurrent, newHighest, now.Unix(), now.Unix(), userID)
	return err
}

// GetAllScheduledProfiles retrieves all user profiles that have a scheduled check-in time
func GetAllScheduledProfiles() ([]UserProfile, error) {
	querySQL := `SELECT id, user_id, platform, timezone, scheduled_checkin_time, cron_job_id, create_time, update_time
		FROM user_profiles WHERE scheduled_checkin_time != ''`
	rows, err := DB.Query(querySQL)
	if err != nil {
		return nil, fmt.Errorf("get scheduled profiles error: %w", err)
	}
	defer rows.Close()

	var profiles []UserProfile
	for rows.Next() {
		var p UserProfile
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.Platform, &p.Timezone, &p.ScheduledCheckinTime,
			&p.CronJobID, &p.CreateTime, &p.UpdateTime,
		); err != nil {
			return nil, fmt.Errorf("scan user profile error: %w", err)
		}
		profiles = append(profiles, p)
	}
	return profiles, nil
}

// GetProfileByGoogleId retrieves a user's profile by their Google user ID
func GetProfileByGoogleId(googleUserID string) (*UserProfile, error) {
	querySQL := `SELECT id, user_id, platform, timezone, scheduled_checkin_time, cron_job_id, create_time, update_time, addictions, link_code, link_code_expires
		FROM user_profiles WHERE google_user_id = ? AND user_id NOT LIKE 'pending:%' LIMIT 1`
	var p UserProfile
	err := DB.QueryRow(querySQL, googleUserID).Scan(
		&p.ID, &p.UserID, &p.Platform, &p.Timezone, &p.ScheduledCheckinTime,
		&p.CronJobID, &p.CreateTime, &p.UpdateTime, &p.Addictions, &p.LinkCode, &p.LinkCodeExpires,
	)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("get profile by google id error: %w", err)
	}
	return &p, nil
}
