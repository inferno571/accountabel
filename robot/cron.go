package robot

import (
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
	"github.com/yincongcyincong/MuseBot/db"
	"github.com/yincongcyincong/MuseBot/logger"
	"github.com/yincongcyincong/MuseBot/param"
	"github.com/yincongcyincong/MuseBot/utils"
)

func InitCron() {
	time.Sleep(10 * time.Second)
	cronJobs, err := db.GetCronsByPage(1, 10000, "", "")
	if err != nil {
		logger.Error("get crons error", "err", err)
		return
	}

	Cron = cron.New(cron.WithSeconds())
	for _, c := range cronJobs {
		if c.CronSpec != "" && c.Status == 1 && c.Type != "" && c.Prompt != "" {
			cronID, err := Cron.AddFunc(c.CronSpec, func() {
				ExecCron(&c)
			})
			if err != nil {
				logger.Error("crontab parse error", "err", err)
				continue
			}

			err = db.UpdateCronJobId(c.ID, int(cronID))
			if err != nil {
				logger.Error("update cron job id error", "err", err)
			}
		}
	}

	// Restore scheduled accountability check-in cron jobs from user_profiles
	initAccountabilityCrons()

	Cron.Start()
}

func ExecCron(c *db.Cron) {
	logger.Info("exec cron", "cron", c.CronName, "cronSpec",
		c.CronSpec, "type", c.Type, "targetId", c.TargetID, "groupId", c.GroupID)
	switch c.Type {
	case param.Telegram:
		ExecTelegram(c)
	}
}

func ExecTelegram(c *db.Cron) {
	for _, targetId := range strings.Split(c.TargetID, ",") {
		targetId = strings.TrimSpace(targetId)
		if targetId == "" {
			continue
		}
		t := &TelegramRobot{
			Bot: TelegramBot,
			Update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					From: &tgbotapi.User{
						ID: int64(utils.ParseInt(targetId)),
					},
					Chat: &tgbotapi.Chat{
						ID: int64(utils.ParseInt(targetId)),
					},
					Text: c.Command + " " + c.Prompt,
				},
			},
		}
		t.Robot = NewRobot(WithRobot(t), WithSkipCheck(true), WithUseRecord(false))
		t.Robot.Exec()
	}

	for _, groupId := range strings.Split(c.GroupID, ",") {
		groupId = strings.TrimSpace(groupId)
		if groupId == "" {
			continue
		}
		t := &TelegramRobot{
			Bot: TelegramBot,
			Update: tgbotapi.Update{
				Message: &tgbotapi.Message{
					From: &tgbotapi.User{
						ID: int64(utils.ParseInt(c.CreateBy)),
					},
					Chat: &tgbotapi.Chat{
						ID: int64(utils.ParseInt(groupId)),
					},
					Text: c.Command + " " + c.Prompt,
				},
			},
		}
		t.Robot = NewRobot(WithRobot(t), WithSkipCheck(true), WithUseRecord(false))
		t.Robot.Exec()
	}

}

func AddCron(cronInfo *db.Cron) error {
	cronJobId, err := Cron.AddFunc(cronInfo.CronSpec, func() {
		ExecCron(cronInfo)
	})
	if err != nil {
		return err
	}

	err = db.UpdateCronJobId(cronInfo.ID, int(cronJobId))
	if err != nil {
		return err
	}

	return nil

}
