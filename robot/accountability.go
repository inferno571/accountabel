package robot

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	cron "github.com/robfig/cron/v3"
	"github.com/sashabaranov/go-openai"
	"github.com/yincongcyincong/MuseBot/db"
	"github.com/yincongcyincong/MuseBot/llm"
	"github.com/yincongcyincong/MuseBot/logger"
	"github.com/yincongcyincong/MuseBot/param"
	"github.com/yincongcyincong/MuseBot/utils"
)

var timeFormatRegex = regexp.MustCompile(`^([01]?\d|2[0-3]):([0-5]\d)$`)

// handleSetup handles the /setup <HH:MM> command.
// Saves the user's preferred check-in time and registers a cron job.
func (r *RobotInfo) handleSetup() {
	chatId, msgId, userId := r.GetChatIdAndMsgIdAndUserID()
	timeArg := strings.TrimSpace(r.Robot.getPrompt())

	if timeArg == "" || !timeFormatRegex.MatchString(timeArg) {
		r.SendMsg(chatId, "⚠️ Please provide a valid time in 24-hour format.\nUsage: `/setup 20:30`",
			msgId, tgbotapi.ModeMarkdown, nil)
		return
	}

	// Determine the platform
	platform := "telegram"
	if _, ok := r.Robot.(*DiscordRobot); ok {
		platform = "discord"
	}

	// Upsert user profile with the scheduled check-in time
	err := db.UpsertUserProfile(userId, platform, "UTC", timeArg)
	if err != nil {
		logger.ErrorCtx(r.Ctx, "upsert user profile fail", "err", err)
		r.SendMsg(chatId, "❌ Failed to save your check-in time. Please try again.",
			msgId, "", nil)
		return
	}

	// Parse HH:MM to build a cron spec (second minute hour * * *)
	parts := strings.Split(timeArg, ":")
	hour := parts[0]
	minute := parts[1]
	cronSpec := fmt.Sprintf("CRON_TZ=UTC 0 %s %s * * *", minute, hour)

	// Remove existing cron job for this user if one exists
	profile, err := db.GetUserProfile(userId)
	if err != nil {
		logger.ErrorCtx(r.Ctx, "get user profile fail", "err", err)
	}
	if profile != nil && profile.CronJobID > 0 {
		Cron.Remove(cronEntryID(profile.CronJobID))
	}

	// Register a new cron job for the proactive check-in
	cronJobID, err := Cron.AddFunc(cronSpec, func() {
		ExecCheckinCron(userId, platform)
	})
	if err != nil {
		logger.ErrorCtx(r.Ctx, "add checkin cron fail", "err", err)
		r.SendMsg(chatId, "❌ Failed to schedule your check-in. Please try again.",
			msgId, "", nil)
		return
	}

	// Save the cron job ID to the profile
	err = db.UpdateUserProfileCronJobID(userId, int64(cronJobID))
	if err != nil {
		logger.ErrorCtx(r.Ctx, "update cron job id fail", "err", err)
	}

	r.SendMsg(chatId, fmt.Sprintf("✅ Daily check-in scheduled at **%s** UTC.\nI'll message you every day at that time to log your status.", timeArg),
		msgId, tgbotapi.ModeMarkdown, nil)
}

// handleLog handles the /log <craving_level_1_10> command.
// Quick-logs a craving level without waiting for a prompt.
func (r *RobotInfo) handleLog() {
	chatId, msgId, userId := r.GetChatIdAndMsgIdAndUserID()
	levelStr := strings.TrimSpace(r.Robot.getPrompt())

	level, err := strconv.Atoi(levelStr)
	if err != nil || level < 1 || level > 10 {
		r.SendMsg(chatId, "⚠️ Please provide a craving level between 1 and 10.\nUsage: `/log 3`",
			msgId, tgbotapi.ModeMarkdown, nil)
		return
	}

	// Insert check-in with the craving level (no relapse, no notes for quick log)
	_, err = db.InsertCheckIn(userId, level, false, "")
	if err != nil {
		logger.ErrorCtx(r.Ctx, "insert check-in fail", "err", err)
		r.SendMsg(chatId, "❌ Failed to log your craving. Please try again.",
			msgId, "", nil)
		return
	}

	// Update streak (no relapse)
	err = db.UpdateStreak(userId, false)
	if err != nil {
		logger.ErrorCtx(r.Ctx, "update streak fail", "err", err)
	}

	// Get current streak for feedback
	streak, err := db.GetOrCreateStreak(userId)
	if err != nil {
		logger.ErrorCtx(r.Ctx, "get streak fail", "err", err)
	}

	emoji := "💪"
	if level <= 3 {
		emoji = "🌟"
	} else if level >= 7 {
		emoji = "🫂"
	}

	msg := fmt.Sprintf("%s Craving level **%d/10** logged.\n\n🔥 Current streak: **%d days**",
		emoji, level, streak.CurrentStreak)
	if streak.CurrentStreak > 0 && streak.CurrentStreak == streak.HighestStreak {
		msg += " *(personal best!)*"
	}

	r.SendMsg(chatId, msg, msgId, tgbotapi.ModeMarkdown, nil)
}

