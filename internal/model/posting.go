package model

import (
	"strings"
)

// Posting 倒排记录（倒排表条目）。
type Posting struct {
	ID        string `json:"id"`
	IndexID   string `json:"index_id"`
	Term      string `json:"term"`
	DocID     string `json:"doc_id"`
	TF        int    `json:"tf"`
	Positions []int  `json:"positions"`
}

// Validate 校验并规范化倒排记录字段。
func (p *Posting) Validate() error {
	p.IndexID = strings.TrimSpace(p.IndexID)
	p.Term = strings.ToLower(strings.TrimSpace(p.Term))
	p.DocID = strings.TrimSpace(p.DocID)
	if p.IndexID == "" {
		return NewValidationError("index_id", "索引不能为空")
	}
	if p.Term == "" {
		return NewValidationError("term", "词条不能为空")
	}
	if p.DocID == "" {
		return NewValidationError("doc_id", "文档不能为空")
	}
	if p.TF < 0 {
		return NewValidationError("tf", "词频不能为负")
	}
	return nil
}
