package store

import "searchengine/internal/model"

// CreateIndex 创建索引。
func (s *MemoryStore) CreateIndex(i *model.Index) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.indexes {
		if exist.Name == i.Name {
			return ErrConflict
		}
	}
	s.indexes[i.ID] = i
	return nil
}

// GetIndex 按 ID 获取索引。
func (s *MemoryStore) GetIndex(id string) (*model.Index, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	i, ok := s.indexes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return i, nil
}

// GetIndexByName 按名称获取索引。
func (s *MemoryStore) GetIndexByName(name string) (*model.Index, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, i := range s.indexes {
		if i.Name == name {
			return i, nil
		}
	}
	return nil, ErrNotFound
}

// ListIndexes 列出全部索引。
func (s *MemoryStore) ListIndexes() []*model.Index {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Index, 0, len(s.indexes))
	for _, i := range s.indexes {
		list = append(list, i)
	}
	return list
}

// UpdateIndex 更新索引。
func (s *MemoryStore) UpdateIndex(i *model.Index) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.indexes[i.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.indexes {
		if exist.ID != i.ID && exist.Name == i.Name {
			return ErrConflict
		}
	}
	s.indexes[i.ID] = i
	return nil
}

// DeleteIndex 删除索引。
func (s *MemoryStore) DeleteIndex(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.indexes[id]; !ok {
		return ErrNotFound
	}
	delete(s.indexes, id)
	return nil
}