// handleStreak handles the /streak command.
// Returns the user's current consecutive clean days.
func (r *RobotInfo) handleStreak() {
	chatId, msgId, userId := r.GetChatIdAndMsgIdAndUserID()

	streak, err := db.GetOrCreateStreak(userId)
	if err != nil {
		logger.ErrorCtx(r.Ctx, "get streak fail", "err", err)
		r.SendMsg(chatId, "❌ Failed to retrieve your streak. Please try again.",
			msgId, "", nil)
		return
	}

	msg := fmt.Sprintf("📊 **Your Streak**\n\n🔥 Current streak: **%d days**\n🏆 Highest streak: **%d days**",
		streak.CurrentStreak, streak.HighestStreak)

	if streak.CurrentStreak >= 7 {
		msg += "\n\n🎉 Amazing progress! Keep going!"
	} else if streak.CurrentStreak >= 3 {
		msg += "\n\n💪 Great momentum! You've got this!"
	} else if streak.CurrentStreak > 0 {
		msg += "\n\n🌱 Every day counts. Stay strong!"
	} else {
		msg += "\n\n🌅 Today is a fresh start. You can do this!"
	}

	r.SendMsg(chatId, msg, msgId, tgbotapi.ModeMarkdown, nil)
}

// ExecCheckinCron is the proactive cron trigger that fires at the user's scheduled time.
// It retrieves the user's current streak and sends an LLM-generated check-in message.
func ExecCheckinCron(userID, platform string) {
	logger.Info("exec checkin cron", "userID", userID, "platform", platform)

	// Get the user's current streak
	streak, err := db.GetOrCreateStreak(userID)
	if err != nil {
		logger.Error("get streak for cron fail", "err", err, "userID", userID)
		return
	}

	// Get the user's profile to check if they have tracked habits/addictions
	profile, err := db.GetUserProfile(userID)
	if err != nil {
		logger.Error("get profile for cron fail", "err", err, "userID", userID)
	}

	addictionsPrompt := ""
	if profile != nil && profile.Addictions != "" {
		addictionsPrompt = fmt.Sprintf(" The user is seeking accountability for: %s.", profile.Addictions)
	}

	// Build the proactive check-in system prompt
	checkinSystemPrompt := fmt.Sprintf(
		"Proactively check in on the user for their scheduled daily log.%s "+
			"Acknowledge that they are currently on a day %d streak. "+
			"Keep it brief and ask them to log their current craving level from 1-10.",
		addictionsPrompt,
		streak.CurrentStreak,
	)

	switch platform {
	case param.Telegram:
		execTelegramCheckin(userID, checkinSystemPrompt)
	case param.Discord:
		execDiscordCheckin(userID, checkinSystemPrompt)
	}
}

