package service

import (
	"sort"
	"time"

	"searchengine/internal/model"
	"searchengine/pkg/idgen"
)

// CreateBoostRule 创建加权规则。
func (s *Service) CreateBoostRule(input model.BoostRule) (*model.BoostRule, error) {
	input.ID = idgen.Hex()
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateBoostRule(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

// GetBoostRule 获取加权规则。
func (s *Service) GetBoostRule(id string) (*model.BoostRule, error) {
	return s.store.GetBoostRule(id)
}

// ListBoostRules 分页列出加权规则。
func (s *Service) ListBoostRules(page, size int) ([]*model.BoostRule, int, error) {
	all := s.store.ListBoostRules()
	sort.Slice(all, func(i, j int) bool { return all[i].CreatedAt.After(all[j].CreatedAt) })
	total := len(all)
	start := (page - 1) * size
	if start >= total {
		return []*model.BoostRule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

// UpdateBoostRule 更新加权规则。
func (s *Service) UpdateBoostRule(id string, input model.BoostRule) (*model.BoostRule, error) {
	existing, err := s.store.GetBoostRule(id)
	if err != nil {
		return nil, err
	}
	existing.Field = input.Field
	existing.Term = input.Term
	existing.Boost = input.Boost
	existing.Status = input.Status
	existing.UpdatedAt = time.Now()
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateBoostRule(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteBoostRule 删除加权规则。
func (s *Service) DeleteBoostRule(id string) error {
	return s.store.DeleteBoostRule(id)
}

// activeBoostRules 返回全部启用的加权规则。
func (s *Service) activeBoostRules() []*model.BoostRule {
	var rules []*model.BoostRule
	for _, b := range s.store.ListBoostRules() {
		if b.Status == model.BoostRuleActive {
			rules = append(rules, b)
		}
	}
	return rules
}
