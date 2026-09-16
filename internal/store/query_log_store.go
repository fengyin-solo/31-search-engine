package store

import "searchengine/internal/model"

// CreateQueryLog 创建查询记录。
func (s *MemoryStore) CreateQueryLog(q *model.QueryLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queryLogs[q.ID] = q
	return nil
}

// GetQueryLog 按 ID 获取查询记录。
func (s *MemoryStore) GetQueryLog(id string) (*model.QueryLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	q, ok := s.queryLogs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return q, nil
}

// ListQueryLogs 列出全部查询记录。
func (s *MemoryStore) ListQueryLogs() []*model.QueryLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.QueryLog, 0, len(s.queryLogs))
	for _, q := range s.queryLogs {
		list = append(list, q)
	}
	return list
}

// DeleteQueryLog 删除查询记录。
func (s *MemoryStore) DeleteQueryLog(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.queryLogs[id]; !ok {
		return ErrNotFound
	}
	delete(s.queryLogs, id)
	return nil
}
