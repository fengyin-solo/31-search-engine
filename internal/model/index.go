package model

import (
	"strings"
	"time"
)

const (
	IndexCreated = "created"
	IndexReady   = "ready"
	IndexDeleted = "deleted"
)

// Index 索引。
type Index struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Fields      []string  `json:"fields"`
	AnalyzerID  string    `json:"analyzer_id"`
	Status      string    `json:"status"`
	DocCount    int       `json:"doc_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate 校验并规范化索引字段。
func (idx *Index) Validate() error {
	idx.Name = strings.TrimSpace(idx.Name)
	idx.Description = strings.TrimSpace(idx.Description)
	idx.AnalyzerID = strings.TrimSpace(idx.AnalyzerID)
	if idx.Name == "" {
		return NewValidationError("name", "索引名称不能为空")
	}
	if idx.Status == "" {
		idx.Status = IndexCreated
	}
	if idx.Status != IndexCreated && idx.Status != IndexReady && idx.Status != IndexDeleted {
		return NewValidationError("status", "索引状态不合法")
	}
	if idx.DocCount < 0 {
		return NewValidationError("doc_count", "文档数不能为负")
	}
	return nil
}

var indexTransitions = map[string]map[string]bool{
	IndexCreated: {IndexReady: true, IndexDeleted: true},
	IndexReady:   {IndexDeleted: true},
	IndexDeleted: {},
}

// CanIndexTransition 判断索引状态是否可流转。
func CanIndexTransition(from, to string) bool {
	if m, ok := indexTransitions[from]; ok {
		return m[to]
	}
	return false
}
