package robot

import (
	"bytes"
	"context"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yincongcyincong/MuseBot/conf"
	"github.com/yincongcyincong/MuseBot/i18n"
	"github.com/yincongcyincong/MuseBot/logger"
	"github.com/yincongcyincong/MuseBot/metrics"
	"github.com/yincongcyincong/MuseBot/param"
	"github.com/yincongcyincong/MuseBot/utils"
)

var (
	DiscordSession *discordgo.Session
)

type DiscordRobot struct {
	Session *discordgo.Session
	Msg     *discordgo.MessageCreate
	Inter   *discordgo.InteractionCreate

	Robot        *RobotInfo
	Prompt       string
	Command      string
	ImageContent []byte
	AudioContent []byte
	UserName     string
}

func StartDiscordRobot(ctx context.Context) {
	var err error
	DiscordSession, err = discordgo.New("Bot " + conf.BaseConfInfo.DiscordBotToken)
	if err != nil {
		logger.ErrorCtx(ctx, "create discord bot", "err", err)
		return
	}
	DiscordSession.Client = utils.GetRobotProxyClient()

	// 添加消息处理函数
	DiscordSession.AddHandler(messageCreate)
	DiscordSession.AddHandler(onInteractionCreate)

	// 打开连接
	err = DiscordSession.Open()
	if err != nil {
		logger.ErrorCtx(ctx, "connect fail", "err", err)
		return
	}

	logger.InfoCtx(ctx, "discordBot Info", "username", DiscordSession.State.User.Username)

	registerSlashCommands(DiscordSession)

	select {
	case <-ctx.Done():
		DiscordSession.Close()
	}
}

func NewDiscordRobot(s *discordgo.Session, msg *discordgo.MessageCreate, i *discordgo.InteractionCreate) *DiscordRobot {
	metrics.AppRequestCount.WithLabelValues("discord").Inc()
	dr := &DiscordRobot{
		Session: s,
		Msg:     msg,
		Inter:   i,
	}

	if msg != nil {
		dr.UserName = msg.Author.Username
	}

	if i != nil {
		dr.UserName = i.User.Username
	}

	return dr
}

