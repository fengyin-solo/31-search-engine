package service

import (
	"sort"
	"time"

	"searchengine/internal/model"
	"searchengine/pkg/idgen"
)

// CreateAnalyzer 创建分词器。
func (s *Service) CreateAnalyzer(input model.Analyzer) (*model.Analyzer, error) {
	input.ID = idgen.Hex()
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateAnalyzer(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

// GetAnalyzer 获取分词器。
func (s *Service) GetAnalyzer(id string) (*model.Analyzer, error) {
	return s.store.GetAnalyzer(id)
}

// ListAnalyzers 分页列出分词器。
func (s *Service) ListAnalyzers(page, size int) ([]*model.Analyzer, int, error) {
	all := s.store.ListAnalyzers()
	sort.Slice(all, func(i, j int) bool { return all[i].CreatedAt.After(all[j].CreatedAt) })
	total := len(all)
	start := (page - 1) * size
	if start >= total {
		return []*model.Analyzer{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

// UpdateAnalyzer 更新分词器。
func (s *Service) UpdateAnalyzer(id string, input model.Analyzer) (*model.Analyzer, error) {
	existing, err := s.store.GetAnalyzer(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Type = input.Type
	existing.StopWords = input.StopWords
	existing.Status = input.Status
	existing.UpdatedAt = time.Now()
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateAnalyzer(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteAnalyzer 删除分词器。
func (s *Service) DeleteAnalyzer(id string) error {
	return s.store.DeleteAnalyzer(id)
}
