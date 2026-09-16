package service

import (
	"sort"

	"searchengine/internal/model"
)

// GetQueryLog 获取查询记录。
func (s *Service) GetQueryLog(id string) (*model.QueryLog, error) {
	return s.store.GetQueryLog(id)
}

// ListQueryLogs 分页列出查询记录。
func (s *Service) ListQueryLogs(indexID string, page, size int) ([]*model.QueryLog, int, error) {
	all := s.store.ListQueryLogs()
	matched := make([]*model.QueryLog, 0, len(all))
	for _, q := range all {
		if indexID == "" || q.IndexID == indexID {
			matched = append(matched, q)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.After(matched[j].CreatedAt) })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.QueryLog{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// DeleteQueryLog 删除查询记录。
func (s *Service) DeleteQueryLog(id string) error {
	return s.store.DeleteQueryLog(id)
}
