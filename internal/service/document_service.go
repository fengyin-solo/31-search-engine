package service

import (
	"sort"
	"time"

	"searchengine/internal/model"
	"searchengine/pkg/idgen"
)

// CreateDocument 创建文档。
func (s *Service) CreateDocument(input model.Document) (*model.Document, error) {
	input.ID = idgen.Hex()
	input.Status = model.DocumentActive
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateDocument(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

// GetDocument 获取文档。
func (s *Service) GetDocument(id string) (*model.Document, error) {
	return s.store.GetDocument(id)
}

// ListDocuments 分页列出文档。
func (s *Service) ListDocuments(filter model.DocumentFilter, page, size int) ([]*model.Document, int, error) {
	all := s.store.ListDocuments()
	matched := make([]*model.Document, 0, len(all))
	for _, d := range all {
		if filter.Match(d) {
			matched = append(matched, d)
		}
	}
	sort.Slice(matched, func(i, j int) bool { return matched[i].CreatedAt.After(matched[j].CreatedAt) })
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Document{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

// UpdateDocument 更新文档内容。
func (s *Service) UpdateDocument(id string, input model.Document) (*model.Document, error) {
	existing, err := s.store.GetDocument(id)
	if err != nil {
		return nil, err
	}
	existing.Title = input.Title
	existing.Body = input.Body
	existing.URL = input.URL
	existing.Source = input.Source
	existing.Status = input.Status
	existing.UpdatedAt = time.Now()
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateDocument(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteDocument 删除文档（软删，并从所有索引移除）。
func (s *Service) DeleteDocument(id string) error {
	existing, err := s.store.GetDocument(id)
	if err != nil {
		return err
	}
	existing.Status = model.DocumentDeleted
	existing.UpdatedAt = time.Now()
	if err := s.store.UpdateDocument(existing); err != nil {
		return err
	}
	// 从所有索引中移除该文档的倒排记录
	for _, idx := range s.store.ListIndexes() {
		_ = s.removeDocumentFromIndex(idx.ID, id)
	}
	return nil
}