// execDiscordCheckin sends a proactive check-in message via Discord DM using the LLM.
func execDiscordCheckin(userID, systemPrompt string) {
	if DiscordSession == nil {
		logger.Error("discord session not initialized for checkin cron")
		return
	}

	// Open a DM channel with the user
	dmChannel, err := DiscordSession.UserChannelCreate(userID)
	if err != nil {
		logger.Error("create discord DM channel fail", "err", err, "userID", userID)
		return
	}

	// Create a minimal LLM call with the check-in system prompt
	messageChan := make(chan *param.MsgInfo, 1)

	go func() {
		defer close(messageChan)

		l := llm.NewLLM(
			llm.WithChatId(userID),
			llm.WithUserId(userID),
			llm.WithContent("Daily check-in"),
			llm.WithMessageChan(messageChan),
			llm.WithContext(context.Background()),
		)

		// Inject the proactive check-in system prompt
		l.LLMClient.GetModel(l)
		l.LLMClient.GetMessage(openai.ChatMessageRoleSystem, systemPrompt)
		l.LLMClient.GetMessage(openai.ChatMessageRoleUser, "Daily check-in time. How am I doing?")

		content, err := l.LLMClient.SyncSend(l.Ctx, l)
		if err != nil {
			logger.Error("discord checkin cron LLM call fail", "err", err, "userID", userID)
			content = fmt.Sprintf("🔔 Daily check-in! You're on a streak — keep going! Use /log <1-10> to log your craving level.")
		}

		messageChan <- &param.MsgInfo{
			Content:  content,
			Finished: true,
		}
	}()

	// Read and send the LLM response as a DM
	for msg := range messageChan {
		if msg.Content != "" {
			_, err := DiscordSession.ChannelMessageSend(dmChannel.ID, msg.Content)
			if err != nil {
				logger.Error("send discord checkin cron message fail", "err", err, "userID", userID)
			}
		}
	}
}

// execTelegramCheckin sends a proactive check-in message via Telegram using the LLM.
func execTelegramCheckin(userID, systemPrompt string) {
	if TelegramBot == nil {
		logger.Error("telegram bot not initialized for checkin cron")
		return
	}

	chatID := int64(utils.ParseInt(userID))

	// Create a minimal LLM call with the check-in system prompt
	messageChan := make(chan *param.MsgInfo, 1)

	go func() {
		defer close(messageChan)

		l := llm.NewLLM(
			llm.WithChatId(userID),
			llm.WithUserId(userID),
			llm.WithContent("Daily check-in"),
			llm.WithMessageChan(messageChan),
			llm.WithContext(context.Background()),
		)

		// Inject the proactive check-in system prompt
		l.LLMClient.GetModel(l)
		l.LLMClient.GetMessage(openai.ChatMessageRoleSystem, systemPrompt)
		l.LLMClient.GetMessage(openai.ChatMessageRoleUser, "Daily check-in time. How am I doing?")

		content, err := l.LLMClient.SyncSend(l.Ctx, l)
		if err != nil {
			logger.Error("checkin cron LLM call fail", "err", err, "userID", userID)
			content = fmt.Sprintf("🔔 Daily check-in! You're on a streak — keep going! Use /log <1-10> to log your craving level.")
		}

		messageChan <- &param.MsgInfo{
			Content:  content,
			Finished: true,
		}
	}()

	// Read and send the LLM response
	for msg := range messageChan {
		if msg.Content != "" {
			tgMsg := tgbotapi.NewMessage(chatID, msg.Content)
			tgMsg.ParseMode = tgbotapi.ModeMarkdown
			_, err := TelegramBot.Send(tgMsg)
			if err != nil {
				logger.Error("send checkin cron message fail", "err", err, "userID", userID)
			}
		}
	}
}

