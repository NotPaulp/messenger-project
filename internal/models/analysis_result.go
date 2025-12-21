package models

import "time"

type AnalysisResult struct {
	MessageID      int64     `json:"message_id" bson:"message_id"`
	Category       string    `json:"category" bson:"category"`
	ToxicityScore  float32   `json:"toxicity_score" bson:"toxicity_score"`
	IsToxic        bool      `json:"toxic" bson:"toxic"`
	AnalyzedAt     time.Time `json:"analyzed_at" bson:"analyzed_at"`
	SenderUsername string    `json:"sender_username" bson:"sender_username"`
}
