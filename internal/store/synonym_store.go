package store

import "searchengine/internal/model"

// CreateSynonym 创建同义词。
func (s *MemoryStore) CreateSynonym(sy *model.Synonym) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.synonyms {
		if exist.Word == sy.Word {
			return ErrConflict
		}
	}
	s.synonyms[sy.ID] = sy
	return nil
}

// GetSynonym 按 ID 获取同义词。
func (s *MemoryStore) GetSynonym(id string) (*model.Synonym, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sy, ok := s.synonyms[id]
	if !ok {
		return nil, ErrNotFound
	}
	return sy, nil
}

// GetSynonymByWord 按原词获取同义词。
func (s *MemoryStore) GetSynonymByWord(word string) (*model.Synonym, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, sy := range s.synonyms {
		if sy.Word == word {
			return sy, nil
		}
	}
	return nil, ErrNotFound
}

// ListSynonyms 列出全部同义词。
func (s *MemoryStore) ListSynonyms() []*model.Synonym {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Synonym, 0, len(s.synonyms))
	for _, sy := range s.synonyms {
		list = append(list, sy)
	}
	return list
}

// UpdateSynonym 更新同义词。
func (s *MemoryStore) UpdateSynonym(sy *model.Synonym) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.synonyms[sy.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.synonyms {
		if exist.ID != sy.ID && exist.Word == sy.Word {
			return ErrConflict
		}
	}
	s.synonyms[sy.ID] = sy
	return nil
}

// DeleteSynonym 删除同义词。
func (s *MemoryStore) DeleteSynonym(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.synonyms[id]; !ok {
		return ErrNotFound
	}
	delete(s.synonyms, id)
	return nil
}
