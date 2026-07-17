package db

import (
	"time"
)

type CravingLog struct {
	ID             int64  `json:"id"`
	UserID         string `json:"user_id"`
	Intensity      int    `json:"intensity"`
	TriggerContext string `json:"trigger_context"` // Should be encrypted in practice
	Tags           string `json:"tags"`
	ActionTaken    string `json:"action_taken"`
	LoggedAt       int64  `json:"logged_at"`
}

type SlipLog struct {
	ID                 int64  `json:"id"`
	UserID             string `json:"user_id"`
	SlipDate           int64  `json:"slip_date"`
	PreviousStreakDays int    `json:"previous_streak_days"`
	Notes              string `json:"notes"` // Should be encrypted in practice
}

type EmergencyContact struct {
	ID           int64  `json:"id"`
	UserID       string `json:"user_id"`
	ContactName  string `json:"contact_name"`
	ContactPhone string `json:"contact_phone"`
	Relationship string `json:"relationship"`
}

type CopingTool struct {
	ID       int64  `json:"id"`
	UserID   string `json:"user_id"`
	ToolType string `json:"tool_type"`
	Content  string `json:"content"`
}

func InsertCravingLog(userID string, intensity int, trigger, tags, action string) (int64, error) {
	query := `INSERT INTO craving_logs (user_id, intensity, trigger_context, tags, action_taken, logged_at) VALUES (?, ?, ?, ?, ?, ?)`
	res, err := DB.Exec(query, userID, intensity, trigger, tags, action, time.Now().Unix())
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func GetCravingLogs(userID string) ([]CravingLog, error) {
	query := `SELECT id, user_id, intensity, trigger_context, tags, action_taken, logged_at FROM craving_logs WHERE user_id = ? ORDER BY logged_at DESC LIMIT 50`
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []CravingLog
	for rows.Next() {
		var l CravingLog
		if err := rows.Scan(&l.ID, &l.UserID, &l.Intensity, &l.TriggerContext, &l.Tags, &l.ActionTaken, &l.LoggedAt); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func InsertSlip(userID string, notes string) error {
	// First get current streak info
	queryStreak := `SELECT current_streak FROM streaks WHERE user_id = ?`
	var prevStreak int
	err := DB.QueryRow(queryStreak, userID).Scan(&prevStreak)
	if err != nil {
		prevStreak = 0 // If not found, defaults to 0
	}

	query := `INSERT INTO slips_log (user_id, slip_date, previous_streak_days, notes) VALUES (?, ?, ?, ?)`
	_, err = DB.Exec(query, userID, time.Now().Unix(), prevStreak, notes)
	if err != nil {
		return err
	}

	// Update the streak record quietly (reset current streak to 0, leaving highest_streak untouched)
	updateStreak := `UPDATE streaks SET current_streak = 0, update_time = ? WHERE user_id = ?`
	_, _ = DB.Exec(updateStreak, time.Now().Unix(), userID)

	return nil
}

func InsertCopingTool(userID string, toolType, content string) error {
	query := `INSERT INTO coping_tools (user_id, tool_type, content) VALUES (?, ?, ?)`
	_, err := DB.Exec(query, userID, toolType, content)
	return err
}

func GetCopingTools(userID string) ([]CopingTool, error) {
	query := `SELECT id, user_id, tool_type, content FROM coping_tools WHERE user_id = ?`
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tools []CopingTool
	for rows.Next() {
		var t CopingTool
		if err := rows.Scan(&t.ID, &t.UserID, &t.ToolType, &t.Content); err != nil {
			return nil, err
		}
		tools = append(tools, t)
	}
	return tools, nil
}

func InsertEmergencyContact(userID, name, phone, relationship string) error {
	query := `INSERT INTO emergency_contacts (user_id, contact_name, contact_phone, relationship) VALUES (?, ?, ?, ?)`
	_, err := DB.Exec(query, userID, name, phone, relationship)
	return err
}

func GetEmergencyContacts(userID string) ([]EmergencyContact, error) {
	query := `SELECT id, user_id, contact_name, contact_phone, relationship FROM emergency_contacts WHERE user_id = ?`
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var contacts []EmergencyContact
	for rows.Next() {
		var c EmergencyContact
		if err := rows.Scan(&c.ID, &c.UserID, &c.ContactName, &c.ContactPhone, &c.Relationship); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

// EnsureStreakExists makes sure there's a row in the streaks table for the user
func EnsureStreakExists(userID string) error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM streaks WHERE user_id = ?", userID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		_, err = DB.Exec("INSERT INTO streaks (user_id, current_streak, highest_streak, last_check_in, update_time) VALUES (?, 0, 0, 0, ?)", userID, time.Now().Unix())
		return err
	}
	return nil
}

func GetStreakInfo(userID string) (map[string]interface{}, error) {
	EnsureStreakExists(userID)
	var current, highest, last int64
	err := DB.QueryRow("SELECT current_streak, highest_streak, last_check_in FROM streaks WHERE user_id = ?", userID).Scan(&current, &highest, &last)
	if err != nil {
		return nil, err
	}
	
	// Increment current streak if last checkin was yesterday (simplified for demonstration)
	// A proper implementation checks timezone and actual days elapsed.
	
	return map[string]interface{}{
		"current_streak": current,
		"highest_streak": highest,
		"last_check_in": last,
	}, nil
}

// GetAllCheckIns returns all check-ins for a user (up to 200)
func GetAllCheckIns(userID string) ([]CheckIn, error) {
	query := `SELECT id, user_id, timestamp, craving_level, relapse_status, notes
		FROM check_ins WHERE user_id = ? ORDER BY timestamp DESC LIMIT 200`
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checkIns []CheckIn
	for rows.Next() {
		var c CheckIn
		if err := rows.Scan(&c.ID, &c.UserID, &c.Timestamp, &c.CravingLevel, &c.RelapseStatus, &c.Notes); err != nil {
			return nil, err
		}
		checkIns = append(checkIns, c)
	}
	return checkIns, nil
}

// GetSlipLogs returns all slip logs for a user
func GetSlipLogs(userID string) ([]SlipLog, error) {
	query := `SELECT id, user_id, slip_date, previous_streak_days, notes
		FROM slips_log WHERE user_id = ? ORDER BY slip_date DESC LIMIT 100`
	rows, err := DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slips []SlipLog
	for rows.Next() {
		var s SlipLog
		if err := rows.Scan(&s.ID, &s.UserID, &s.SlipDate, &s.PreviousStreakDays, &s.Notes); err != nil {
			return nil, err
		}
		slips = append(slips, s)
	}
	return slips, nil
}

