package conf

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/yincongcyincong/MuseBot/logger"
)

type BaseConf struct {
	StartTime int64 `json:"-"`
	ImageDay  int   `json:"-"`

	TelegramBotToken        string `json:"telegram_bot_token"`
	DiscordBotToken         string `json:"discord_bot_token"`
	SlackBotToken           string `json:"slack_bot_token"`
	SlackAppToken           string `json:"slack_app_token"`
	LarkAPPID               string `json:"lark_app_id"`
	LarkAppSecret           string `json:"lark_app_secret"`
	DingClientId            string `json:"ding_client_id"`
	DingClientSecret        string `json:"ding_client_secret"`
	ComWechatToken          string `json:"com_wechat_token"`
	ComWechatEncodingAESKey string `json:"com_wechat_encoding_aes_key"`
	ComWechatCorpID         string `json:"com_wechat_corp_id"`
	ComWechatSecret         string `json:"com_wechat_secret"`
	ComWechatAgentID        string `json:"com_wechat_agent_id"`
	WechatAppID             string `json:"wechat_app_id"`
	WechatAppSecret         string `json:"wechat_app_secret"`
	WechatToken             string `json:"wechat_token"`
	WechatEncodingAESKey    string `json:"wechat_encoding_aes_key"`
	WechatActive            bool   `json:"wechat_active"`
	QQAppID                 string `json:"qq_app_id"`
	QQAppSecret             string `json:"qq_app_secret"`
	QQOneBotReceiveToken    string `json:"qq_one_bot_receive_token"`
	QQOneBotSendToken       string `json:"qq_one_bot_send_token"`
	QQOneBotHttpServer      string `json:"qq_one_bot_http_server"`

	DeepseekToken     string `json:"deepseek_token"`
	OpenAIToken       string `json:"openai_token"`
	GeminiToken       string `json:"gemini_token"`
	OpenRouterToken   string `json:"open_router_token"`
	AI302Token        string `json:"ai_302_token"`
	VolToken          string `json:"vol_token"`
	AliyunToken       string `json:"aliyun_token"`
	ChatAnyWhereToken string `json:"chat_any_where_token"`
	ErnieAK           string `json:"ernie_ak"`
	ErnieSK           string `json:"ernie_sk"`

	BotName           string `json:"bot_name"`
	Type              string `json:"type"`
	MediaType         string `json:"media_type"`
	CustomUrl         string `json:"custom_url"`
	CustomPath        string `json:"custom_path"`
	VolcAK            string `json:"volc_ak"`
	VolcSK            string `json:"volc_sk"`
	DBType            string `json:"db_type"`
	DBConf            string `json:"db_conf"`
	LLMProxy          string `json:"llm_proxy"`
	RobotProxy        string `json:"robot_proxy"`
	Lang              string `json:"lang"`
	TokenPerUser      int    `json:"token_per_user"`
	MaxUserChat       int    `json:"max_user_chat"`
	HTTPHost          string `json:"http_host"`
	UseTools          bool   `json:"use_tools"`
	MaxQAPair         int    `json:"max_qa_pari"`
	Character         string `json:"character"`
	SmartMode         bool   `json:"smart_mode"`
	ContextExpireTime int    `json:"context_expire_time"`
	Powered           string `json:"powered"`
	SendMcpRes        bool   `json:"send_mcp_res"`
	SendMcpMediaToLLM bool   `json:"send_mcp_media_to_llm"`
	DefaultModel      string `json:"default_model"`
	LLMRetryTimes     int    `json:"llm_retry_times"`
	LLMRetryInterval  int    `json:"llm_retry_interval"`
	LLMOptionParam    bool   `json:"llm_option_param"`
	ImagePath         string `json:"image_path"`
	IsStreaming       bool   `json:"is_streaming"`

	CrtFile string `json:"crt_file"`
	KeyFile string `json:"key_file"`
	CaFile  string `json:"ca_file"`

	AllowedUserIds  map[string]bool `json:"allowed_user_ids"`
	AllowedGroupIds map[string]bool `json:"allowed_group_ids"`
}

var (
	BaseConfInfo = new(BaseConf)
	AllConf      = make(map[string]interface{})
)

