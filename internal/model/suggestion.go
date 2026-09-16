package model

import (
	"strings"
	"time"
)

// Suggestion 搜索建议词（自动补全候选）。
type Suggestion struct {
	ID        string    `json:"id"`
	Term      string    `json:"term"`
	Weight    int       `json:"weight"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate 校验并规范化建议词字段。
func (s *Suggestion) Validate() error {
	s.Term = strings.ToLower(strings.TrimSpace(s.Term))
	if s.Term == "" {
		return NewValidationError("term", "建议词不能为空")
	}
	if s.Weight < 0 {
		return NewValidationError("weight", "权重不能为负")
	}
	return nil
}
