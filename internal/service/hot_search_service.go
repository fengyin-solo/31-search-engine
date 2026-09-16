package service

import (
	"sort"

	"searchengine/internal/model"
)

// GetHotSearch 获取热门搜索词。
func (s *Service) GetHotSearch(id string) (*model.HotSearch, error) {
	return s.store.GetHotSearch(id)
}

// ListHotSearches 分页列出热门搜索词（按热度倒序）。
func (s *Service) ListHotSearches(page, size int) ([]*model.HotSearch, int, error) {
	all := s.store.ListHotSearches()
	sort.Slice(all, func(i, j int) bool { return all[i].Count > all[j].Count })
	total := len(all)
	start := (page - 1) * size
	if start >= total {
		return []*model.HotSearch{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

// DeleteHotSearch 删除热门搜索词。
func (s *Service) DeleteHotSearch(id string) error {
	return s.store.DeleteHotSearch(id)
}