func InitConf() {
	BaseConfInfo.StartTime = time.Now().Unix()
	if loadConf() {
		logConf("", "")
		return
	}

	flag.StringVar(&BaseConfInfo.TelegramBotToken, "telegram_bot_token", "", "Telegram bot tokens")
	flag.StringVar(&BaseConfInfo.DiscordBotToken, "discord_bot_token", "", "Discord bot tokens")
	flag.StringVar(&BaseConfInfo.SlackBotToken, "slack_bot_token", "", "Slack bot tokens")
	flag.StringVar(&BaseConfInfo.SlackAppToken, "slack_app_token", "", "Slack app tokens")
	flag.StringVar(&BaseConfInfo.LarkAPPID, "lark_app_id", "", "Lark app id")
	flag.StringVar(&BaseConfInfo.LarkAppSecret, "lark_app_secret", "", "Lark app secret")
	flag.StringVar(&BaseConfInfo.DingClientId, "ding_client_id", "", "Dingding client id")
	flag.StringVar(&BaseConfInfo.DingClientSecret, "ding_client_secret", "", "Dingding app secret")
	flag.StringVar(&BaseConfInfo.ComWechatToken, "com_wechat_token", "", "ComWechat token")
	flag.StringVar(&BaseConfInfo.ComWechatEncodingAESKey, "com_wechat_encoding_aes_key", "", "ComWechat encoding aes key")
	flag.StringVar(&BaseConfInfo.ComWechatCorpID, "com_wechat_corp_id", "", "ComWechat corp id")
	flag.StringVar(&BaseConfInfo.ComWechatSecret, "com_wechat_secret", "", "ComWechat secret")
	flag.StringVar(&BaseConfInfo.ComWechatAgentID, "com_wechat_agent_id", "", "ComWechat agent id")
	flag.StringVar(&BaseConfInfo.WechatAppID, "wechat_app_id", "", "Wechat app id")
	flag.StringVar(&BaseConfInfo.WechatAppSecret, "wechat_app_secret", "", "Wechat app secret")
	flag.StringVar(&BaseConfInfo.WechatEncodingAESKey, "wechat_encoding_aes_key", "", "Wechat encoding aes key")
	flag.StringVar(&BaseConfInfo.WechatToken, "wechat_token", "", "Wechat token")
	flag.BoolVar(&BaseConfInfo.WechatActive, "wechat_active", false, "Wechat active")
	flag.StringVar(&BaseConfInfo.QQAppID, "qq_app_id", "", "QQ app id")
	flag.StringVar(&BaseConfInfo.QQAppSecret, "qq_app_secret", "", "QQ app secret")
	flag.StringVar(&BaseConfInfo.QQOneBotReceiveToken, "qq_one_bot_receive_token", "Accountabel AI", "onebot receive token")
	flag.StringVar(&BaseConfInfo.QQOneBotSendToken, "qq_one_bot_send_token", "Accountabel AI", "onebot send token")
	flag.StringVar(&BaseConfInfo.QQOneBotHttpServer, "qq_one_bot_http_server", "http://127.0.0.1:3000", "onebot http server")
	flag.BoolVar(&BaseConfInfo.SmartMode, "smart_mode", false, "Smart mode")
	flag.IntVar(&BaseConfInfo.ContextExpireTime, "context_expire_time", 86400, "Context expire time")

	flag.StringVar(&BaseConfInfo.DeepseekToken, "deepseek_token", "", "deepseek auth token")
	flag.StringVar(&BaseConfInfo.OpenAIToken, "openai_token", "", "openai auth token")
	flag.StringVar(&BaseConfInfo.GeminiToken, "gemini_token", "", "gemini auth token")
	flag.StringVar(&BaseConfInfo.OpenRouterToken, "open_router_token", "", "openrouter auth token")
	flag.StringVar(&BaseConfInfo.AI302Token, "ai_302_token", "", "302.ai token")
	flag.StringVar(&BaseConfInfo.VolToken, "vol_token", "", "vol auth token")
	flag.StringVar(&BaseConfInfo.AliyunToken, "aliyun_token", "", "aliyun auth token")
	flag.StringVar(&BaseConfInfo.ErnieAK, "ernie_ak", "", "ernie ak")
	flag.StringVar(&BaseConfInfo.ErnieSK, "ernie_sk", "", "ernie sk")
	flag.StringVar(&BaseConfInfo.VolcAK, "volc_ak", "", "volc ak")
	flag.StringVar(&BaseConfInfo.VolcSK, "volc_sk", "", "volc sk")
	flag.StringVar(&BaseConfInfo.ChatAnyWhereToken, "chat_any_where_token", "", "chatAnyWhere Token")

	flag.StringVar(&BaseConfInfo.BotName, "bot_name", "Accountabel AI", "bot name")
	flag.StringVar(&BaseConfInfo.CustomUrl, "custom_url", "", "custom url")
	flag.StringVar(&BaseConfInfo.CustomPath, "custom_path", "", "custom path")
	flag.StringVar(&BaseConfInfo.Type, "type", "", "llm type: deepseek gemini openai openrouter vol chatanywhere")
	flag.StringVar(&BaseConfInfo.MediaType, "media_type", "", "media type: vol gemini openai aliyun 302-ai openrouter")
	flag.StringVar(&BaseConfInfo.DBType, "db_type", "sqlite3", "db type")
	flag.StringVar(&BaseConfInfo.DBConf, "db_conf", GetAbsPath("data/accountabel_bot.db"), "db conf")
	flag.StringVar(&BaseConfInfo.LLMProxy, "llm_proxy", "", "llm proxy: http://127.0.0.1:7890")
	flag.StringVar(&BaseConfInfo.RobotProxy, "robot_proxy", "", "robot proxy: http://127.0.0.1:7890")
	flag.StringVar(&BaseConfInfo.Lang, "lang", "en", "lang")
	flag.IntVar(&BaseConfInfo.TokenPerUser, "token_per_user", 10000, "token per user")
	flag.IntVar(&BaseConfInfo.MaxUserChat, "max_user_chat", 2, "max chat per user")
	flag.StringVar(&BaseConfInfo.HTTPHost, "http_host", ":36060", "http server port")
	flag.BoolVar(&BaseConfInfo.UseTools, "use_tools", false, "use function tools")
	flag.IntVar(&BaseConfInfo.MaxQAPair, "max_qa_pari", 100, "max qa pair")
	flag.StringVar(&BaseConfInfo.Character, "character", "", "ai's character")
	flag.StringVar(&BaseConfInfo.Powered, "powered", "", "powered by")
	flag.StringVar(&BaseConfInfo.ImagePath, "image_path", "./conf/img/", "image path")

	flag.StringVar(&BaseConfInfo.CrtFile, "crt_file", "", "public key file")
	flag.StringVar(&BaseConfInfo.KeyFile, "key_file", "", "secret key file")
	flag.StringVar(&BaseConfInfo.CaFile, "ca_file", "", "ca file")
	flag.BoolVar(&BaseConfInfo.SendMcpRes, "send_mcp_res", false, "send mcp res")
	flag.BoolVar(&BaseConfInfo.SendMcpMediaToLLM, "send_mcp_media_to_llm", false, "send mcp media to llm")
	flag.StringVar(&BaseConfInfo.DefaultModel, "default_model", "", "default model")
	flag.IntVar(&BaseConfInfo.LLMRetryTimes, "llm_retry_times", 3, "llm retry times")
	flag.IntVar(&BaseConfInfo.LLMRetryInterval, "llm_retry_interval", 100, "llm retry interval")
	flag.BoolVar(&BaseConfInfo.LLMOptionParam, "llm_option_param", false, "llm option param")
	flag.BoolVar(&BaseConfInfo.IsStreaming, "is_streaming", false, "is streaming")

	allowedUserIds := flag.String("allowed_user_ids", "", "allowed user ids")
	allowedGroupIds := flag.String("allowed_group_ids", "", "allowed group ids")

	BaseConfInfo.AllowedUserIds = make(map[string]bool)
	BaseConfInfo.AllowedGroupIds = make(map[string]bool)

	InitLLMConf()
	InitToolsConf()

	flag.CommandLine.Init(os.Args[0], flag.ContinueOnError)
	flag.Parse()

	if os.Getenv("TELEGRAM_BOT_TOKEN") != "" {
		BaseConfInfo.TelegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	}

	if os.Getenv("CHAT_ANY_WHERE_TOKEN") != "" {
		BaseConfInfo.ChatAnyWhereToken = os.Getenv("CHAT_ANY_WHERE_TOKEN")
	}

	if os.Getenv("DISCORD_BOT_TOKEN") != "" {
		BaseConfInfo.DiscordBotToken = os.Getenv("DISCORD_BOT_TOKEN")
	}

	if os.Getenv("SLACK_BOT_TOKEN") != "" {
		BaseConfInfo.SlackBotToken = os.Getenv("SLACK_BOT_TOKEN")
	}

	if os.Getenv("SLACK_APP_TOKEN") != "" {
		BaseConfInfo.SlackAppToken = os.Getenv("SLACK_APP_TOKEN")
	}

	if os.Getenv("LARK_APP_ID") != "" {
		BaseConfInfo.LarkAPPID = os.Getenv("LARK_APP_ID")
	}

	if os.Getenv("LARK_APP_SECRET") != "" {
		BaseConfInfo.LarkAppSecret = os.Getenv("LARK_APP_SECRET")
	}

	if os.Getenv("DING_CLIENT_ID") != "" {
		BaseConfInfo.DingClientId = os.Getenv("DING_CLIENT_ID")
	}

	if os.Getenv("DING_CLIENT_SECRET") != "" {
		BaseConfInfo.DingClientSecret = os.Getenv("DING_CLIENT_SECRET")
	}

	if os.Getenv("COM_WECHAT_TOKEN") != "" {
		BaseConfInfo.ComWechatToken = os.Getenv("COM_WECHAT_TOKEN")
	}

	if os.Getenv("WECHAT_TOKEN") != "" {
		BaseConfInfo.WechatToken = os.Getenv("WECHAT_TOKEN")
	}

	if os.Getenv("WECHAT_APP_ID") != "" {
		BaseConfInfo.WechatAppID = os.Getenv("WECHAT_APP_ID")
	}

	if os.Getenv("WECHAT_APP_SECRET") != "" {
		BaseConfInfo.WechatAppSecret = os.Getenv("WECHAT_APP_SECRET")
	}

	if os.Getenv("WECHAT_ENCODING_AES_KEY") != "" {
		BaseConfInfo.WechatEncodingAESKey = os.Getenv("WECHAT_ENCODING_AES_KEY")
	}

	if os.Getenv("WECHAT_ACTIVE") != "" {
		BaseConfInfo.WechatActive = os.Getenv("WECHAT_ACTIVE") == "true"
	}

	if os.Getenv("COM_WECHAT_ENCODING_AES_KEY") != "" {
		BaseConfInfo.ComWechatEncodingAESKey = os.Getenv("COM_WECHAT_ENCODING_AES_KEY")
	}

	if os.Getenv("COM_WECHAT_CORP_ID") != "" {
		BaseConfInfo.ComWechatCorpID = os.Getenv("COM_WECHAT_CORP_ID")
	}

	if os.Getenv("COM_WECHAT_SECRET") != "" {
		BaseConfInfo.ComWechatSecret = os.Getenv("COM_WECHAT_SECRET")
	}

	if os.Getenv("COM_WECHAT_AGENT_ID") != "" {
		BaseConfInfo.ComWechatAgentID = os.Getenv("COM_WECHAT_AGENT_ID")
	}

	if os.Getenv("QQ_APP_ID") != "" {
		BaseConfInfo.QQAppID = os.Getenv("QQ_APP_ID")
	}

	if os.Getenv("QQ_APP_SECRET") != "" {
		BaseConfInfo.QQAppSecret = os.Getenv("QQ_APP_SECRET")
	}

	if os.Getenv("QQ_ONEBOT_SEND_TOKEN") != "" {
		BaseConfInfo.QQOneBotSendToken = os.Getenv("QQ_ONEBOT_SEND_TOKEN")
	}

	if os.Getenv("QQ_ONEBOT_RECEIVE_TOKEN") != "" {
		BaseConfInfo.QQOneBotReceiveToken = os.Getenv("QQ_ONEBOT_RECEIVE_TOKEN")
	}

	if os.Getenv("QQ_ONEBOT_HTTP_SERVER") != "" {
		BaseConfInfo.QQOneBotHttpServer = os.Getenv("QQ_ONEBOT_HTTP_SERVER")
	}

	if os.Getenv("DEEPSEEK_TOKEN") != "" {
		BaseConfInfo.DeepseekToken = os.Getenv("DEEPSEEK_TOKEN")
	}

	if os.Getenv("CUSTOM_URL") != "" {
		BaseConfInfo.CustomUrl = os.Getenv("CUSTOM_URL")
	}

	if os.Getenv("BOT_NAME") != "" {
		BaseConfInfo.BotName = os.Getenv("BOT_NAME")
	}

	if os.Getenv("TYPE") != "" {
		BaseConfInfo.Type = os.Getenv("TYPE")
	}

	if os.Getenv("VOLC_AK") != "" {
		BaseConfInfo.VolcAK = os.Getenv("VOLC_AK")
	}

	if os.Getenv("VOLC_SK") != "" {
		BaseConfInfo.VolcSK = os.Getenv("VOLC_SK")
	}

	if os.Getenv("DB_TYPE") != "" {
		BaseConfInfo.DBType = os.Getenv("DB_TYPE")
	}

	if os.Getenv("DB_CONF") != "" {
		BaseConfInfo.DBConf = os.Getenv("DB_CONF")
	}

	if os.Getenv("ALLOWED_USER_IDS") != "" {
		*allowedUserIds = os.Getenv("ALLOWED_USER_IDS")
	}

	if os.Getenv("ALLOWED_GROUP_IDS") != "" {
		*allowedGroupIds = os.Getenv("ALLOWED_GROUP_IDS")
	}

	if os.Getenv("LLM_PROXY") != "" {
		BaseConfInfo.LLMProxy = os.Getenv("LLM_PROXY")
	}

	if os.Getenv("ROBOT_PROXY") != "" {
		BaseConfInfo.RobotProxy = os.Getenv("ROBOT_PROXY")
	}

	if os.Getenv("LANG") != "" {
		BaseConfInfo.Lang = os.Getenv("LANG")
	}

	if os.Getenv("TOKEN_PER_USER") != "" {
		BaseConfInfo.TokenPerUser, _ = strconv.Atoi(os.Getenv("TOKEN_PER_USER"))
	}

	if os.Getenv("MAX_USER_CHAT") != "" {
		BaseConfInfo.MaxUserChat, _ = strconv.Atoi(os.Getenv("MAX_USER_CHAT"))
	}

	if os.Getenv("HTTP_HOST") != "" {
		BaseConfInfo.HTTPHost = os.Getenv("HTTP_HOST")
	}

	if os.Getenv("USE_TOOLS") != "" {
		BaseConfInfo.UseTools = os.Getenv("USE_TOOLS") == "true"
	}

	if os.Getenv("OPENAI_TOKEN") != "" {
		BaseConfInfo.OpenAIToken = os.Getenv("OPENAI_TOKEN")
	}

	if os.Getenv("GEMINI_TOKEN") != "" {
		BaseConfInfo.GeminiToken = os.Getenv("GEMINI_TOKEN")
	}

	if os.Getenv("VOL_TOKEN") != "" {
		BaseConfInfo.VolToken = os.Getenv("VOL_TOKEN")
	}

	if os.Getenv("ALIYUN_TOKEN") != "" {
		BaseConfInfo.AliyunToken = os.Getenv("ALIYUN_TOKEN")
	}

	if os.Getenv("ERNIE_AK") != "" {
		BaseConfInfo.ErnieAK = os.Getenv("ERNIE_AK")
	}

	if os.Getenv("ERNIE_SK") != "" {
		BaseConfInfo.ErnieSK = os.Getenv("ERNIE_SK")
	}

	if os.Getenv("OPEN_ROUTER_TOKEN") != "" {
		BaseConfInfo.OpenRouterToken = os.Getenv("OPEN_ROUTER_TOKEN")
	}

	if os.Getenv("AI_302_TOKEN") != "" {
		BaseConfInfo.AI302Token = os.Getenv("AI_302_TOKEN")
	}

	if os.Getenv("MAX_QA_PAIR") != "" {
		BaseConfInfo.MaxQAPair, _ = strconv.Atoi(os.Getenv("MAX_QA_PAIR"))
	}

	if os.Getenv("CHARACTER") != "" {
		BaseConfInfo.Character = os.Getenv("CHARACTER")
	}

	if os.Getenv("CRT_FILE") != "" {
		BaseConfInfo.CrtFile = os.Getenv("CRT_FILE")
	}

	if os.Getenv("KEY_FILE") != "" {
		BaseConfInfo.KeyFile = os.Getenv("KEY_FILE")
	}

	if os.Getenv("CA_FILE") != "" {
		BaseConfInfo.CaFile = os.Getenv("CA_FILE")
	}

	if os.Getenv("MEDIA_TYPE") != "" {
		BaseConfInfo.MediaType = os.Getenv("MEDIA_TYPE")
	}

	if os.Getenv("SMART_MODE") != "" {
		BaseConfInfo.SmartMode = os.Getenv("SMART_MODE") == "true"
	}

	if os.Getenv("CONTEXT_EXPIRE_TIME") != "" {
		BaseConfInfo.ContextExpireTime, _ = strconv.Atoi(os.Getenv("CONTEXT_EXPIRE_TIME"))
	}

	if os.Getenv("POWERED") != "" {
		BaseConfInfo.Powered = os.Getenv("POWERED")
	}

	if os.Getenv("SEND_MCP_RES") != "" {
		BaseConfInfo.SendMcpRes = os.Getenv("SEND_MCP_RES") == "true"
	}

	if os.Getenv("DEFAULT_MODEL") != "" {
		BaseConfInfo.DefaultModel = os.Getenv("DEFAULT_MODEL")
	}

	if os.Getenv("LLM_RETRY_TIMES") != "" {
		BaseConfInfo.LLMRetryTimes, _ = strconv.Atoi(os.Getenv("LLM_RETRY_TIMES"))
	}

	if os.Getenv("LLM_RETRY_INTERVAL") != "" {
		BaseConfInfo.LLMRetryInterval, _ = strconv.Atoi(os.Getenv("LLM_RETRY_INTERVAL"))
	}

	if os.Getenv("LLM_OPTION_PARAM") != "" {
		BaseConfInfo.LLMOptionParam = os.Getenv("LLM_OPTION_PARAM") == "true"
	}

	if os.Getenv("IMAGE_PATH") != "" {
		BaseConfInfo.ImagePath = os.Getenv("IMAGE_PATH")
	}

	if os.Getenv("IS_STREAMING") != "" {
		BaseConfInfo.IsStreaming = os.Getenv("IS_STREAMING") == "true"
	}

	if os.Getenv("SEND_MCP_MEDIA_TO_LLM") == "true" {
		BaseConfInfo.SendMcpMediaToLLM = true
	}

	EnvLLMConf()
	EnvToolsConf()

	logConf(*allowedUserIds, *allowedGroupIds)
	SaveConf()

}

