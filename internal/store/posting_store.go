package store

import "searchengine/internal/model"

// CreatePosting 创建倒排记录。
func (s *MemoryStore) CreatePosting(p *model.Posting) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.postings[p.ID] = p
	return nil
}

// GetPosting 按 ID 获取倒排记录。
func (s *MemoryStore) GetPosting(id string) (*model.Posting, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.postings[id]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

// GetPostingByKey 按 索引+词+文档 获取倒排记录。
func (s *MemoryStore) GetPostingByKey(indexID, term, docID string) (*model.Posting, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.postings {
		if p.IndexID == indexID && p.Term == term && p.DocID == docID {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

// ListPostings 列出全部倒排记录。
func (s *MemoryStore) ListPostings() []*model.Posting {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Posting, 0, len(s.postings))
	for _, p := range s.postings {
		list = append(list, p)
	}
	return list
}

// UpdatePosting 更新倒排记录。
func (s *MemoryStore) UpdatePosting(p *model.Posting) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.postings[p.ID]; !ok {
		return ErrNotFound
	}
	s.postings[p.ID] = p
	return nil
}

// DeletePosting 删除倒排记录。
func (s *MemoryStore) DeletePosting(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.postings[id]; !ok {
		return ErrNotFound
	}
	delete(s.postings, id)
	return nil
}
