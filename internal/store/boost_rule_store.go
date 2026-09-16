package store

import "searchengine/internal/model"

// CreateBoostRule 创建加权规则。
func (s *MemoryStore) CreateBoostRule(b *model.BoostRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.boostRules[b.ID] = b
	return nil
}

// GetBoostRule 按 ID 获取加权规则。
func (s *MemoryStore) GetBoostRule(id string) (*model.BoostRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.boostRules[id]
	if !ok {
		return nil, ErrNotFound
	}
	return b, nil
}

// ListBoostRules 列出全部加权规则。
func (s *MemoryStore) ListBoostRules() []*model.BoostRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.BoostRule, 0, len(s.boostRules))
	for _, b := range s.boostRules {
		list = append(list, b)
	}
	return list
}

// UpdateBoostRule 更新加权规则。
func (s *MemoryStore) UpdateBoostRule(b *model.BoostRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.boostRules[b.ID]; !ok {
		return ErrNotFound
	}
	s.boostRules[b.ID] = b
	return nil
}

// DeleteBoostRule 删除加权规则。
func (s *MemoryStore) DeleteBoostRule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.boostRules[id]; !ok {
		return ErrNotFound
	}
	delete(s.boostRules, id)
	return nil
}
