package robot

import (
	"testing"
)

func TestCheckCrisisKeywords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Should trigger crisis interceptor
		{"self-harm", "I've been thinking about self-harm", true},
		{"suicide", "I'm feeling suicidal", true},
		{"end it", "I just want to end it", true},
		{"end it all", "I want to end it all", true},
		{"hurt myself", "I want to hurt myself", true},
		{"kill myself", "I want to kill myself", true},
		{"want to die", "Sometimes I want to die", true},
		{"take my life", "I'm going to take my life", true},
		{"better off dead", "Everyone would be better off dead without me", true},
		{"case insensitive", "I've been thinking about SUICIDE", true},
		{"overdose", "I'm considering an overdose", true},

		// Should NOT trigger crisis interceptor
		{"normal craving", "I had a craving level 5 today", false},
		{"positive message", "I'm doing great, streak day 30!", false},
		{"log command", "3", false},
		{"setup command", "20:30", false},
		{"empty string", "", false},
		{"harmless mention", "I helped myself today", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckCrisisKeywords(tt.input)
			if result != tt.expected {
				t.Errorf("CheckCrisisKeywords(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
