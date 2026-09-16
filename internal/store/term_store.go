package store

import "searchengine/internal/model"

// CreateTerm 创建词典词条。
func (s *MemoryStore) CreateTerm(t *model.Term) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := t.IndexID + "|" + t.Term
	for _, exist := range s.terms {
		if exist.IndexID+"|"+exist.Term == key {
			return ErrConflict
		}
	}
	s.terms[t.ID] = t
	return nil
}

// GetTerm 按 ID 获取词条。
func (s *MemoryStore) GetTerm(id string) (*model.Term, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.terms[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

// GetTermByKey 按 索引+词 获取词条。
func (s *MemoryStore) GetTermByKey(indexID, term string) (*model.Term, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, t := range s.terms {
		if t.IndexID == indexID && t.Term == term {
			return t, nil
		}
	}
	return nil, ErrNotFound
}

// ListTerms 列出全部词条。
func (s *MemoryStore) ListTerms() []*model.Term {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Term, 0, len(s.terms))
	for _, t := range s.terms {
		list = append(list, t)
	}
	return list
}

// UpdateTerm 更新词条。
func (s *MemoryStore) UpdateTerm(t *model.Term) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.terms[t.ID]; !ok {
		return ErrNotFound
	}
	s.terms[t.ID] = t
	return nil
}

// DeleteTerm 删除词条。
func (s *MemoryStore) DeleteTerm(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.terms[id]; !ok {
		return ErrNotFound
	}
	delete(s.terms, id)
	return nil
}