func (d *DiscordRobot) checkValid() bool {
	chatId, msgId, _ := d.Robot.GetChatIdAndMsgIdAndUserID()

	if d.Msg != nil {
		if d.skipThisMsg() {
			logger.WarnCtx(d.Robot.Ctx, "skip this msg", "msgId", msgId, "chat", chatId, "content", d.Msg.Content)
			return false
		}
		d.Command, d.Prompt = ParseCommand(d.Msg.Content)
		if d.Session != nil && d.Session.State != nil && d.Session.State.User != nil {
			d.Command = strings.ReplaceAll(d.Command, "<@"+d.Session.State.User.ID+">", "")
		}
		d.getMessageContent()
		return true
	}

	if d.Inter != nil {
		switch d.Inter.Type {
		case discordgo.InteractionApplicationCommand, discordgo.InteractionApplicationCommandAutocomplete:
			d.Command = d.Inter.ApplicationCommandData().Name
		}

		if d.Inter != nil && d.Inter.Type == discordgo.InteractionApplicationCommand && len(d.Inter.ApplicationCommandData().Options) > 0 {
			opt := d.Inter.ApplicationCommandData().Options[0]
			switch opt.Type {
			case discordgo.ApplicationCommandOptionInteger:
				d.Prompt = strconv.FormatInt(opt.IntValue(), 10)
			default:
				d.Prompt = opt.StringValue()
			}
		}
		if d.Session != nil && d.Session.State != nil && d.Session.State.User != nil {
			d.Command = strings.ReplaceAll(d.Command, "<@"+d.Session.State.User.ID+">", "")
		}

		err := d.Session.InteractionRespond(d.Inter.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
		if err != nil {
			logger.ErrorCtx(d.Robot.Ctx, "Failed to defer interaction response", "err", err)
		}
		return true
	}

	return false
}

func (d *DiscordRobot) getMsgContent() string {
	if d.Msg != nil {
		return d.Msg.Content
	}
	return ""
}

func (d *DiscordRobot) getMessageContent() {
	var err error
	chatId, msgId, _ := d.Robot.GetChatIdAndMsgIdAndUserID()
	if d.Inter != nil && d.Inter.ApplicationCommandData().GetOption("image") != nil {
		if attachment, ok := d.Inter.ApplicationCommandData().GetOption("image").Value.(string); ok {
			d.ImageContent, err = utils.DownloadFile(d.Inter.ApplicationCommandData().Resolved.Attachments[attachment].URL)
			if err != nil {
				logger.WarnCtx(d.Robot.Ctx, "download image fail", "err", err)
			}
		}
	}

	if d.Msg != nil {
		attachments := d.Msg.Attachments
		if len(attachments) > 0 {
			for _, att := range attachments {
				if strings.HasPrefix(att.ContentType, "audio/") {
					d.AudioContent, err = utils.DownloadFile(att.URL)
					if d.AudioContent == nil || err != nil {
						logger.ErrorCtx(d.Robot.Ctx, "audio url empty", "url", att.URL, "err", err)
						d.Robot.SendMsg(chatId, err.Error(), msgId, "", nil)
						return
					}
					if d.AudioContent != nil {
						d.Prompt, err = d.Robot.GetAudioContent(d.AudioContent)
						if err != nil {
							logger.WarnCtx(d.Robot.Ctx, "get audio content err", "err", err)
							d.Robot.SendMsg(chatId, err.Error(), msgId, "", nil)
							return
						}
					}
				}

				if strings.HasPrefix(att.ContentType, "image/") {
					d.ImageContent, err = utils.DownloadFile(att.URL)
					if d.ImageContent == nil || err != nil {
						logger.ErrorCtx(d.Robot.Ctx, "image url empty", "url", att.URL, "err", err)
						d.Robot.SendMsg(chatId, err.Error(), msgId, "", nil)
						return
					}
				}
			}
		}
	}
}

func (d *DiscordRobot) requestLLM(content string) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.ErrorCtx(d.Robot.Ctx, "DiscordRobot panic", "err", r, "stack", string(debug.Stack()))
			}
		}()
		switch d.Command {
		case "talk":
			d.Talk()
			return
		}

		d.Robot.ExecCmd(d.Command, d.sendChatMessage, nil, nil)
	}()
}

func (d *DiscordRobot) executeLLM() {
	messageChan := &MsgChan{
		NormalMessageChan: make(chan *param.MsgInfo),
	}

	go d.Robot.ExecLLM(d.Prompt, messageChan)

	go d.Robot.HandleUpdate(messageChan, "mp3")
}

