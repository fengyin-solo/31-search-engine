package store

import "searchengine/internal/model"

// CreateSuggestion 创建建议词。
func (s *MemoryStore) CreateSuggestion(v *model.Suggestion) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.suggestions {
		if exist.Term == v.Term {
			return ErrConflict
		}
	}
	s.suggestions[v.ID] = v
	return nil
}

// GetSuggestion 按 ID 获取建议词。
func (s *MemoryStore) GetSuggestion(id string) (*model.Suggestion, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.suggestions[id]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

// GetSuggestionByTerm 按词获取建议词。
func (s *MemoryStore) GetSuggestionByTerm(term string) (*model.Suggestion, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.suggestions {
		if v.Term == term {
			return v, nil
		}
	}
	return nil, ErrNotFound
}

// ListSuggestions 列出全部建议词。
func (s *MemoryStore) ListSuggestions() []*model.Suggestion {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Suggestion, 0, len(s.suggestions))
	for _, v := range s.suggestions {
		list = append(list, v)
	}
	return list
}

// UpdateSuggestion 更新建议词。
func (s *MemoryStore) UpdateSuggestion(v *model.Suggestion) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.suggestions[v.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.suggestions {
		if exist.ID != v.ID && exist.Term == v.Term {
			return ErrConflict
		}
	}
	s.suggestions[v.ID] = v
	return nil
}

// DeleteSuggestion 删除建议词。
func (s *MemoryStore) DeleteSuggestion(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.suggestions[id]; !ok {
		return ErrNotFound
	}
	delete(s.suggestions, id)
	return nil
}
