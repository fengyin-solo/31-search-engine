// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"searchengine/internal/model"
)

var (
	// ErrNotFound 表示记录不存在。
	ErrNotFound = errors.New("记录不存在")
	// ErrConflict 表示记录已存在或状态冲突。
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法。
type Store interface {
	CreateDocument(d *model.Document) error
	GetDocument(id string) (*model.Document, error)
	ListDocuments() []*model.Document
	UpdateDocument(d *model.Document) error
	DeleteDocument(id string) error

	CreateIndex(i *model.Index) error
	GetIndex(id string) (*model.Index, error)
	GetIndexByName(name string) (*model.Index, error)
	ListIndexes() []*model.Index
	UpdateIndex(i *model.Index) error
	DeleteIndex(id string) error

	CreateAnalyzer(a *model.Analyzer) error
	GetAnalyzer(id string) (*model.Analyzer, error)
	GetAnalyzerByName(name string) (*model.Analyzer, error)
	ListAnalyzers() []*model.Analyzer
	UpdateAnalyzer(a *model.Analyzer) error
	DeleteAnalyzer(id string) error

	CreateStopWord(s *model.StopWord) error
	GetStopWord(id string) (*model.StopWord, error)
	GetStopWordByWord(word string) (*model.StopWord, error)
	ListStopWords() []*model.StopWord
	DeleteStopWord(id string) error

	CreateTerm(t *model.Term) error
	GetTerm(id string) (*model.Term, error)
	GetTermByKey(indexID, term string) (*model.Term, error)
	ListTerms() []*model.Term
	UpdateTerm(t *model.Term) error
	DeleteTerm(id string) error

	CreatePosting(p *model.Posting) error
	GetPosting(id string) (*model.Posting, error)
	GetPostingByKey(indexID, term, docID string) (*model.Posting, error)
	ListPostings() []*model.Posting
	UpdatePosting(p *model.Posting) error
	DeletePosting(id string) error

	CreateSynonym(s *model.Synonym) error
	GetSynonym(id string) (*model.Synonym, error)
	GetSynonymByWord(word string) (*model.Synonym, error)
	ListSynonyms() []*model.Synonym
	UpdateSynonym(s *model.Synonym) error
	DeleteSynonym(id string) error

	CreateQueryLog(q *model.QueryLog) error
	GetQueryLog(id string) (*model.QueryLog, error)
	ListQueryLogs() []*model.QueryLog
	DeleteQueryLog(id string) error

	CreateHotSearch(h *model.HotSearch) error
	GetHotSearch(id string) (*model.HotSearch, error)
	GetHotSearchByTerm(term string) (*model.HotSearch, error)
	ListHotSearches() []*model.HotSearch
	UpdateHotSearch(h *model.HotSearch) error
	DeleteHotSearch(id string) error

	CreateCategory(c *model.Category) error
	GetCategory(id string) (*model.Category, error)
	GetCategoryByName(name string) (*model.Category, error)
	ListCategories() []*model.Category
	UpdateCategory(c *model.Category) error
	DeleteCategory(id string) error

	CreateSuggestion(s *model.Suggestion) error
	GetSuggestion(id string) (*model.Suggestion, error)
	GetSuggestionByTerm(term string) (*model.Suggestion, error)
	ListSuggestions() []*model.Suggestion
	UpdateSuggestion(s *model.Suggestion) error
	DeleteSuggestion(id string) error

	CreateSpellCorrection(s *model.SpellCorrection) error
	GetSpellCorrection(id string) (*model.SpellCorrection, error)
	GetSpellCorrectionByWord(word string) (*model.SpellCorrection, error)
	ListSpellCorrections() []*model.SpellCorrection
	UpdateSpellCorrection(s *model.SpellCorrection) error
	DeleteSpellCorrection(id string) error

	CreateBoostRule(b *model.BoostRule) error
	GetBoostRule(id string) (*model.BoostRule, error)
	ListBoostRules() []*model.BoostRule
	UpdateBoostRule(b *model.BoostRule) error
	DeleteBoostRule(id string) error
}
