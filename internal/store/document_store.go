package store

import "searchengine/internal/model"

// CreateDocument 创建文档。
func (s *MemoryStore) CreateDocument(d *model.Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.documents[d.ID] = d
	return nil
}

// GetDocument 按 ID 获取文档。
func (s *MemoryStore) GetDocument(id string) (*model.Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.documents[id]
	if !ok {
		return nil, ErrNotFound
	}
	return d, nil
}

// ListDocuments 列出全部文档。
func (s *MemoryStore) ListDocuments() []*model.Document {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Document, 0, len(s.documents))
	for _, d := range s.documents {
		list = append(list, d)
	}
	return list
}

// UpdateDocument 更新文档。
func (s *MemoryStore) UpdateDocument(d *model.Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.documents[d.ID]; !ok {
		return ErrNotFound
	}
	s.documents[d.ID] = d
	return nil
}

// DeleteDocument 删除文档。
func (s *MemoryStore) DeleteDocument(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.documents[id]; !ok {
		return ErrNotFound
	}
	delete(s.documents, id)
	return nil
}
