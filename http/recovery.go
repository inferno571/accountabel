package http

import (
	"encoding/json"
	"net/http"

	"github.com/yincongcyincong/MuseBot/db"
	"github.com/yincongcyincong/MuseBot/logger"
)

// In a real app, middleware would extract user ID from the session token.
// For this prototype, we'll accept it in headers or assume it's passed.
func getUserIDFromHeader(r *http.Request) string {
	return r.Header.Get("X-User-Id")
}

// LogCravingHandler handles saving a new craving log
func LogCravingHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromHeader(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Intensity      int    `json:"intensity"`
		TriggerContext string `json:"trigger_context"`
		Tags           string `json:"tags"`
		ActionTaken    string `json:"action_taken"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	_, err := db.InsertCravingLog(userID, req.Intensity, req.TriggerContext, req.Tags, req.ActionTaken)
	if err != nil {
		logger.Error("Failed to log craving", "err", err)
		http.Error(w, "Failed to save data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// GetCravingsHandler returns craving logs
func GetCravingsHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromHeader(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	logs, err := db.GetCravingLogs(userID)
	if err != nil {
		logger.Error("Failed to fetch cravings", "err", err)
		http.Error(w, "Failed to load data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"status": "success", "data": logs})
}

// LogSlipHandler handles saving a slip/relapse neutrally
func LogSlipHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromHeader(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Notes string `json:"notes"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	err := db.InsertSlip(userID, req.Notes)
	if err != nil {
		logger.Error("Failed to log slip", "err", err)
		http.Error(w, "Failed to save slip", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Slip logged. That's okay — every day is a new start."})
}

// AddCopingToolHandler saves a new distraction or tool
func AddCopingToolHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromHeader(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		ToolType string `json:"tool_type"`
		Content  string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	err := db.InsertCopingTool(userID, req.ToolType, req.Content)
	if err != nil {
		http.Error(w, "Failed to save tool", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// GetDashboardDataHandler aggregates streaks, craving counts, tools, check-ins, slips, and profile.
func GetDashboardDataHandler(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromHeader(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// If the user ID starts with "google_", look up the Telegram user ID from the profile
	actualUserID := userID
	var profile *db.UserProfile
	if len(userID) > 7 && userID[:7] == "google_" {
		p, _ := db.GetProfileByGoogleId(userID)
		if p != nil {
			actualUserID = p.UserID
			profile = p
		}
	}

	if profile == nil {
		profile, _ = db.GetUserProfile(actualUserID)
	}

	streakInfo, _ := db.GetStreakInfo(actualUserID)
	tools, _ := db.GetCopingTools(actualUserID)
	contacts, _ := db.GetEmergencyContacts(actualUserID)
	cravings, _ := db.GetCravingLogs(actualUserID)
	checkIns, _ := db.GetAllCheckIns(actualUserID)
	slips, _ := db.GetSlipLogs(actualUserID)

	response := map[string]interface{}{
		"status":    "success",
		"streak":    streakInfo,
		"tools":     tools,
		"contacts":  contacts,
		"cravings":  cravings,
		"check_ins": checkIns,
		"slips":     slips,
		"profile":   profile,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
