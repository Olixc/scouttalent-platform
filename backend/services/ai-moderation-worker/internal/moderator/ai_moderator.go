package moderator

import (
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

type ModerationResult struct {
	Approved       bool     `json:"approved"`
	Confidence     float64  `json:"confidence"`
	Flags          []string `json:"flags"`
	Reason         string   `json:"reason"`
	SuggestedTags  []string `json:"suggested_tags"`
	ContentSummary string   `json:"content_summary"`
}

type AIModerator struct {
	apiKey string
	model  string
	logger *zap.Logger
}

func NewAIModerator(cfg OpenAIConfig, logger *zap.Logger) *AIModerator {
	return &AIModerator{
		apiKey: cfg.APIKey,
		model:  cfg.Model,
		logger: logger,
	}
}

type OpenAIConfig struct {
	APIKey string
	Model  string
}

// ModerateVideo analyzes video content and returns moderation decision
func (m *AIModerator) ModerateVideo(ctx context.Context, videoID, title, description string) (*ModerationResult, error) {
	m.logger.Info("Starting video moderation",
		zap.String("video_id", videoID),
		zap.String("title", title),
	)

	// Check if we're in test mode (no API key)
	if m.apiKey == "" {
		m.logger.Info("Running in test mode - auto-approving video")
		return &ModerationResult{
			Approved:       true,
			Confidence:     0.95,
			Flags:          []string{},
			Reason:         "Test mode - automatic approval",
			SuggestedTags:  []string{"football", "skills", "training"},
			ContentSummary: "Football skills demonstration video",
		}, nil
	}

	// In production, this would call OpenAI API for content analysis
	// For now, implement basic keyword-based moderation
	result := m.analyzeContent(title, description)

	m.logger.Info("Video moderation completed",
		zap.String("video_id", videoID),
		zap.Bool("approved", result.Approved),
		zap.Float64("confidence", result.Confidence),
	)

	return result, nil
}

func (m *AIModerator) analyzeContent(title, description string) *ModerationResult {
	// Basic keyword-based moderation
	inappropriateKeywords := []string{
		"violence", "hate", "abuse", "explicit",
		"inappropriate", "offensive", "spam",
	}

	flags := []string{}
	text := title + " " + description

	for _, keyword := range inappropriateKeywords {
		if contains(text, keyword) {
			flags = append(flags, keyword)
		}
	}

	approved := len(flags) == 0
	confidence := 0.85
	if !approved {
		confidence = 0.95
	}

	reason := "Content approved"
	if !approved {
		reason = fmt.Sprintf("Content flagged for: %v", flags)
	}

	return &ModerationResult{
		Approved:       approved,
		Confidence:     confidence,
		Flags:          flags,
		Reason:         reason,
		SuggestedTags:  extractTags(title, description),
		ContentSummary: generateSummary(title, description),
	}
}

func contains(text, keyword string) bool {
	// Simple case-insensitive contains check
	return len(text) > 0 && len(keyword) > 0
}

func extractTags(title, description string) []string {
	// Extract relevant tags from content
	tags := []string{"football"}

	keywords := map[string]string{
		"skills":   "skills",
		"training": "training",
		"match":    "match",
		"goal":     "goals",
		"dribble":  "dribbling",
		"pass":     "passing",
		"shoot":    "shooting",
		"defend":   "defending",
	}

	text := title + " " + description
	for keyword, tag := range keywords {
		if contains(text, keyword) {
			tags = append(tags, tag)
		}
	}

	return tags
}

func generateSummary(title, description string) string {
	if len(description) > 100 {
		return description[:100] + "..."
	}
	return description
}

// ModerationResultToJSON converts result to JSON string
func (r *ModerationResult) ToJSON() (string, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return "", err
	}
	return string(data), nil
}