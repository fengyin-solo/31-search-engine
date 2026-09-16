package service

import (
	"sort"
	"strings"
	"time"

	"searchengine/internal/model"
	"searchengine/pkg/idgen"
)

// CreateSuggestion 创建建议词。
func (s *Service) CreateSuggestion(input model.Suggestion) (*model.Suggestion, error) {
	input.ID = idgen.Hex()
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateSuggestion(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

// GetSuggestion 获取建议词。
func (s *Service) GetSuggestion(id string) (*model.Suggestion, error) {
	return s.store.GetSuggestion(id)
}

// ListSuggestions 分页列出建议词。
func (s *Service) ListSuggestions(page, size int) ([]*model.Suggestion, int, error) {
	all := s.store.ListSuggestions()
	sort.Slice(all, func(i, j int) bool { return all[i].Weight > all[j].Weight })
	total := len(all)
	start := (page - 1) * size
	if start >= total {
		return []*model.Suggestion{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

// UpdateSuggestion 更新建议词。
func (s *Service) UpdateSuggestion(id string, input model.Suggestion) (*model.Suggestion, error) {
	existing, err := s.store.GetSuggestion(id)
	if err != nil {
		return nil, err
	}
	existing.Term = input.Term
	existing.Weight = input.Weight
	existing.UpdatedAt = time.Now()
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateSuggestion(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteSuggestion 删除建议词。
func (s *Service) DeleteSuggestion(id string) error {
	return s.store.DeleteSuggestion(id)
}

// Suggest 根据前缀返回搜索建议（综合建议词库与词典词条）。
func (s *Service) Suggest(prefix string, limit int) []string {
	prefix = strings.ToLower(strings.TrimSpace(prefix))
	if prefix == "" {
		return []string{}
	}
	seen := make(map[string]bool)
	result := make([]string, 0)

	candidates := make(map[string]int)
	for _, sug := range s.store.ListSuggestions() {
		if strings.HasPrefix(sug.Term, prefix) {
			candidates[sug.Term] = sug.Weight
		}
	}
	for _, t := range s.store.ListTerms() {
		if strings.HasPrefix(t.Term, prefix) {
			if _, ok := candidates[t.Term]; !ok {
				candidates[t.Term] = t.DocCount
			}
		}
	}

	type pair struct {
		term   string
		weight int
	}
	pairs := make([]pair, 0, len(candidates))
	for term, w := range candidates {
		pairs = append(pairs, pair{term: term, weight: w})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].weight > pairs[j].weight })
	for _, p := range pairs {
		if seen[p.term] {
			continue
		}
		seen[p.term] = true
		result = append(result, p.term)
		if limit > 0 && len(result) >= limit {
			break
		}
	}
	return result
}
