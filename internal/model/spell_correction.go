package model

import (
	"strings"
	"time"
)

const (
	SpellCorrectionActive   = "active"
	SpellCorrectionInactive = "inactive"
)

// SpellCorrection 拼写纠错词对。
type SpellCorrection struct {
	ID        string    `json:"id"`
	Word      string    `json:"word"`
	Correct   string    `json:"correct"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate 校验并规范化拼写纠错字段。
func (s *SpellCorrection) Validate() error {
	s.Word = strings.ToLower(strings.TrimSpace(s.Word))
	s.Correct = strings.ToLower(strings.TrimSpace(s.Correct))
	if s.Word == "" {
		return NewValidationError("word", "错误拼写不能为空")
	}
	if s.Correct == "" {
		return NewValidationError("correct", "正确拼写不能为空")
	}
	if s.Status == "" {
		s.Status = SpellCorrectionActive
	}
	if s.Status != SpellCorrectionActive && s.Status != SpellCorrectionInactive {
		return NewValidationError("status", "纠错状态不合法")
	}
	return nil
}
