package http

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/yincongcyincong/MuseBot/db"
	"github.com/yincongcyincong/MuseBot/logger"
)

type GoogleAuthRequest struct {
	Credential string `json:"credential"`
}

type GoogleAuthResponse struct {
	Status string `json:"status"`
	UserId string `json:"user_id"`
	Token  string `json:"token"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

// GoogleAuthHandler handles the Google Sign-In verification
func GoogleAuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GoogleAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Basic parsing of the JWT credential (header.payload.signature)
	parts := strings.Split(req.Credential, ".")
	if len(parts) != 3 {
		http.Error(w, "Invalid credential format", http.StatusBadRequest)
		return
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		// Fallback to padded decoding if RawURLEncoding fails
		payloadSegment := parts[1]
		if l := len(payloadSegment) % 4; l > 0 {
			payloadSegment += strings.Repeat("=", 4-l)
		}
		payloadBytes, err = base64.URLEncoding.DecodeString(payloadSegment)
	}

	if err != nil {
		logger.Error("Failed to decode JWT payload", "err", err)
		http.Error(w, "Invalid credential encoding", http.StatusBadRequest)
		return
	}

	var payload struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		logger.Error("Failed to parse JWT payload", "err", err)
		http.Error(w, "Invalid credential payload", http.StatusBadRequest)
		return
	}

	if payload.Sub == "" {
		http.Error(w, "Invalid credential, no subject", http.StatusBadRequest)
		return
	}

	// Use Google's sub as the unique user ID. Prefix to avoid collisions.
	userId := "google_" + payload.Sub

	// Ensure user exists in our DB
	_, err = db.InsertUser(userId, "{}")
	if err != nil {
		logger.Error("Failed to create/get user during auth", "err", err, "userId", userId)
	}

	// Check if user has a linked profile (completed onboarding + bot link)
	profile, _ := db.GetProfileByGoogleId(userId)
	hasProfile := profile != nil

	// Generate a simple session token (in a real app, this should be a signed JWT)
	sessionToken := uuid.New().String()

	resp := map[string]interface{}{
		"status":      "success",
		"user_id":     userId,
		"token":       sessionToken,
		"name":        payload.Name,
		"email":       payload.Email,
		"has_profile": hasProfile,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