func logConf(allowedUserIds, allowedGroupIds string) {
	for _, userIdStr := range strings.Split(allowedUserIds, ",") {
		if userIdStr == "" {
			continue
		}
		BaseConfInfo.AllowedUserIds[userIdStr] = true
	}

	for _, groupIdStr := range strings.Split(allowedGroupIds, ",") {
		if groupIdStr == "" {
			continue
		}
		BaseConfInfo.AllowedGroupIds[groupIdStr] = true
	}

	logger.Info("CONF", "TelegramBotToken", BaseConfInfo.TelegramBotToken)
	logger.Info("CONF", "DiscordBotToken", BaseConfInfo.DiscordBotToken)
	logger.Info("CONF", "SlackBotToken", BaseConfInfo.SlackBotToken)
	logger.Info("CONF", "SlackAppToken", BaseConfInfo.SlackAppToken)
	logger.Info("CONF", "LarkAPPID", BaseConfInfo.LarkAPPID)
	logger.Info("CONF", "LarkAppSecret", BaseConfInfo.LarkAppSecret)
	logger.Info("CONF", "DingClientId", BaseConfInfo.DingClientId)
	logger.Info("CONF", "DingClientSecret", BaseConfInfo.DingClientSecret)
	logger.Info("CONF", "ComWechatToken", BaseConfInfo.ComWechatToken)
	logger.Info("CONF", "ComWechatEncodingAESKey", BaseConfInfo.ComWechatEncodingAESKey)
	logger.Info("CONF", "ComWechatCorpID", BaseConfInfo.ComWechatCorpID)
	logger.Info("CONF", "ComWechatSecret", BaseConfInfo.ComWechatSecret)
	logger.Info("CONF", "ComWechatAgentID", BaseConfInfo.ComWechatAgentID)
	logger.Info("CONF", "WechatToken", BaseConfInfo.WechatToken)
	logger.Info("CONF", "WechatAppSecret", BaseConfInfo.WechatAppSecret)
	logger.Info("CONF", "WechatAppID", BaseConfInfo.WechatAppID)
	logger.Info("CONF", "WechatActive", BaseConfInfo.WechatActive)
	logger.Info("CONF", "WechatEncodingAESKey", BaseConfInfo.WechatEncodingAESKey)
	logger.Info("CONF", "QQAppID", BaseConfInfo.QQAppID)
	logger.Info("CONF", "QQAppSecret", BaseConfInfo.QQAppSecret)
	logger.Info("CONF", "QQOneBotHttpServer", BaseConfInfo.QQOneBotHttpServer)
	logger.Info("CONF", "QQOneBotReceiveToken", BaseConfInfo.QQOneBotReceiveToken)
	logger.Info("CONF", "QQOneBotSendToken", BaseConfInfo.QQOneBotSendToken)
	logger.Info("CONF", "DeepseekToken", BaseConfInfo.DeepseekToken)
	logger.Info("CONF", "CustomUrl", BaseConfInfo.CustomUrl)
	logger.Info("CONF", "Type", BaseConfInfo.Type)
	logger.Info("CONF", "VolcAK", BaseConfInfo.VolcAK)
	logger.Info("CONF", "VolcSK", BaseConfInfo.VolcSK)
	logger.Info("CONF", "AliyunToken", BaseConfInfo.AliyunToken)
	logger.Info("CONF", "DBType", BaseConfInfo.DBType)
	logger.Info("CONF", "DBConf", BaseConfInfo.DBConf)
	logger.Info("CONF", "AllowedUserIds", BaseConfInfo.AllowedUserIds)
	logger.Info("CONF", "AllowedGroupIds", BaseConfInfo.AllowedGroupIds)
	logger.Info("CONF", "LLMProxy", BaseConfInfo.LLMProxy)
	logger.Info("CONF", "RobotProxy", BaseConfInfo.RobotProxy)
	logger.Info("CONF", "Lang", BaseConfInfo.Lang)
	logger.Info("CONF", "TokenPerUser", BaseConfInfo.TokenPerUser)
	logger.Info("CONF", "MaxUserChat", BaseConfInfo.MaxUserChat)
	logger.Info("CONF", "HTTPHost", BaseConfInfo.HTTPHost)
	logger.Info("CONF", "UseTools", BaseConfInfo.UseTools)
	logger.Info("CONF", "OpenAIToken", BaseConfInfo.OpenAIToken)
	logger.Info("CONF", "GeminiToken", BaseConfInfo.GeminiToken)
	logger.Info("CONF", "OpenRouterToken", BaseConfInfo.OpenRouterToken)
	logger.Info("CONF", "AI302Token", BaseConfInfo.AI302Token)
	logger.Info("CONF", "ErnieAK", BaseConfInfo.ErnieAK)
	logger.Info("CONF", "ErnieSK", BaseConfInfo.ErnieSK)
	logger.Info("CONF", "VolToken", BaseConfInfo.VolToken)
	logger.Info("CONF", "CrtFile", BaseConfInfo.CrtFile)
	logger.Info("CONF", "KeyFile", BaseConfInfo.KeyFile)
	logger.Info("CONF", "CaFile", BaseConfInfo.CaFile)
	logger.Info("CONF", "MediaType", BaseConfInfo.MediaType)
	logger.Info("CONF", "BotName", BaseConfInfo.BotName)
	logger.Info("CONF", "MaxQAPair", BaseConfInfo.MaxQAPair)
	logger.Info("CONF", "SmartMode", BaseConfInfo.SmartMode)
	logger.Info("CONF", "Powered", BaseConfInfo.Powered)
	logger.Info("CONF", "Character", BaseConfInfo.Character)
	logger.Info("CONF", "ContextExpireTime", BaseConfInfo.ContextExpireTime)
	logger.Info("CONF", "SendMcpRes", BaseConfInfo.SendMcpRes)
	logger.Info("CONF", "DefaultModel", BaseConfInfo.DefaultModel)
	logger.Info("CONF", "LLMRetryTimes", BaseConfInfo.LLMRetryTimes)
	logger.Info("CONF", "LLMRetryInterval", BaseConfInfo.LLMRetryInterval)
	logger.Info("CONF", "LLMOptionParam", BaseConfInfo.LLMOptionParam)
	logger.Info("CONF", "ImagePath", BaseConfInfo.ImagePath)
	logger.Info("CONF", "IsStreaming", BaseConfInfo.IsStreaming)
	logger.Info("CONF", "SendMcpMediaToLLM", BaseConfInfo.SendMcpMediaToLLM)

	logger.Info("AUDIO_CONF", "AudioAppID", AudioConfInfo.VolAudioAppID)
	logger.Info("AUDIO_CONF", "AudioToken", AudioConfInfo.VolAudioToken)
	logger.Info("AUDIO_CONF", "AudioCluster", AudioConfInfo.VolAudioRecCluster)
	logger.Info("AUDIO_CONF", "AudioVoiceType", AudioConfInfo.VolAudioVoiceType)
	logger.Info("AUDIO_CONF", "AudioTTSCluster", AudioConfInfo.VolAudioTTSCluster)
	logger.Info("AUDIO_CONF", "GeminiAudioModel", AudioConfInfo.GeminiAudioModel)
	logger.Info("AUDIO_CONF", "GeminiVoiceName", AudioConfInfo.GeminiVoiceName)
	logger.Info("AUDIO_CONF", "OpenAIAudioModel", AudioConfInfo.OpenAIAudioModel)
	logger.Info("AUDIO_CONF", "OpenAIVoiceName", AudioConfInfo.OpenAIVoiceName)
	logger.Info("AUDIO_CONF", "TTSType", AudioConfInfo.TTSType)
	logger.Info("AUDIO_CONF", "VolEndSmoothWindow", AudioConfInfo.VolEndSmoothWindow)
	logger.Info("AUDIO_CONF", "VolTTSSpeaker", AudioConfInfo.VolTTSSpeaker)
	logger.Info("AUDIO_CONF", "VolBotName", AudioConfInfo.VolBotName)
	logger.Info("AUDIO_CONF", "VolSystemRole", AudioConfInfo.VolSystemRole)
	logger.Info("AUDIO_CONF", "VolSpeakingStyle", AudioConfInfo.VolSpeakingStyle)
	logger.Info("AUDIO_CONF", "AliyunAudioModel", AudioConfInfo.AliyunAudioModel)
	logger.Info("AUDIO_CONF", "AliyunAudioVoice", AudioConfInfo.AliyunAudioVoice)
	logger.Info("AUDIO_CONF", "AliyunAudioRecModel", AudioConfInfo.AliyunAudioRecModel)

	logger.Info("RAG_CONF", "EmbeddingType", RagConfInfo.EmbeddingType)
	logger.Info("RAG_CONF", "KnowledgePath", RagConfInfo.KnowledgePath)
	logger.Info("RAG_CONF", "VectorDBType", RagConfInfo.VectorDBType)
	logger.Info("RAG_CONF", "ChromaURL", RagConfInfo.ChromaURL)
	logger.Info("RAG_CONF", "ChromaSpace", RagConfInfo.Space)
	logger.Info("RAG_CONF", "MilvusURL", RagConfInfo.MilvusURL)
	logger.Info("RAG_CONF", "WeaviateURL", RagConfInfo.WeaviateURL)
	logger.Info("RAG_CONF", "WeaviateScheme", RagConfInfo.WeaviateScheme)

	logger.Info("PHOTO_CONF", "ReqKey", PhotoConfInfo.ReqKey)
	logger.Info("PHOTO_CONF", "ModelVersion", PhotoConfInfo.ModelVersion)
	logger.Info("PHOTO_CONF", "ReqScheduleConf", PhotoConfInfo.ReqScheduleConf)
	logger.Info("PHOTO_CONF", "Seed", PhotoConfInfo.Seed)
	logger.Info("PHOTO_CONF", "Width", PhotoConfInfo.Width)
	logger.Info("PHOTO_CONF", "Height", PhotoConfInfo.Height)
	logger.Info("PHOTO_CONF", "Scale", PhotoConfInfo.Scale)
	logger.Info("PHOTO_CONF", "DDIMSteps", PhotoConfInfo.DDIMSteps)
	logger.Info("PHOTO_CONF", "UsePreLLM", PhotoConfInfo.UsePreLLM)
	logger.Info("PHOTO_CONF", "UseSr", PhotoConfInfo.UseSr)
	logger.Info("PHOTO_CONF", "ReturnUrl", PhotoConfInfo.ReturnUrl)
	logger.Info("PHOTO_CONF", "AddLogo", PhotoConfInfo.AddLogo)
	logger.Info("PHOTO_CONF", "Position", PhotoConfInfo.Position)
	logger.Info("PHOTO_CONF", "Language", PhotoConfInfo.Language)
	logger.Info("PHOTO_CONF", "Opacity", PhotoConfInfo.Opacity)
	logger.Info("PHOTO_CONF", "LogoTextContent", PhotoConfInfo.LogoTextContent)
	logger.Info("PHOTO_CONF", "GeminiImageModel", PhotoConfInfo.GeminiImageModel)
	logger.Info("PHOTO_CONF", "GeminiRecModel", PhotoConfInfo.GeminiRecModel)
	logger.Info("PHOTO_CONF", "OpenAIImageStyle", PhotoConfInfo.OpenAIImageStyle)
	logger.Info("PHOTO_CONF", "OpenAIImageModel", PhotoConfInfo.OpenAIImageModel)
	logger.Info("PHOTO_CONF", "OpenAIImageSize", PhotoConfInfo.OpenAIImageSize)
	logger.Info("PHOTO_CONF", "OpenAIRecModel", PhotoConfInfo.OpenAIRecModel)
	logger.Info("PHOTO_CONF", "VolImageModel", PhotoConfInfo.VolImageModel)
	logger.Info("PHOTO_CONF", "VolRecModel", PhotoConfInfo.VolRecModel)
	logger.Info("PHOTO_CONF", "AI302ImageModel", PhotoConfInfo.MixRecModel)
	logger.Info("PHOTO_CONF", "AI302RecModel", PhotoConfInfo.MixRecModel)
	logger.Info("PHOTO_CONF", "AliyunImageModel", PhotoConfInfo.AliyunImageModel)
	logger.Info("PHOTO_CONF", "AliyunRecModel", PhotoConfInfo.AliyunRecModel)

	logger.Info("VIDEO_CONF", "VOL_VIDEO_MODEL", VideoConfInfo.VolVideoModel)
	logger.Info("VIDEO_CONF", "RADIO", VideoConfInfo.Radio)
	logger.Info("VIDEO_CONF", "DURATION", VideoConfInfo.Duration)
	logger.Info("VIDEO_CONF", "FPS", VideoConfInfo.FPS)
	logger.Info("VIDEO_CONF", "RESOLUTION", VideoConfInfo.Resolution)
	logger.Info("VIDEO_CONF", "WATERMARK", VideoConfInfo.Watermark)
	logger.Info("AUDIO_CONF", "GeminiVideoModel", VideoConfInfo.GeminiVideoModel)
	logger.Info("AUDIO_CONF", "AI302VideoModel", VideoConfInfo.AI302VideoModel)
	logger.Info("AUDIO_CONF", "AliyunVideoModel", VideoConfInfo.AliyunVideoModel)

	logger.Info("REGISTER_CONF", "Type", RegisterConfInfo.Type)
	logger.Info("REGISTER_CONF", "EtcdURLs", RegisterConfInfo.EtcdURLs)
	logger.Info("REGISTER_CONF", "EtcdUsername", RegisterConfInfo.EtcdUsername)
	logger.Info("REGISTER_CONF", "EtcdPassword", RegisterConfInfo.EtcdPassword)

	logger.Info("LLM_CONF", "FrequencyPenalty", LLMConfInfo.FrequencyPenalty)
	logger.Info("LLM_CONF", "MaxTokens", LLMConfInfo.MaxTokens)
	logger.Info("LLM_CONF", "PresencePenalty", LLMConfInfo.PresencePenalty)
	logger.Info("LLM_CONF", "Temperature", LLMConfInfo.Temperature)
	logger.Info("LLM_CONF", "TopP", LLMConfInfo.TopP)
	logger.Info("LLM_CONF", "Stop", LLMConfInfo.Stop)
	logger.Info("LLM_CONF", "LogProbs", LLMConfInfo.LogProbs)
	logger.Info("LLM_CONF", "TopLogProbs", LLMConfInfo.TopLogProbs)

	logger.Info("TOOLS_CONF", "McpConfPath", *ToolsConfInfo.McpConfPath)
}

