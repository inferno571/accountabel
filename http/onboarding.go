package http

import (
	"encoding/json"
	"net/http"

	"github.com/yincongcyincong/MuseBot/db"
	"github.com/yincongcyincong/MuseBot/logger"
	"github.com/yincongcyincong/MuseBot/param"
	"github.com/yincongcyincong/MuseBot/utils"
)

type OnboardingRequest struct {
	Timezone     string `json:"timezone"`
	Time         string `json:"time"`
	Addictions   string `json:"addictions"`
	Platform     string `json:"platform"`
	GoogleUserId string `json:"google_user_id"`
	GoogleName   string `json:"google_name"`
}

// OnboardingStartHandler handles creating a pending onboarding configuration
func OnboardingStartHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req OnboardingRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.ErrorCtx(ctx, "failed to decode onboarding request", "err", err)
		utils.Failure(ctx, w, r, param.CodeParamError, "Invalid JSON body", err)
		return
	}

	if req.Time == "" || req.Platform == "" {
		utils.Failure(ctx, w, r, param.CodeParamError, "Missing required fields (time/platform)", nil)
		return
	}

	code, err := db.CreatePendingProfile(req.Timezone, req.Time, req.Addictions, req.Platform, req.GoogleUserId, req.GoogleName)
	if err != nil {
		logger.ErrorCtx(ctx, "failed to create pending profile", "err", err)
		utils.Failure(ctx, w, r, param.CodeDBQueryFail, "Database error", err)
		return
	}

	utils.Success(ctx, w, r, map[string]string{
		"code": code,
	})
}
