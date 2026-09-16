package model

import (
	"strings"
	"time"
)

const (
	AnalyzerTypeStandard   = "standard"
	AnalyzerTypeWhitespace = "whitespace"

	AnalyzerActive   = "active"
	AnalyzerInactive = "inactive"
)

// Analyzer 分词器配置。
type Analyzer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	StopWords []string  `json:"stop_words"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate 校验并规范化分词器字段。
func (a *Analyzer) Validate() error {
	a.Name = strings.TrimSpace(a.Name)
	a.Type = strings.TrimSpace(a.Type)
	if a.Name == "" {
		return NewValidationError("name", "分词器名称不能为空")
	}
	if a.Type == "" {
		a.Type = AnalyzerTypeStandard
	}
	if a.Type != AnalyzerTypeStandard && a.Type != AnalyzerTypeWhitespace {
		return NewValidationError("type", "分词器类型不合法")
	}
	if a.Status == "" {
		a.Status = AnalyzerActive
	}
	if a.Status != AnalyzerActive && a.Status != AnalyzerInactive {
		return NewValidationError("status", "分词器状态不合法")
	}
	// 去重停用词
	seen := make(map[string]bool)
	words := make([]string, 0, len(a.StopWords))
	for _, w := range a.StopWords {
		w = strings.ToLower(strings.TrimSpace(w))
		if w == "" || seen[w] {
			continue
		}
		seen[w] = true
		words = append(words, w)
	}
	a.StopWords = words
	return nil
}

// StopWordSet 返回停用词集合。
func (a *Analyzer) StopWordSet() map[string]bool {
	set := make(map[string]bool, len(a.StopWords))
	for _, w := range a.StopWords {
		set[w] = true
	}
	return set
}
