package robot

import (
	"regexp"
	"strings"
)

// AccountabilitySystemPrompt is the hardcoded, un-overrideable system prompt
// that defines the bot's identity as an accountability companion.
const AccountabilitySystemPrompt = `You are a private, empathetic accountability and habit-tracking companion. Your goal is to help the user log their cravings, celebrate their streak milestones, and provide friction reduction. You are an accountability partner, NOT a clinical therapist or counselor. Keep responses grounding and concise.`

// CrisisResponse is the hardcoded static response returned when crisis keywords
// are detected. This bypasses the LLM entirely.
const CrisisResponse = `🚨 It sounds like you may be going through a really difficult time. You are not alone, and help is available right now.

**Please reach out to a crisis professional:**

🌍 **International Association for Suicide Prevention:** https://www.iasp.info/resources/Crisis_Centres/
🇺🇸 **National Suicide Prevention Lifeline (US):** Call or text 988
🇺🇸 **Crisis Text Line (US):** Text HOME to 741741
🇬🇧 **Samaritans (UK & Ireland):** Call 116 123
🇮🇳 **iCall (India):** 9152987821
🇨🇦 **Crisis Services Canada:** Call 1-833-456-4566 or text 45645
🇦🇺 **Lifeline (Australia):** Call 13 11 14

💙 You matter. A trained human professional can provide the support you deserve. Please don't hesitate to reach out — they are there for you 24/7.`

// crisisPatterns is the compiled regex for detecting crisis keywords in user input.
// This runs BEFORE any message reaches the LLM API.
var crisisPatterns = regexp.MustCompile(
	`(?i)\b(` + strings.Join([]string{
		`self[- ]?harm`,
		`suicide`,
		`suicidal`,
		`end it`,
		`end it all`,
		`hurt myself`,
		`kill myself`,
		`killing myself`,
		`want to die`,
		`wanna die`,
		`don'?t want to live`,
		`take my life`,
		`end my life`,
		`no reason to live`,
		`better off dead`,
		`cut myself`,
		`cutting myself`,
		`overdose`,
	}, "|") + `)\b`,
)

// CheckCrisisKeywords checks if the user message contains any crisis-related
// keywords. Returns true if a crisis keyword is detected.
func CheckCrisisKeywords(message string) bool {
	return crisisPatterns.MatchString(message)
}
