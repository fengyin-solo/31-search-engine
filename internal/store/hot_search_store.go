package store

import "searchengine/internal/model"

// CreateHotSearch 创建热门搜索词。
func (s *MemoryStore) CreateHotSearch(h *model.HotSearch) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.hotSearches[h.ID] = h
	return nil
}

// GetHotSearch 按 ID 获取热门搜索词。
func (s *MemoryStore) GetHotSearch(id string) (*model.HotSearch, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	h, ok := s.hotSearches[id]
	if !ok {
		return nil, ErrNotFound
	}
	return h, nil
}

// GetHotSearchByTerm 按词获取热门搜索词。
func (s *MemoryStore) GetHotSearchByTerm(term string) (*model.HotSearch, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, h := range s.hotSearches {
		if h.Term == term {
			return h, nil
		}
	}
	return nil, ErrNotFound
}

// ListHotSearches 列出全部热门搜索词。
func (s *MemoryStore) ListHotSearches() []*model.HotSearch {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.HotSearch, 0, len(s.hotSearches))
	for _, h := range s.hotSearches {
		list = append(list, h)
	}
	return list
}

// UpdateHotSearch 更新热门搜索词。
func (s *MemoryStore) UpdateHotSearch(h *model.HotSearch) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.hotSearches[h.ID]; !ok {
		return ErrNotFound
	}
	s.hotSearches[h.ID] = h
	return nil
}

// DeleteHotSearch 删除热门搜索词。
func (s *MemoryStore) DeleteHotSearch(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.hotSearches[id]; !ok {
		return ErrNotFound
	}
	delete(s.hotSearches, id)
	return nil
}
