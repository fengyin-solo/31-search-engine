package store

import "searchengine/internal/model"

// CreateSpellCorrection 创建拼写纠错词对。
func (s *MemoryStore) CreateSpellCorrection(v *model.SpellCorrection) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.spellCorrections {
		if exist.Word == v.Word {
			return ErrConflict
		}
	}
	s.spellCorrections[v.ID] = v
	return nil
}

// GetSpellCorrection 按 ID 获取纠错词对。
func (s *MemoryStore) GetSpellCorrection(id string) (*model.SpellCorrection, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.spellCorrections[id]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

// GetSpellCorrectionByWord 按错误拼写获取纠错词对。
func (s *MemoryStore) GetSpellCorrectionByWord(word string) (*model.SpellCorrection, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, v := range s.spellCorrections {
		if v.Word == word {
			return v, nil
		}
	}
	return nil, ErrNotFound
}

// ListSpellCorrections 列出全部纠错词对。
func (s *MemoryStore) ListSpellCorrections() []*model.SpellCorrection {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.SpellCorrection, 0, len(s.spellCorrections))
	for _, v := range s.spellCorrections {
		list = append(list, v)
	}
	return list
}

// UpdateSpellCorrection 更新纠错词对。
func (s *MemoryStore) UpdateSpellCorrection(v *model.SpellCorrection) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.spellCorrections[v.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.spellCorrections {
		if exist.ID != v.ID && exist.Word == v.Word {
			return ErrConflict
		}
	}
	s.spellCorrections[v.ID] = v
	return nil
}

// DeleteSpellCorrection 删除纠错词对。
func (s *MemoryStore) DeleteSpellCorrection(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.spellCorrections[id]; !ok {
		return ErrNotFound
	}
	delete(s.spellCorrections, id)
	return nil
}