func GetAbsPath(relPath string) string {
	exe, err := os.Executable()
	if err != nil {
		logger.Error("Failed to get executable path", "err", err)
		return ""
	}
	dir := filepath.Dir(exe)
	return filepath.Join(dir, relPath)
}

func loadConf() bool {
	m := make(map[string]string)
	for _, part := range os.Args {
		if strings.HasPrefix(part, "-") {
			kv := strings.SplitN(part[1:], "=", 2)
			if len(kv) == 2 {
				m[kv[0]] = kv[1]
			}
		}
	}

	if !(len(m) == 0 || (len(m) == 1 && (m["bot_name"] != "" || m["http_host"] != "")) ||
		(len(m) == 2 && m["bot_name"] != "" && m["http_host"] != "")) {
		return false
	}

	data, err := os.ReadFile(getSaveConf(m))
	if err != nil {
		return false
	}

	err = json.Unmarshal(data, &AllConf)
	if err != nil {
		logger.Error("Failed to parse config file", "err", err)
		return false
	}

	err = TransferMapToConf(AllConf["base"].(map[string]interface{}), BaseConfInfo)
	if err != nil {
		logger.Error("Failed to transfer map to base conf", "err", err)
		return false
	}

	err = TransferMapToConf(AllConf["audio"].(map[string]interface{}), AudioConfInfo)
	if err != nil {
		logger.Error("Failed to transfer map to audio conf", "err", err)
		return false
	}

	err = TransferMapToConf(AllConf["llm"].(map[string]interface{}), LLMConfInfo)
	if err != nil {
		logger.Error("Failed to transfer map to llm conf", "err", err)
		return false
	}

	err = TransferMapToConf(AllConf["photo"].(map[string]interface{}), PhotoConfInfo)
	if err != nil {
		logger.Error("Failed to transfer map to photo conf", "err", err)
		return false
	}

	err = TransferMapToConf(AllConf["rag"].(map[string]interface{}), RagConfInfo)
	if err != nil {
		logger.Error("Failed to transfer map to rag conf", "err", err)
		return false
	}

	err = TransferMapToConf(AllConf["video"].(map[string]interface{}), VideoConfInfo)
	if err != nil {
		logger.Error("Failed to transfer map to video conf", "err", err)
		return false
	}

	err = TransferMapToConf(AllConf["register"].(map[string]interface{}), RegisterConfInfo)
	if err != nil {
		logger.Error("Failed to transfer map to register conf", "err", err)
		return false
	}

	err = TransferMapToConf(AllConf["tools"].(map[string]interface{}), ToolsConfInfo)
	if err != nil {
		logger.Error("Failed to transfer map to tools conf", "err", err)
		return false
	}

	return true
}