func (d *DiscordRobot) sendTextStream(messageChan *MsgChan) {

	var originalMsgID string
	var channelID string
	var err error

	if d.Msg != nil {
		channelID = d.Msg.ChannelID

		thinkingMsg, err := d.Session.ChannelMessageSend(channelID, i18n.GetMessage("thinking", nil))
		if err != nil {
			logger.WarnCtx(d.Robot.Ctx, "Sending thinking message failed", "err", err)
		} else {
			originalMsgID = thinkingMsg.ID
		}

	} else if d.Inter != nil {
		channelID = d.Inter.ChannelID

		err = d.Session.InteractionRespond(d.Inter.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
		if err != nil {
			logger.ErrorCtx(d.Robot.Ctx, "Failed to defer interaction response", "err", err)
		}
	} else {
		logger.ErrorCtx(d.Robot.Ctx, "Unknown Discord message type")
		return
	}

	var msg *param.MsgInfo
	for msg = range messageChan.NormalMessageChan {
		if len(msg.Content) == 0 {
			msg.Content = "get nothing from llm!"
		}

		if msg.MsgId == "" && originalMsgID != "" {
			msg.MsgId = originalMsgID
		}

		if d.Msg != nil {
			if msg.MsgId == "" && originalMsgID == "" {
				_, err = d.Session.ChannelMessageSend(channelID, msg.Content)
				if err != nil {
					logger.ErrorCtx(d.Robot.Ctx, "Sending message failed", "err", err)
				}
			} else {
				_, err = d.Session.ChannelMessageEdit(channelID, msg.MsgId, msg.Content)
				if err != nil {
					logger.ErrorCtx(d.Robot.Ctx, "Editing message failed", "msgID", msg.MsgId, "err", err)
				}
				originalMsgID = ""
			}
		} else if d.Inter != nil {
			if msg.MsgId == "" && originalMsgID == "" {
				_, err = d.Session.InteractionResponseEdit(d.Inter.Interaction, &discordgo.WebhookEdit{
					Content: &msg.Content,
				})
				if err != nil {
					logger.ErrorCtx(d.Robot.Ctx, "Sending interaction response failed", "err", err)
				}
			} else {
				_, err = d.Session.FollowupMessageCreate(d.Inter.Interaction, true, &discordgo.WebhookParams{
					Content: msg.Content,
				})
				if err != nil {
					logger.ErrorCtx(d.Robot.Ctx, "Editing followup interaction message failed", "err", err)
				}
				originalMsgID = ""
			}
		}
	}
}

func (d *DiscordRobot) skipThisMsg() bool {
	if d.Msg == nil || d.Msg.Author == nil ||
		d.Session == nil || d.Msg.Author.ID == d.Session.State.User.ID {
		return true
	}

	if d.Msg.GuildID == "" {
		if strings.TrimSpace(d.Msg.Content) == "" && len(d.Msg.Attachments) == 0 {
			return true
		}
		return false
	}

	mentionedBot := false
	for _, user := range d.Msg.Mentions {
		if user.ID == d.Session.State.User.ID {
			mentionedBot = true
			break
		}
	}

	if !mentionedBot {
		return true
	}

	contentWithoutMention := strings.TrimSpace(strings.ReplaceAll(d.Msg.Content, "<@"+d.Session.State.User.ID+">", ""))
	if contentWithoutMention == "" && len(d.Msg.Attachments) == 0 {
		return true
	}

	return false
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	d := NewDiscordRobot(s, m, nil)
	d.Robot = NewRobot(WithRobot(d))
	d.Robot.Exec()
}

func registerSlashCommands(s *discordgo.Session) {
	commands := []*discordgo.ApplicationCommand{
		{Name: param.Chat, Description: i18n.GetMessage("commands.chat.description", nil), Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "prompt", Description: "Prompt", Required: true},
			{Type: discordgo.ApplicationCommandOptionAttachment, Name: "image", Description: "upload a image", Required: false},
		}},
		{Name: param.TxtType, Description: i18n.GetMessage("commands.mode.description", nil), Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "type", Description: "Type", Required: false},
		}},
		{Name: param.PhotoType, Description: i18n.GetMessage("commands.mode.description", nil), Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "type", Description: "Type", Required: false},
		}},
		{Name: param.VideoType, Description: i18n.GetMessage("commands.mode.description", nil), Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "type", Description: "Type", Required: false},
		}},
		{Name: param.TxtModel, Description: i18n.GetMessage("commands.mode.description", nil), Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "type", Description: "Type", Required: false},
		}},
		{Name: param.PhotoModel, Description: i18n.GetMessage("commands.mode.description", nil), Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "type", Description: "Type", Required: false},
		}},
		{Name: param.VideoModel, Description: i18n.GetMessage("commands.mode.description", nil), Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "type", Description: "Type", Required: false},
		}},
		{Name: "talk", Description: i18n.GetMessage("commands.talk.description", nil)},
		{Name: param.State, Description: i18n.GetMessage("commands.state.description", nil)},
		{Name: param.Clear, Description: i18n.GetMessage("commands.clear.description", nil)},
		{Name: param.Retry, Description: i18n.GetMessage("commands.retry.description", nil)},
		{Name: param.Photo, Description: i18n.GetMessage("commands.photo.description", nil), Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "prompt", Description: "Prompt", Required: true},
		}},
		{Name: param.EditPhoto, Description: i18n.GetMessage("commands.photo.description", nil), Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "prompt", Description: "Prompt", Required: true},
			{Type: discordgo.ApplicationCommandOptionAttachment, Name: "image", Description: "upload a image", Required: false},
		}},
		{Name: param.Video, Description: i18n.GetMessage("commands.video.description", nil), Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "prompt", Description: "Prompt", Required: true},
		}},
		{Name: param.Help, Description: i18n.GetMessage("commands.help.description", nil)},
		{Name: param.Task, Description: i18n.GetMessage("commands.task.description", nil), Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "prompt", Description: "Prompt", Required: true},
		}},
		{Name: param.Mcp, Description: i18n.GetMessage("commands.mcp.description", nil), Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "prompt", Description: "Prompt", Required: true},
		}},
		{Name: param.CronDel, Description: i18n.GetMessage("commands.cron.description", nil), Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "id", Description: "id", Required: true},
		}},
		{Name: param.CronClear, Description: i18n.GetMessage("commands.cron.description", nil)},
		{Name: param.CronDel, Description: i18n.GetMessage("commands.cron.description", nil)},
		{Name: param.Setup, Description: "Set daily check-in time (e.g. /setup 20:30)", Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "time", Description: "Check-in time in HH:MM format", Required: true},
		}},
		{Name: param.Log, Description: "Log craving level 1-10", Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionInteger, Name: "level", Description: "Craving level (1-10)", Required: true},
		}},
		{Name: param.StreakCmd, Description: "View your current sobriety streak"},
		{Name: param.LinkCmd, Description: "Link onboarding web preferences", Options: []*discordgo.ApplicationCommandOption{
			{Type: discordgo.ApplicationCommandOptionString, Name: "code", Description: "6-digit pairing code", Required: true},
		}},
	}

	for _, cmd := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			logger.Error("Cannot create command", "cmd", cmd.Name, "err", err)
		}
	}
}

func onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("onInteractionCreate panic err", "err", err, "stack", string(debug.Stack()))
		}
	}()

	d := NewDiscordRobot(s, nil, i)
	d.Robot = NewRobot(WithRobot(d))
	d.Robot.Exec()
}

func (d *DiscordRobot) sendChatMessage() {
	d.Robot.TalkingPreCheck(func() {
		d.executeLLM()
	})
}

func (d *DiscordRobot) sendImg() {
	d.Robot.TalkingPreCheck(func() {
		chatId, msgId, _ := d.Robot.GetChatIdAndMsgIdAndUserID()

		prompt := strings.TrimSpace(d.getPrompt())
		if prompt == "" {
			d.Robot.SendMsg(chatId, i18n.GetMessage("video_empty_content", nil),
				msgId, tgbotapi.ModeMarkdown, nil)
			return
		}

		d.Robot.SendMsg(chatId, i18n.GetMessage("thinking", nil),
			msgId, tgbotapi.ModeMarkdown, nil)

		var lastImageContent = d.ImageContent
		var err error
		if len(lastImageContent) == 0 && strings.Contains(d.Command, "edit_photo") {
			lastImageContent, err = d.Robot.GetLastImageContent()
			if err != nil {
				logger.WarnCtx(d.Robot.Ctx, "get last image record fail", "err", err)
			}
		}

		imageContent, totalToken, err := d.Robot.CreatePhoto(prompt, lastImageContent)
		if err != nil {
			logger.WarnCtx(d.Robot.Ctx, "generate image fail", "err", err)
			d.Robot.SendMsg(chatId, err.Error(), msgId, "", nil)
			return
		}

		if err != nil {
			logger.WarnCtx(d.Robot.Ctx, "send image fail", "err", err)
			d.Robot.SendMsg(chatId, err.Error(), msgId, "", nil)
			return
		}

		d.Robot.saveRecord(imageContent, lastImageContent, param.ImageRecordType, totalToken)
	})
}

