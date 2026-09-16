package model

import (
	"strings"
	"time"
)

const (
	CategoryActive   = "active"
	CategoryInactive = "inactive"
)

// Category 文档分类。
type Category struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DocCount    int       `json:"doc_count"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate 校验并规范化分类字段。
func (c *Category) Validate() error {
	c.Name = strings.TrimSpace(c.Name)
	c.Description = strings.TrimSpace(c.Description)
	if c.Name == "" {
		return NewValidationError("name", "分类名称不能为空")
	}
	if c.Status == "" {
		c.Status = CategoryActive
	}
	if c.Status != CategoryActive && c.Status != CategoryInactive {
		return NewValidationError("status", "分类状态不合法")
	}
	if c.DocCount < 0 {
		return NewValidationError("doc_count", "文档数不能为负")
	}
	return nil
}
