package service

import (
	"sort"

	"searchengine/internal/model"
)

// GetPosting 获取倒排记录。
func (s *Service) GetPosting(id string) (*model.Posting, error) {
	return s.store.GetPosting(id)
}

// ListPostings 分页列出倒排记录。
func (s *Service) ListPostings(indexID, term, docID string, page, size int) ([]*model.Posting, int, error) {
	all := s.store.ListPostings()
	matched := make([]*model.Posting, 0, len(all))
	for _, p := range all {
		if indexID != "" && p.IndexID != indexID {
			continue
		}
		if term != "" && p.Term != term {
			continue
		}
		if docID != "" && p.DocID != docID {
			continue
		}
		matched = append(matched, p)
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].Term < matched[j].Term })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Posting{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}