func (d *DiscordRobot) sendMedia(media []byte, contentType, sType string) error {
	chatId, msgId, _ := d.Robot.GetChatIdAndMsgIdAndUserID()
	var err error
	if sType == "image" {
		file := &discordgo.File{
			Name:   "image." + contentType,
			Reader: bytes.NewReader(media),
		}

		if d.Inter != nil {
			_, err = d.Session.InteractionResponseEdit(d.Inter.Interaction, &discordgo.WebhookEdit{
				Files: []*discordgo.File{file},
			})
			if err != nil {
				logger.ErrorCtx(d.Robot.Ctx, "Error sending message:", "err", err)
				return err
			}
		} else {
			messageSend := &discordgo.MessageSend{
				Reference: &discordgo.MessageReference{
					MessageID: msgId,
					ChannelID: chatId,
				},
				Files: []*discordgo.File{file},
			}
			_, err = d.Session.ChannelMessageSendComplex(chatId, messageSend)
			if err != nil {
				logger.ErrorCtx(d.Robot.Ctx, "Error sending message:", "err", err)
				return err
			}
		}
	} else {
		file := &discordgo.File{
			Name:   "video." + contentType,
			Reader: bytes.NewReader(media),
		}

		if d.Inter != nil {
			_, err = d.Session.InteractionResponseEdit(d.Inter.Interaction, &discordgo.WebhookEdit{
				Files: []*discordgo.File{file},
			})
			if err != nil {
				logger.ErrorCtx(d.Robot.Ctx, "Error sending message:", "err", err)
				return err
			}
		} else {
			messageSend := &discordgo.MessageSend{
				Reference: &discordgo.MessageReference{
					MessageID: msgId,
					ChannelID: chatId,
				},
				Files: []*discordgo.File{file},
			}
			_, err = d.Session.ChannelMessageSendComplex(chatId, messageSend)
			if err != nil {
				logger.ErrorCtx(d.Robot.Ctx, "Error sending message:", "err", err)
				return err
			}
		}
	}

	return nil
}

func (d *DiscordRobot) sendVideo() {
	d.Robot.TalkingPreCheck(func() {
		chatId, msgId, _ := d.Robot.GetChatIdAndMsgIdAndUserID()

		prompt := strings.TrimSpace(d.getPrompt())
		if prompt == "" {
			d.Robot.SendMsg(chatId, i18n.GetMessage("video_empty_content", nil),
				msgId, tgbotapi.ModeMarkdown, nil)
			return
		}

		d.Robot.SendMsg(chatId, i18n.GetMessage("thinking", nil),
			msgId, tgbotapi.ModeMarkdown, nil)

		var imageContent = d.ImageContent
		videoContent, totalToken, err := d.Robot.CreateVideo(prompt, imageContent)
		if err != nil {
			logger.WarnCtx(d.Robot.Ctx, "generate video fail", "err", err)
			d.Robot.SendMsg(chatId, err.Error(), msgId, "", nil)
			return
		}

		if err != nil {
			logger.WarnCtx(d.Robot.Ctx, "send video fail", "err", err)
			d.Robot.SendMsg(chatId, err.Error(), msgId, "", nil)
			return
		}

		d.Robot.saveRecord(videoContent, imageContent, param.VideoRecordType, totalToken)
	})
}

func (d *DiscordRobot) getPrompt() string {
	return d.Prompt
}

func (d *DiscordRobot) getPerMsgLen() int {
	return 1800
}

func (d *DiscordRobot) sendVoiceContent(voiceContent []byte, duration int) error {
	var err error
	if d.Msg != nil {
		_, err = d.Session.ChannelFileSend(d.Msg.ChannelID, "voice."+utils.DetectAudioFormat(voiceContent), bytes.NewReader(voiceContent))

	} else if d.Inter != nil {
		_, err = d.Session.InteractionResponseEdit(d.Inter.Interaction, &discordgo.WebhookEdit{
			Files: []*discordgo.File{
				{
					Name:   "voice." + utils.DetectAudioFormat(voiceContent),
					Reader: bytes.NewReader(voiceContent),
				},
			},
		})
	}

	return err
}

func (d *DiscordRobot) Talk() {
	chatId, msgId, _ := d.Robot.GetChatIdAndMsgIdAndUserID()
	d.Robot.SendMsg(chatId, "Voice chat feature is not supported in this build.", msgId, "", nil)
}

func (d *DiscordRobot) setCommand(command string) {
	d.Command = command
}

func (d *DiscordRobot) getCommand() string {
	return d.Command
}

func (d *DiscordRobot) getUserName() string {
	return d.UserName
}

func (d *DiscordRobot) setPrompt(prompt string) {
	d.Prompt = prompt
}

func (d *DiscordRobot) getAudio() []byte {
	return d.AudioContent
}

func (d *DiscordRobot) getImage() []byte {
	return d.ImageContent
}

func (d *DiscordRobot) setImage(image []byte) {
	d.ImageContent = image
}