func SaveConf() {
	AllConf["base"] = BaseConfInfo
	AllConf["audio"] = AudioConfInfo
	AllConf["llm"] = LLMConfInfo
	AllConf["photo"] = PhotoConfInfo
	AllConf["rag"] = RagConfInfo
	AllConf["video"] = VideoConfInfo
	AllConf["register"] = RegisterConfInfo
	AllConf["tools"] = ToolsConfInfo

	fileName := getSaveConf(map[string]string{
		"bot_name":  BaseConfInfo.BotName,
		"http_host": BaseConfInfo.HTTPHost,
	})

	confData, err := json.Marshal(AllConf)
	if err != nil {
		logger.Error("Failed to marshal config data", "err", err)
		return
	}

	err = os.WriteFile(fileName, confData, 0644)
	if err != nil {
		logger.Error("Failed to write config file", "err", err)
		return
	}

}

func getSaveConf(m map[string]string) string {
	botName := m["bot_name"]
	if botName == "" {
		botName = "Accountabel AI"
	}

	httpHost := m["http_host"]
	if httpHost == "" {
		httpHost = ":36060"
	}
	httpHost = NormalizeHTTP(httpHost)

	hash := md5.Sum([]byte(httpHost))
	md5Str := hex.EncodeToString(hash[:])
	return GetAbsPath(botName + md5Str + ".json")
}

