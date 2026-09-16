package model

import (
	"strings"
)

// Term 词典词条（倒排词典）。
type Term struct {
	ID       string `json:"id"`
	IndexID  string `json:"index_id"`
	Term     string `json:"term"`
	DocCount int    `json:"doc_count"`
	TotalTF  int    `json:"total_tf"`
}

// Validate 校验并规范化词条字段。
func (t *Term) Validate() error {
	t.IndexID = strings.TrimSpace(t.IndexID)
	t.Term = strings.ToLower(strings.TrimSpace(t.Term))
	if t.IndexID == "" {
		return NewValidationError("index_id", "索引不能为空")
	}
	if t.Term == "" {
		return NewValidationError("term", "词条不能为空")
	}
	if t.DocCount < 0 || t.TotalTF < 0 {
		return NewValidationError("count", "统计值不能为负")
	}
	return nil
}
