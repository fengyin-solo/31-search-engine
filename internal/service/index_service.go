package service

import (
	"sort"
	"time"

	"searchengine/internal/model"
	"searchengine/internal/store"
	"searchengine/pkg/idgen"
)

// CreateIndex 创建索引。
func (s *Service) CreateIndex(input model.Index) (*model.Index, error) {
	input.ID = idgen.Hex()
	input.Status = model.IndexCreated
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateIndex(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

// GetIndex 获取索引。
func (s *Service) GetIndex(id string) (*model.Index, error) {
	return s.store.GetIndex(id)
}

// ListIndexes 分页列出索引。
func (s *Service) ListIndexes(page, size int) ([]*model.Index, int, error) {
	all := s.store.ListIndexes()
	sort.Slice(all, func(i, j int) bool { return all[i].CreatedAt.After(all[j].CreatedAt) })
	total := len(all)
	start := (page - 1) * size
	if start >= total {
		return []*model.Index{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

// UpdateIndex 更新索引信息。
func (s *Service) UpdateIndex(id string, input model.Index) (*model.Index, error) {
	existing, err := s.store.GetIndex(id)
	if err != nil {
		return nil, err
	}
	existing.Name = input.Name
	existing.Description = input.Description
	existing.Fields = input.Fields
	existing.AnalyzerID = input.AnalyzerID
	existing.UpdatedAt = time.Now()
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateIndex(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// ActivateIndex 激活索引（created→ready）。
func (s *Service) ActivateIndex(id string) (*model.Index, error) {
	return s.transitionIndex(id, model.IndexReady)
}

// DeactivateIndex 停用索引（ready→deleted）。
func (s *Service) DeactivateIndex(id string) (*model.Index, error) {
	return s.transitionIndex(id, model.IndexDeleted)
}

func (s *Service) transitionIndex(id, to string) (*model.Index, error) {
	existing, err := s.store.GetIndex(id)
	if err != nil {
		return nil, err
	}
	if !model.CanIndexTransition(existing.Status, to) {
		return nil, store.ErrConflict
	}
	existing.Status = to
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateIndex(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteIndex 删除索引。
func (s *Service) DeleteIndex(id string) error {
	return s.store.DeleteIndex(id)
}
