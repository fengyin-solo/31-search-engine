package store

import (
	"sync"

	"searchengine/internal/model"
)

// MemoryStore 基于内存的 Store 实现。
type MemoryStore struct {
	mu sync.RWMutex

	documents   map[string]*model.Document
	indexes     map[string]*model.Index
	analyzers   map[string]*model.Analyzer
	stopWords   map[string]*model.StopWord
	terms       map[string]*model.Term
	postings    map[string]*model.Posting
	synonyms    map[string]*model.Synonym
	queryLogs   map[string]*model.QueryLog
	hotSearches map[string]*model.HotSearch

	categories       map[string]*model.Category
	suggestions      map[string]*model.Suggestion
	spellCorrections map[string]*model.SpellCorrection
	boostRules       map[string]*model.BoostRule
}

// NewMemoryStore 构造内存存储。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		documents:   make(map[string]*model.Document),
		indexes:     make(map[string]*model.Index),
		analyzers:   make(map[string]*model.Analyzer),
		stopWords:   make(map[string]*model.StopWord),
		terms:       make(map[string]*model.Term),
		postings:    make(map[string]*model.Posting),
		synonyms:    make(map[string]*model.Synonym),
		queryLogs:   make(map[string]*model.QueryLog),
		hotSearches: make(map[string]*model.HotSearch),

		categories:       make(map[string]*model.Category),
		suggestions:      make(map[string]*model.Suggestion),
		spellCorrections: make(map[string]*model.SpellCorrection),
		boostRules:       make(map[string]*model.BoostRule),
	}
}

var _ Store = (*MemoryStore)(nil)
