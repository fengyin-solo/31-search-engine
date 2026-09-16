package store

import "searchengine/internal/model"

// CreateAnalyzer 创建分词器。
func (s *MemoryStore) CreateAnalyzer(a *model.Analyzer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.analyzers {
		if exist.Name == a.Name {
			return ErrConflict
		}
	}
	s.analyzers[a.ID] = a
	return nil
}

// GetAnalyzer 按 ID 获取分词器。
func (s *MemoryStore) GetAnalyzer(id string) (*model.Analyzer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.analyzers[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

// GetAnalyzerByName 按名称获取分词器。
func (s *MemoryStore) GetAnalyzerByName(name string) (*model.Analyzer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, a := range s.analyzers {
		if a.Name == name {
			return a, nil
		}
	}
	return nil, ErrNotFound
}

// ListAnalyzers 列出全部分词器。
func (s *MemoryStore) ListAnalyzers() []*model.Analyzer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Analyzer, 0, len(s.analyzers))
	for _, a := range s.analyzers {
		list = append(list, a)
	}
	return list
}

// UpdateAnalyzer 更新分词器。
func (s *MemoryStore) UpdateAnalyzer(a *model.Analyzer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.analyzers[a.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.analyzers {
		if exist.ID != a.ID && exist.Name == a.Name {
			return ErrConflict
		}
	}
	s.analyzers[a.ID] = a
	return nil
}

// DeleteAnalyzer 删除分词器。
func (s *MemoryStore) DeleteAnalyzer(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.analyzers[id]; !ok {
		return ErrNotFound
	}
	delete(s.analyzers, id)
	return nil
}
