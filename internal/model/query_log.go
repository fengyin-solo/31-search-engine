package model

import (
	"strings"
	"time"
)

// QueryLog 查询记录。
type QueryLog struct {
	ID          string    `json:"id"`
	IndexID     string    `json:"index_id"`
	Query       string    `json:"query"`
	ResultCount int       `json:"result_count"`
	DurationMs  int       `json:"duration_ms"`
	CreatedAt   time.Time `json:"created_at"`
}

// Validate 校验并规范化查询记录字段。
func (q *QueryLog) Validate() error {
	q.IndexID = strings.TrimSpace(q.IndexID)
	q.Query = strings.ToLower(strings.TrimSpace(q.Query))
	if q.Query == "" {
		return NewValidationError("query", "查询词不能为空")
	}
	if q.ResultCount < 0 || q.DurationMs < 0 {
		return NewValidationError("count", "统计值不能为负")
	}
	return nil
}
