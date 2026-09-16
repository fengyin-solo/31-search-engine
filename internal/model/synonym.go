package model

import (
	"strings"
	"time"
)

const (
	SynonymActive   = "active"
	SynonymInactive = "inactive"
)

// Synonym 同义词。
type Synonym struct {
	ID        string    `json:"id"`
	Word      string    `json:"word"`
	Synonyms  []string  `json:"synonyms"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate 校验并规范化同义词字段。
func (s *Synonym) Validate() error {
	s.Word = strings.ToLower(strings.TrimSpace(s.Word))
	if s.Word == "" {
		return NewValidationError("word", "原词不能为空")
	}
	if s.Status == "" {
		s.Status = SynonymActive
	}
	if s.Status != SynonymActive && s.Status != SynonymInactive {
		return NewValidationError("status", "同义词状态不合法")
	}
	seen := make(map[string]bool)
	words := make([]string, 0, len(s.Synonyms))
	for _, w := range s.Synonyms {
		w = strings.ToLower(strings.TrimSpace(w))
		if w == "" || w == s.Word || seen[w] {
			continue
		}
		seen[w] = true
		words = append(words, w)
	}
	s.Synonyms = words
	return nil
}