// handleLink handles the /link <code> command.
// Retrieves the pending configuration and pairs it to this user ID.
func (r *RobotInfo) handleLink() {
	chatId, msgId, userId := r.GetChatIdAndMsgIdAndUserID()
	code := strings.TrimSpace(r.Robot.getPrompt())

	if code == "" {
		r.SendMsg(chatId, "⚠️ Please provide your 6-digit onboarding code.\nUsage: `/link 123456`",
			msgId, tgbotapi.ModeMarkdown, nil)
		return
	}

	platform := "telegram"
	if _, ok := r.Robot.(*DiscordRobot); ok {
		platform = "discord"
	}

	profile, err := db.LinkProfileByCode(userId, platform, code)
	if err != nil {
		logger.ErrorCtx(r.Ctx, "link profile by code fail", "err", err)
		r.SendMsg(chatId, "❌ Database error during link. Please try onboarding again.",
			msgId, "", nil)
		return
	}

	if profile == nil {
		r.SendMsg(chatId, "❌ Invalid or expired 6-digit code. Please check the code on the web page and try again.",
			msgId, "", nil)
		return
	}

	// Schedule the new daily check-in cron job
	parts := strings.Split(profile.ScheduledCheckinTime, ":")
	if len(parts) == 2 {
		tz := profile.Timezone
		if tz == "" || tz == "IST" || tz == "EST" || tz == "CST" || tz == "PST" || tz == "BST" || tz == "AEST" {
			// fallback map for old invalid values if they exist in DB
			timezoneMap := map[string]string{
				"EST": "America/New_York",
				"CST": "America/Chicago",
				"PST": "America/Los_Angeles",
				"IST": "Asia/Kolkata",
				"BST": "Europe/London",
				"AEST": "Australia/Sydney",
			}
			if mapped, ok := timezoneMap[tz]; ok {
				tz = mapped
			} else {
				tz = "UTC"
			}
		}
		cronSpec := fmt.Sprintf("CRON_TZ=%s 0 %s %s * * *", tz, parts[1], parts[0])

		// Remove existing cron job if one exists
		p, _ := db.GetUserProfile(userId)
		if p != nil && p.CronJobID > 0 {
			Cron.Remove(cronEntryID(p.CronJobID))
		}

		cronJobID, err := Cron.AddFunc(cronSpec, func() {
			ExecCheckinCron(userId, platform)
		})
		if err != nil {
			logger.ErrorCtx(r.Ctx, "add checkin cron fail", "err", err)
		} else {
			err = db.UpdateUserProfileCronJobID(userId, int64(cronJobID))
			if err != nil {
				logger.ErrorCtx(r.Ctx, "update cron job id fail", "err", err)
			}
		}
	}

	r.SendMsg(chatId, fmt.Sprintf("✅ **Onboarding Complete!**\n\n🎯 **Habits tracked:** %s\n⏰ **Daily check-in:** %s %s\n\nI will message you daily at that time to track your streak. You've got this! 💪",
		profile.Addictions, profile.ScheduledCheckinTime, profile.Timezone), msgId, tgbotapi.ModeMarkdown, nil)
}

// cronEntryID converts an int64 cron job ID to the cron.EntryID type.
func cronEntryID(id int64) cronEntryIDType {
	return cronEntryIDType(id)
}

// cronEntryIDType is a type alias to match the cron library's EntryID type.
type cronEntryIDType = cron.EntryID

// initAccountabilityCrons restores scheduled check-in cron jobs from user_profiles at startup.
func initAccountabilityCrons() {
	profiles, err := db.GetAllScheduledProfiles()
	if err != nil {
		logger.Error("get scheduled profiles error", "err", err)
		return
	}

	for _, p := range profiles {
		if p.ScheduledCheckinTime == "" {
			continue
		}

		parts := strings.Split(p.ScheduledCheckinTime, ":")
		if len(parts) != 2 {
			continue
		}

		tz := p.Timezone
		if tz == "" || tz == "IST" || tz == "EST" || tz == "CST" || tz == "PST" || tz == "BST" || tz == "AEST" {
			timezoneMap := map[string]string{
				"EST": "America/New_York",
				"CST": "America/Chicago",
				"PST": "America/Los_Angeles",
				"IST": "Asia/Kolkata",
				"BST": "Europe/London",
				"AEST": "Australia/Sydney",
			}
			if mapped, ok := timezoneMap[tz]; ok {
				tz = mapped
			} else {
				tz = "UTC"
			}
		}
		cronSpec := fmt.Sprintf("CRON_TZ=%s 0 %s %s * * *", tz, parts[1], parts[0])
		userID := p.UserID
		platform := p.Platform

		cronJobID, err := Cron.AddFunc(cronSpec, func() {
			ExecCheckinCron(userID, platform)
		})
		if err != nil {
			logger.Error("add accountability cron fail", "err", err, "userID", userID)
			continue
		}

		err = db.UpdateUserProfileCronJobID(userID, int64(cronJobID))
		if err != nil {
			logger.Error("update user profile cron job id fail", "err", err, "userID", userID)
		}

		logger.Info("restored accountability cron", "userID", userID, "time", p.ScheduledCheckinTime)
	}
}
