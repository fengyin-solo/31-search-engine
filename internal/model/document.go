package model

import (
	"strings"
	"time"
)

const (
	DocumentActive  = "active"
	DocumentDeleted = "deleted"
)

// Document 待索引文档。
type Document struct {
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	URL        string    `json:"url"`
	Source     string    `json:"source"`
	CategoryID string    `json:"category_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Validate 校验并规范化文档字段。
func (d *Document) Validate() error {
	d.Title = strings.TrimSpace(d.Title)
	d.Body = strings.TrimSpace(d.Body)
	d.URL = strings.TrimSpace(d.URL)
	d.Source = strings.TrimSpace(d.Source)
	d.CategoryID = strings.TrimSpace(d.CategoryID)
	if d.Title == "" {
		return NewValidationError("title", "文档标题不能为空")
	}
	if d.Body == "" {
		return NewValidationError("body", "文档正文不能为空")
	}
	if d.Status == "" {
		d.Status = DocumentActive
	}
	if d.Status != DocumentActive && d.Status != DocumentDeleted {
		return NewValidationError("status", "文档状态不合法")
	}
	return nil
}

// DocumentFilter 文档列表筛选条件。
type DocumentFilter struct {
	Status     string
	Source     string
	CategoryID string
	Keyword    string
}

// Match 判断文档是否命中筛选条件。
func (f DocumentFilter) Match(d *Document) bool {
	if f.Status != "" && d.Status != f.Status {
		return false
	}
	if f.Source != "" && d.Source != f.Source {
		return false
	}
	if f.CategoryID != "" && d.CategoryID != f.CategoryID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(d.Title), k) &&
			!strings.Contains(strings.ToLower(d.Body), k) {
			return false
		}
	}
	return true
}
