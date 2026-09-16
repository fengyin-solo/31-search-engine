package model

import (
	"strings"
	"time"
)

const (
	BoostFieldTitle = "title"
	BoostFieldBody  = "body"

	BoostRuleActive   = "active"
	BoostRuleInactive = "inactive"
)

// BoostRule 加权规则：命中指定字段的词条在打分时乘以权重。
type BoostRule struct {
	ID        string    `json:"id"`
	Field     string    `json:"field"`
	Term      string    `json:"term"`
	Boost     float64   `json:"boost"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate 校验并规范化加权规则字段。
func (b *BoostRule) Validate() error {
	b.Field = strings.TrimSpace(b.Field)
	b.Term = strings.ToLower(strings.TrimSpace(b.Term))
	if b.Field == "" {
		b.Field = BoostFieldTitle
	}
	if b.Field != BoostFieldTitle && b.Field != BoostFieldBody {
		return NewValidationError("field", "加权字段不合法")
	}
	if b.Term == "" {
		return NewValidationError("term", "加权词条不能为空")
	}
	if b.Boost <= 0 {
		return NewValidationError("boost", "权重必须大于 0")
	}
	if b.Status == "" {
		b.Status = BoostRuleActive
	}
	if b.Status != BoostRuleActive && b.Status != BoostRuleInactive {
		return NewValidationError("status", "加权规则状态不合法")
	}
	return nil
}