func NormalizeHTTP(addr string) string {
	if strings.HasPrefix(addr, ":") {
		addr = "127.0.0.1" + addr
	}
	if !strings.HasPrefix(addr, "http://") {
		addr = "http://" + addr
	}
	return addr
}

func TransferMapToConf(m map[string]interface{}, conf interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, conf)
}

// Stubs for legacy config options that are no longer used but referenced
type AudioConf struct {
	VolAudioAppID        string `json:"vol_audio_app_id"`
	VolAudioToken        string `json:"vol_audio_token"`
	VolAudioRecCluster   string `json:"vol_audio_rec_cluster"`
	VolAudioVoiceType    string `json:"vol_audio_voice_type"`
	VolAudioTTSCluster   string `json:"vol_audio_tts_cluster"`
	GeminiAudioModel     string `json:"gemini_audio_model"`
	GeminiVoiceName      string `json:"gemini_voice_name"`
	OpenAIAudioModel     string `json:"openai_audio_model"`
	OpenAIVoiceName      string `json:"openai_voice_name"`
	TTSType              string `json:"tts_type"`
	VolEndSmoothWindow   int    `json:"vol_end_smooth_window"`
	VolTTSSpeaker        string `json:"vol_tts_speaker"`
	VolBotName           string `json:"vol_bot_name"`
	VolSystemRole        string `json:"vol_system_role"`
	VolSpeakingStyle     string `json:"vol_speaking_style"`
	AliyunAudioModel     string `json:"aliyun_audio_model"`
	AliyunAudioVoice     string `json:"aliyun_audio_voice"`
	AliyunAudioRecModel  string `json:"aliyun_audio_rec_model"`
}

type PhotoConf struct {
	ReqKey           string  `json:"req_key"`
	ModelVersion     string  `json:"model_version"`
	ReqScheduleConf  string  `json:"req_schedule_conf"`
	Seed             int     `json:"seed"`
	Width            int     `json:"width"`
	Height           int     `json:"height"`
	Scale            float64 `json:"scale"`
	DDIMSteps        int     `json:"ddim_steps"`
	UsePreLLM        bool    `json:"use_pre_llm"`
	UseSr            bool    `json:"use_sr"`
	ReturnUrl        bool    `json:"return_url"`
	AddLogo          bool    `json:"add_logo"`
	Position         string  `json:"position"`
	Language         int     `json:"photo_language"`
	Opacity          float64 `json:"opacity"`
	LogoTextContent  string  `json:"logo_text_content"`
	GeminiImageModel string  `json:"gemini_image_model"`
	GeminiRecModel   string  `json:"gemini_rec_model"`
	OpenAIImageStyle string  `json:"openai_image_style"`
	OpenAIImageModel string  `json:"openai_image_model"`
	OpenAIImageSize  string  `json:"openai_image_size"`
	OpenAIRecModel   string  `json:"openai_rec_model"`
	VolImageModel    string  `json:"vol_image_model"`
	VolRecModel      string  `json:"vol_rec_model"`
	MixRecModel      string  `json:"mix_rec_model"`
	AliyunImageModel string  `json:"aliyun_image_model"`
	AliyunRecModel   string  `json:"aliyun_rec_model"`
}

type RagConf struct {
	EmbeddingType string      `json:"embedding_type"`
	KnowledgePath string      `json:"knowledge_path"`
	VectorDBType  string      `json:"vector_db_type"`
	ChromaURL     string      `json:"chroma_url"`
	Space         string      `json:"space"`
	ChunkSize     int         `json:"chunk_size"`
	ChunkOverlap  int         `json:"chunk_overlap"`
	MilvusURL     string      `json:"milvus_url"`
	WeaviateURL   string      `json:"weaviate_url"`
	WeaviateScheme string     `json:"weaviate_scheme"`
	Store         interface{} `json:"-"`
}

type VideoConf struct {
	VolVideoModel    string `json:"vol_video_model"`
	Radio            string `json:"radio"`
	Duration         int    `json:"duration"`
	FPS              int    `json:"fps"`
	Resolution       string `json:"resolution"`
	Watermark        bool   `json:"watermark"`
	GeminiVideoModel string `json:"gemini_video_model"`
	AI302VideoModel  string `json:"ai_302_video_model"`
	AliyunVideoModel string `json:"aliyun_video_model"`
}

type RegisterConf struct {
	Type         string   `json:"type"`
	EtcdURLs     []string `json:"etcd_urls"`
	EtcdUsername string   `json:"etcd_username"`
	EtcdPassword string   `json:"etcd_password"`
}

var (
	AudioConfInfo    = new(AudioConf)
	PhotoConfInfo    = new(PhotoConf)
	RagConfInfo      = new(RagConf)
	VideoConfInfo    = new(VideoConf)
	RegisterConfInfo = new(RegisterConf)
)

func InitAudioConf() {
	flag.StringVar(&AudioConfInfo.VolAudioAppID, "vol_audio_app_id", "", "vol audio app id")
	flag.StringVar(&AudioConfInfo.VolAudioToken, "vol_audio_token", "", "vol audio token")
	flag.StringVar(&AudioConfInfo.VolAudioRecCluster, "vol_audio_rec_cluster", "volcengine_input_common", "vol audio cluster")
	flag.StringVar(&AudioConfInfo.VolAudioVoiceType, "vol_audio_voice_type", "", "vol audio voice type")
	flag.StringVar(&AudioConfInfo.VolAudioTTSCluster, "vol_audio_tts_cluster", "volcano_tts", "vol audio tts cluster")
	flag.IntVar(&AudioConfInfo.VolEndSmoothWindow, "vol_end_smooth_window", 1500, "vol end smooth window")
	flag.StringVar(&AudioConfInfo.VolTTSSpeaker, "vol_tts_speaker", "zh_female_vv_jupiter_bigtts", "vol tts speaker")
	flag.StringVar(&AudioConfInfo.VolBotName, "vol_bot_name", "豆包", "vol bot name")
	flag.StringVar(&AudioConfInfo.VolSystemRole, "vol_system_role", "你使用活泼灵动的女声，性格开朗，热爱生活。", "vol system role")
	flag.StringVar(&AudioConfInfo.VolSpeakingStyle, "vol_speaking_style", "你的说话风格简洁明了，语速适中，语调自然。", "vol speaking style")
	flag.StringVar(&AudioConfInfo.GeminiAudioModel, "gemini_audio_model", "gemini-2.5-flash-preview-tts", "gemini audio model")
	flag.StringVar(&AudioConfInfo.GeminiVoiceName, "gemini_voice_name", "Kore", "gemini voice name")
	flag.StringVar(&AudioConfInfo.OpenAIAudioModel, "openai_audio_model", "tts-1", "openai audio model")
	flag.StringVar(&AudioConfInfo.OpenAIVoiceName, "openai_voice_name", "alloy", "openai voice name")
	flag.StringVar(&AudioConfInfo.AliyunAudioModel, "aliyun_audio_model", "qwen3-tts-flash", "aliyun audio model")
	flag.StringVar(&AudioConfInfo.AliyunAudioVoice, "aliyun_audio_voice", "Cherry", "aliyun audio voice")
	flag.StringVar(&AudioConfInfo.AliyunAudioRecModel, "aliyun_audio_rec_model", "qwen-audio-turbo-latest", "aliyun audio rec model")
	flag.StringVar(&AudioConfInfo.TTSType, "tts_type", "", "vol tts type: 1. vol 2. gemini")
}

