package model

import (
	"strings"
	"time"
)

// HotSearch 热门搜索词。
type HotSearch struct {
	ID        string    `json:"id"`
	Term      string    `json:"term"`
	Count     int       `json:"count"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate 校验并规范化热门搜索词字段。
func (h *HotSearch) Validate() error {
	h.Term = strings.ToLower(strings.TrimSpace(h.Term))
	if h.Term == "" {
		return NewValidationError("term", "搜索词不能为空")
	}
	if h.Count < 0 {
		return NewValidationError("count", "热度不能为负")
	}
	return nil
}
