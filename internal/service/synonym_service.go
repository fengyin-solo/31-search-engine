package service

import (
	"sort"
	"time"

	"searchengine/internal/model"
	"searchengine/pkg/idgen"
)

// CreateSynonym 创建同义词。
func (s *Service) CreateSynonym(input model.Synonym) (*model.Synonym, error) {
	input.ID = idgen.Hex()
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateSynonym(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

// GetSynonym 获取同义词。
func (s *Service) GetSynonym(id string) (*model.Synonym, error) {
	return s.store.GetSynonym(id)
}

// ListSynonyms 分页列出同义词。
func (s *Service) ListSynonyms(page, size int) ([]*model.Synonym, int, error) {
	all := s.store.ListSynonyms()
	sort.Slice(all, func(i, j int) bool { return all[i].Word < all[j].Word })
	total := len(all)
	start := (page - 1) * size
	if start >= total {
		return []*model.Synonym{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

// UpdateSynonym 更新同义词。
func (s *Service) UpdateSynonym(id string, input model.Synonym) (*model.Synonym, error) {
	existing, err := s.store.GetSynonym(id)
	if err != nil {
		return nil, err
	}
	existing.Word = input.Word
	existing.Synonyms = input.Synonyms
	existing.Status = input.Status
	existing.UpdatedAt = time.Now()
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateSynonym(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteSynonym 删除同义词。
func (s *Service) DeleteSynonym(id string) error {
	return s.store.DeleteSynonym(id)
}

// synonymMap 返回 原词→同义词列表 的映射（供查询扩展使用）。
func (s *Service) synonymMap() map[string][]string {
	m := make(map[string][]string)
	for _, sy := range s.store.ListSynonyms() {
		if sy.Status != model.SynonymActive {
			continue
		}
		m[sy.Word] = sy.Synonyms
	}
	return m
}