func InitPhotoConf() {
	flag.StringVar(&PhotoConfInfo.ReqKey, "req_key", "high_aes_general_v21_L", "request key")
	flag.StringVar(&PhotoConfInfo.ModelVersion, "model_version", "general_v2.1_L", "model version")
	flag.StringVar(&PhotoConfInfo.ReqScheduleConf, "req_schedule_conf", "general_v20_9B_pe", "request schedule conf")
	flag.IntVar(&PhotoConfInfo.Seed, "seed", -1, "seed for random seed")
	flag.Float64Var(&PhotoConfInfo.Scale, "scale", 3.5, "scale factor")
	flag.IntVar(&PhotoConfInfo.DDIMSteps, "ddim_steps", 25, "ddim steps")
	flag.IntVar(&PhotoConfInfo.Width, "width", 512, "width of the image")
	flag.IntVar(&PhotoConfInfo.Height, "height", 512, "height of the image")
	flag.BoolVar(&PhotoConfInfo.UsePreLLM, "use_pre_llm", true, "use pre llm")
	flag.BoolVar(&PhotoConfInfo.UseSr, "use_sr", true, "use super resolution")
	flag.BoolVar(&PhotoConfInfo.ReturnUrl, "return_url", true, "return url")
	flag.BoolVar(&PhotoConfInfo.AddLogo, "add_logo", false, "add logo")
	flag.StringVar(&PhotoConfInfo.Position, "position", "", "position")
	flag.IntVar(&PhotoConfInfo.Language, "photo_language", 1, "language")
	flag.Float64Var(&PhotoConfInfo.Opacity, "opacity", 0.3, "opacity")
	flag.StringVar(&PhotoConfInfo.LogoTextContent, "logo_text_content", "", "logo text content")
	flag.StringVar(&PhotoConfInfo.GeminiImageModel, "gemini_image_model", "gemini-2.0-flash-preview-image-generation", "gemini create photo model")
	flag.StringVar(&PhotoConfInfo.GeminiRecModel, "gemini_rec_model", "gemini-2.0-flash", "gemini recognize photo model")
	flag.StringVar(&PhotoConfInfo.OpenAIRecModel, "openai_rec_model", "chatgpt-4o-latest", "openai create photo model")
	flag.StringVar(&PhotoConfInfo.OpenAIImageModel, "openai_image_model", "gpt-image-1", "openai create photo model")
	flag.StringVar(&PhotoConfInfo.OpenAIImageSize, "openai_image_size", "1024x1024", "openai image size")
	flag.StringVar(&PhotoConfInfo.OpenAIImageStyle, "openai_image_style", "", "openai image style")
	flag.StringVar(&PhotoConfInfo.VolImageModel, "vol_image_model", "doubao-seed-1-6-250615", "vol image model")
	flag.StringVar(&PhotoConfInfo.VolRecModel, "vol_rec_model", "doubao-seed-1-6-250615", "vol recognize photo model")
	flag.StringVar(&PhotoConfInfo.MixRecModel, "mix_rec_model", "chatgpt-4o-latest", "ai302/openrouter recognize photo model")
	flag.StringVar(&PhotoConfInfo.AliyunImageModel, "aliyun_image_model", "qwen-image-plus", "aliyun image model")
	flag.StringVar(&PhotoConfInfo.AliyunRecModel, "aliyun_rec_model", "qwen-vl-max-latest", "aliyun recognize photo model")
}

func InitRagConf() {
	flag.StringVar(&RagConfInfo.EmbeddingType, "embedding_type", "", "embedding split api: openai gemini ernie")
	flag.StringVar(&RagConfInfo.KnowledgePath, "knowledge_path", GetAbsPath("data/knowledge"), "knowledge")
	flag.StringVar(&RagConfInfo.VectorDBType, "vector_db_type", "milvus", "vector db type: chroma weaviate milvus")
	flag.StringVar(&RagConfInfo.ChromaURL, "chroma_url", "http://localhost:8000", "chroma url")
	flag.StringVar(&RagConfInfo.MilvusURL, "milvus_url", "http://localhost:19530", "milvus url")
	flag.StringVar(&RagConfInfo.WeaviateURL, "weaviate_url", "localhost:8000", "weaviate url localhost:8000")
	flag.StringVar(&RagConfInfo.WeaviateScheme, "weaviate_scheme", "http", "weaviate scheme: http")
	flag.StringVar(&RagConfInfo.Space, "space", "Accountabel AI", "chroma space")
	flag.IntVar(&RagConfInfo.ChunkSize, "chunk_size", 500, "rag file chunk size")
	flag.IntVar(&RagConfInfo.ChunkOverlap, "chunk_overlap", 50, "rag file chunk overlap")
}

func InitVideoConf() {
	flag.StringVar(&VideoConfInfo.VolVideoModel, "vol_video_model", "doubao-seedance-1-0-pro-250528", "video model")
	flag.StringVar(&VideoConfInfo.Radio, "radio", "1:1", "the width to height ratio")
	flag.IntVar(&VideoConfInfo.Duration, "duration", 5, "the duration in seconds, only support 5s / 10s")
	flag.IntVar(&VideoConfInfo.FPS, "fps", 24, "the frame per second")
	flag.StringVar(&VideoConfInfo.Resolution, "resolution", "480p", "the resolution of video, only support 480p / 720p")
	flag.BoolVar(&VideoConfInfo.Watermark, "watermark", false, "include watermark")
	flag.StringVar(&VideoConfInfo.GeminiVideoModel, "gemini_video_model", "veo-2.0-generate-001", "create video model")
	flag.StringVar(&VideoConfInfo.AI302VideoModel, "ai_302_video_model", "luma_video", "create video model")
	flag.StringVar(&VideoConfInfo.AliyunVideoModel, "aliyun_video_model", "wan2.5-t2v-preview", "create video model")
}

func InitRegisterConf() {}

