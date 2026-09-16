package service

import (
	"sort"

	"searchengine/internal/model"
)

// GetTerm 获取词条。
func (s *Service) GetTerm(id string) (*model.Term, error) {
	return s.store.GetTerm(id)
}

// ListTerms 分页列出词条（按文档频率倒序）。
func (s *Service) ListTerms(indexID string, page, size int) ([]*model.Term, int, error) {
	all := s.store.ListTerms()
	matched := make([]*model.Term, 0, len(all))
	for _, t := range all {
		if indexID == "" || t.IndexID == indexID {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].DocCount > matched[j].DocCount })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Term{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
