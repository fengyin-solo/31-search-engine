package store

import "searchengine/internal/model"

// CreateCategory 创建分类。
func (s *MemoryStore) CreateCategory(c *model.Category) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.categories {
		if exist.Name == c.Name {
			return ErrConflict
		}
	}
	s.categories[c.ID] = c
	return nil
}

// GetCategory 按 ID 获取分类。
func (s *MemoryStore) GetCategory(id string) (*model.Category, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.categories[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

// GetCategoryByName 按名称获取分类。
func (s *MemoryStore) GetCategoryByName(name string) (*model.Category, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.categories {
		if c.Name == name {
			return c, nil
		}
	}
	return nil, ErrNotFound
}

// ListCategories 列出全部分类。
func (s *MemoryStore) ListCategories() []*model.Category {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Category, 0, len(s.categories))
	for _, c := range s.categories {
		list = append(list, c)
	}
	return list
}

// UpdateCategory 更新分类。
func (s *MemoryStore) UpdateCategory(c *model.Category) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.categories[c.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.categories {
		if exist.ID != c.ID && exist.Name == c.Name {
			return ErrConflict
		}
	}
	s.categories[c.ID] = c
	return nil
}

// DeleteCategory 删除分类。
func (s *MemoryStore) DeleteCategory(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.categories[id]; !ok {
		return ErrNotFound
	}
	delete(s.categories, id)
	return nil
}