func EnvAudioConf() {
	if os.Getenv("VOL_AUDIO_APP_ID") != "" {
		AudioConfInfo.VolAudioAppID = os.Getenv("VOL_AUDIO_APP_ID")
	}
	if os.Getenv("VOL_AUDIO_TOKEN") != "" {
		AudioConfInfo.VolAudioToken = os.Getenv("VOL_AUDIO_TOKEN")
	}
	if os.Getenv("VOL_AUDIO_REC_CLUSTER") != "" {
		AudioConfInfo.VolAudioRecCluster = os.Getenv("VOL_AUDIO_REC_CLUSTER")
	}
	if os.Getenv("VOL_AUDIO_VOICE_TYPE") != "" {
		AudioConfInfo.VolAudioVoiceType = os.Getenv("VOL_AUDIO_VOICE_TYPE")
	}
	if os.Getenv("VOL_AUDIO_TTS_CLUSTER") != "" {
		AudioConfInfo.VolAudioTTSCluster = os.Getenv("VOL_AUDIO_TTS_CLUSTER")
	}
	if os.Getenv("GEMINI_AUDIO_MODEL") != "" {
		AudioConfInfo.GeminiAudioModel = os.Getenv("GEMINI_AUDIO_MODEL")
	}
	if os.Getenv("GEMINI_VOICE_NAME") != "" {
		AudioConfInfo.GeminiVoiceName = os.Getenv("GEMINI_VOICE_NAME")
	}
	if os.Getenv("OPENAI_AUDIO_MODEL") != "" {
		AudioConfInfo.OpenAIAudioModel = os.Getenv("OPENAI_AUDIO_MODEL")
	}
	if os.Getenv("OPENAI_VOICE_NAME") != "" {
		AudioConfInfo.OpenAIVoiceName = os.Getenv("OPENAI_VOICE_NAME")
	}
	if os.Getenv("TTS_TYPE") != "" {
		AudioConfInfo.TTSType = os.Getenv("TTS_TYPE")
	}
	if os.Getenv("VOL_END_SMOOTH_WINDOW") != "" {
		AudioConfInfo.VolEndSmoothWindow, _ = strconv.Atoi(os.Getenv("VOL_END_SMOOTH_WINDOW"))
	}
	if os.Getenv("VOL_TTS_SPEAKER") != "" {
		AudioConfInfo.VolTTSSpeaker = os.Getenv("VOL_TTS_SPEAKER")
	}
	if os.Getenv("VOL_BOT_NAME") != "" {
		AudioConfInfo.VolBotName = os.Getenv("VOL_BOT_NAME")
	}
	if os.Getenv("VOL_SYSTEM_ROLE") != "" {
		AudioConfInfo.VolSystemRole = os.Getenv("VOL_SYSTEM_ROLE")
	}
	if os.Getenv("VOL_SPEAKING_STYLE") != "" {
		AudioConfInfo.VolSpeakingStyle = os.Getenv("VOL_SPEAKING_STYLE")
	}
	if os.Getenv("ALIYUN_AUDIO_MODEL") != "" {
		AudioConfInfo.AliyunAudioModel = os.Getenv("ALIYUN_AUDIO_MODEL")
	}
	if os.Getenv("ALIYUN_AUDIO_VOICE") != "" {
		AudioConfInfo.AliyunAudioVoice = os.Getenv("ALIYUN_AUDIO_VOICE")
	}
	if os.Getenv("ALIYUN_AUDIO_REC_MODEL") != "" {
		AudioConfInfo.AliyunAudioRecModel = os.Getenv("ALIYUN_AUDIO_REC_MODEL")
	}
}

func EnvPhotoConf() {
	if os.Getenv("REQ_KEY") != "" {
		PhotoConfInfo.ReqKey = os.Getenv("REQ_KEY")
	}
	if os.Getenv("MODEL_VERSION") != "" {
		PhotoConfInfo.ModelVersion = os.Getenv("MODEL_VERSION")
	}
	if os.Getenv("REQ_SCHEDULE_CONF") != "" {
		PhotoConfInfo.ReqScheduleConf = os.Getenv("REQ_SCHEDULE_CONF")
	}
	if os.Getenv("SEED") != "" {
		PhotoConfInfo.Seed, _ = strconv.Atoi(os.Getenv("SEED"))
	}
	if os.Getenv("SCALE") != "" {
		PhotoConfInfo.Scale, _ = strconv.ParseFloat(os.Getenv("SCALE"), 64)
	}
	if os.Getenv("DDIM_Steps") != "" {
		PhotoConfInfo.DDIMSteps, _ = strconv.Atoi(os.Getenv("DDIM_Steps"))
	}
	if os.Getenv("WIDTH") != "" {
		PhotoConfInfo.Width, _ = strconv.Atoi(os.Getenv("WIDTH"))
	}
	if os.Getenv("HEIGHT") != "" {
		PhotoConfInfo.Height, _ = strconv.Atoi(os.Getenv("HEIGHT"))
	}
	if os.Getenv("USE_PER_LLM") != "" {
		PhotoConfInfo.UsePreLLM, _ = strconv.ParseBool(os.Getenv("USE_PER_LLM"))
	}
	if os.Getenv("USE_SR") != "" {
		PhotoConfInfo.UseSr, _ = strconv.ParseBool(os.Getenv("USE_SR"))
	}
	if os.Getenv("RETURN_URL") != "" {
		PhotoConfInfo.ReturnUrl, _ = strconv.ParseBool(os.Getenv("RETURN_URL"))
	}
	if os.Getenv("ADD_LOGO") != "" {
		PhotoConfInfo.AddLogo, _ = strconv.ParseBool(os.Getenv("ADD_LOGO"))
	}
	if os.Getenv("POSITION") != "" {
		PhotoConfInfo.Position = os.Getenv("POSITION")
	}
	if os.Getenv("PHOTO_LANGUAGE") != "" {
		PhotoConfInfo.Language, _ = strconv.Atoi(os.Getenv("PHOTO_LANGUAGE"))
	}
	if os.Getenv("OPACITY") != "" {
		PhotoConfInfo.Opacity, _ = strconv.ParseFloat(os.Getenv("OPACITY"), 64)
	}
	if os.Getenv("LOGO_TEXT_CONTENT") != "" {
		PhotoConfInfo.LogoTextContent = os.Getenv("LOGO_TEXT_CONTENT")
	}
}

func EnvRagConf() {
	if os.Getenv("EMBEDDING_TYPE") != "" {
		RagConfInfo.EmbeddingType = os.Getenv("EMBEDDING_TYPE")
	}
	if os.Getenv("KNOWLEDGE_PATH") != "" {
		RagConfInfo.KnowledgePath = os.Getenv("KNOWLEDGE_PATH")
	}
	if os.Getenv("VECTOR_DB_TYPE") != "" {
		RagConfInfo.VectorDBType = os.Getenv("VECTOR_DB_TYPE")
	}
	if os.Getenv("CHROMA_URL") != "" {
		RagConfInfo.ChromaURL = os.Getenv("CHROMA_URL")
	}
	if os.Getenv("MILVUS_URL") != "" {
		RagConfInfo.MilvusURL = os.Getenv("MILVUS_URL")
	}
	if os.Getenv("WEAVIATE_SCHEME") != "" {
		RagConfInfo.WeaviateScheme = os.Getenv("WEAVIATE_SCHEME")
	}
	if os.Getenv("WEAVIATE_URL") != "" {
		RagConfInfo.WeaviateURL = os.Getenv("WEAVIATE_URL")
	}
	if os.Getenv("SPACE") != "" {
		RagConfInfo.Space = os.Getenv("SPACE")
	}
	if os.Getenv("CHUNK_SIZE") != "" {
		RagConfInfo.ChunkSize, _ = strconv.Atoi(os.Getenv("CHUNK_SIZE"))
	}
	if os.Getenv("CHUNK_OVERLAP") != "" {
		RagConfInfo.ChunkOverlap, _ = strconv.Atoi(os.Getenv("CHUNK_OVERLAP"))
	}
}

func EnvVideoConf() {
	if os.Getenv("VOL_VIDEO_MODEL") != "" {
		VideoConfInfo.VolVideoModel = os.Getenv("VOL_VIDEO_MODEL")
	}
	if os.Getenv("RADIO") != "" {
		VideoConfInfo.Radio = os.Getenv("RADIO")
	}
	if os.Getenv("DURATION") != "" {
		VideoConfInfo.Duration, _ = strconv.Atoi(os.Getenv("DURATION"))
	}
	if os.Getenv("FPS") != "" {
		VideoConfInfo.FPS, _ = strconv.Atoi(os.Getenv("FPS"))
	}
	if os.Getenv("RESOLUTION") != "" {
		VideoConfInfo.Resolution = os.Getenv("RESOLUTION")
	}
	if os.Getenv("WATERMARK") != "" {
		VideoConfInfo.Watermark, _ = strconv.ParseBool(os.Getenv("WATERMARK"))
	}
}

func EnvRegisterConf() {}

