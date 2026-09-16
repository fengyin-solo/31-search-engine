package store

import "searchengine/internal/model"

// CreateStopWord 创建停用词。
func (s *MemoryStore) CreateStopWord(w *model.StopWord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.stopWords {
		if exist.Word == w.Word {
			return ErrConflict
		}
	}
	s.stopWords[w.ID] = w
	return nil
}

// GetStopWord 按 ID 获取停用词。
func (s *MemoryStore) GetStopWord(id string) (*model.StopWord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	w, ok := s.stopWords[id]
	if !ok {
		return nil, ErrNotFound
	}
	return w, nil
}

// GetStopWordByWord 按词获取停用词。
func (s *MemoryStore) GetStopWordByWord(word string) (*model.StopWord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, w := range s.stopWords {
		if w.Word == word {
			return w, nil
		}
	}
	return nil, ErrNotFound
}

// ListStopWords 列出全部停用词。
func (s *MemoryStore) ListStopWords() []*model.StopWord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.StopWord, 0, len(s.stopWords))
	for _, w := range s.stopWords {
		list = append(list, w)
	}
	return list
}

// DeleteStopWord 删除停用词。
func (s *MemoryStore) DeleteStopWord(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.stopWords[id]; !ok {
		return ErrNotFound
	}
	delete(s.stopWords, id)
	return nil
}
