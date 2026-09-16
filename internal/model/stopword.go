package model

import (
	"strings"
	"time"
)

// StopWord 停用词。
type StopWord struct {
	ID        string    `json:"id"`
	Word      string    `json:"word"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate 校验并规范化停用词字段。
func (s *StopWord) Validate() error {
	s.Word = strings.ToLower(strings.TrimSpace(s.Word))
	if s.Word == "" {
		return NewValidationError("word", "停用词不能为空")
	}
	return nil
}
