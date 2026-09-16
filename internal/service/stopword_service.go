package service

import (
	"sort"
	"time"

	"searchengine/internal/model"
	"searchengine/pkg/idgen"
)

// CreateStopWord 创建停用词。
func (s *Service) CreateStopWord(input model.StopWord) (*model.StopWord, error) {
	input.ID = idgen.Hex()
	input.CreatedAt = time.Now()
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateStopWord(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

// GetStopWord 获取停用词。
func (s *Service) GetStopWord(id string) (*model.StopWord, error) {
	return s.store.GetStopWord(id)
}

// ListStopWords 分页列出停用词。
func (s *Service) ListStopWords(page, size int) ([]*model.StopWord, int, error) {
	all := s.store.ListStopWords()
	sort.Slice(all, func(i, j int) bool { return all[i].Word < all[j].Word })
	total := len(all)
	start := (page - 1) * size
	if start >= total {
		return []*model.StopWord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

// DeleteStopWord 删除停用词。
func (s *Service) DeleteStopWord(id string) error {
	return s.store.DeleteStopWord(id)
}

// stopWordSet 返回全部停用词集合（供分词使用）。
func (s *Service) stopWordSet() map[string]bool {
	set := make(map[string]bool)
	for _, w := range s.store.ListStopWords() {
		set[w.Word] = true
	}
	return set
}
