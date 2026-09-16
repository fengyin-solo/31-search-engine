package service

import "searchengine/internal/model"

// ListAllDocumentsForExport 导出全部文档。
func (s *Service) ListAllDocumentsForExport() []*model.Document { return s.store.ListDocuments() }

// ListAllIndexesForExport 导出全部索引。
func (s *Service) ListAllIndexesForExport() []*model.Index { return s.store.ListIndexes() }

// ListAllAnalyzersForExport 导出全部分词器。
func (s *Service) ListAllAnalyzersForExport() []*model.Analyzer { return s.store.ListAnalyzers() }

// ListAllStopWordsForExport 导出全部停用词。
func (s *Service) ListAllStopWordsForExport() []*model.StopWord { return s.store.ListStopWords() }

// ListAllTermsForExport 导出全部词条。
func (s *Service) ListAllTermsForExport() []*model.Term { return s.store.ListTerms() }

// ListAllPostingsForExport 导出全部倒排记录。
func (s *Service) ListAllPostingsForExport() []*model.Posting { return s.store.ListPostings() }

// ListAllSynonymsForExport 导出全部同义词。
func (s *Service) ListAllSynonymsForExport() []*model.Synonym { return s.store.ListSynonyms() }

// ListAllQueryLogsForExport 导出全部查询记录。
func (s *Service) ListAllQueryLogsForExport() []*model.QueryLog { return s.store.ListQueryLogs() }

// ListAllHotSearchesForExport 导出全部热门搜索词。
func (s *Service) ListAllHotSearchesForExport() []*model.HotSearch { return s.store.ListHotSearches() }

// ListAllCategoriesForExport 导出全部分类。
func (s *Service) ListAllCategoriesForExport() []*model.Category { return s.store.ListCategories() }

// ListAllSuggestionsForExport 导出全部建议词。
func (s *Service) ListAllSuggestionsForExport() []*model.Suggestion { return s.store.ListSuggestions() }

// ListAllSpellCorrectionsForExport 导出全部纠错词对。
func (s *Service) ListAllSpellCorrectionsForExport() []*model.SpellCorrection {
	return s.store.ListSpellCorrections()
}

// ListAllBoostRulesForExport 导出全部加权规则。
func (s *Service) ListAllBoostRulesForExport() []*model.BoostRule { return s.store.ListBoostRules() }
