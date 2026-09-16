package service

import (
	"sort"
)

// OverviewStats 引擎总览统计。
type OverviewStats struct {
	DocumentCount int `json:"document_count"`
	IndexCount    int `json:"index_count"`
	AnalyzerCount int `json:"analyzer_count"`
	TermCount     int `json:"term_count"`
	PostingCount  int `json:"posting_count"`
	QueryCount    int `json:"query_count"`
	HotSearchCount int `json:"hot_search_count"`
	StopWordCount int `json:"stop_word_count"`
}

// Overview 计算引擎总览统计。
func (s *Service) Overview() OverviewStats {
	return OverviewStats{
		DocumentCount:  len(s.store.ListDocuments()),
		IndexCount:     len(s.store.ListIndexes()),
		AnalyzerCount:  len(s.store.ListAnalyzers()),
		TermCount:      len(s.store.ListTerms()),
		PostingCount:   len(s.store.ListPostings()),
		QueryCount:     len(s.store.ListQueryLogs()),
		HotSearchCount: len(s.store.ListHotSearches()),
		StopWordCount:  len(s.store.ListStopWords()),
	}
}

// IndexStatItem 索引统计条目。
type IndexStatItem struct {
	IndexID   string `json:"index_id"`
	Name      string `json:"name"`
	Status    string `json:"status"`
	DocCount  int    `json:"doc_count"`
	TermCount int    `json:"term_count"`
}

// IndexStats 统计每个索引的文档数与词条数。
func (s *Service) IndexStats() []IndexStatItem {
	termCounts := make(map[string]int)
	for _, t := range s.store.ListTerms() {
		termCounts[t.IndexID]++
	}
	items := make([]IndexStatItem, 0)
	for _, idx := range s.store.ListIndexes() {
		items = append(items, IndexStatItem{
			IndexID: idx.ID, Name: idx.Name, Status: idx.Status,
			DocCount: idx.DocCount, TermCount: termCounts[idx.ID],
		})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return items
}

// QueryStatItem 查询统计条目。
type QueryStatItem struct {
	Term  string `json:"term"`
	Count int    `json:"count"`
}

// TopQueryTerms 返回查询次数最高的 N 个查询词。
func (s *Service) TopQueryTerms(limit int) []QueryStatItem {
	agg := make(map[string]int)
	for _, q := range s.store.ListQueryLogs() {
		if q.Query != "" {
			agg[q.Query]++
		}
	}
	items := make([]QueryStatItem, 0, len(agg))
	for term, count := range agg {
		items = append(items, QueryStatItem{Term: term, Count: count})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Count > items[j].Count })
	if limit > 0 && len(items) > limit {
		items = items[:limit]
	}
	return items
}

// AvgQueryDuration 返回查询平均耗时（毫秒）。
func (s *Service) AvgQueryDuration() float64 {
	logs := s.store.ListQueryLogs()
	if len(logs) == 0 {
		return 0
	}
	var sum int
	for _, q := range logs {
		sum += q.DurationMs
	}
	return float64(sum) / float64(len(logs))
}

// TopTermItem 词条统计条目。
type TopTermItem struct {
	Term     string `json:"term"`
	DocCount int    `json:"doc_count"`
	TotalTF  int    `json:"total_tf"`
}

// TopTerms 返回文档频率最高的 N 个词条。
func (s *Service) TopTerms(limit int) []TopTermItem {
	all := s.store.ListTerms()
	sort.Slice(all, func(i, j int) bool { return all[i].DocCount > all[j].DocCount })
	items := make([]TopTermItem, 0, len(all))
	for _, t := range all {
		items = append(items, TopTermItem{Term: t.Term, DocCount: t.DocCount, TotalTF: t.TotalTF})
	}
	if limit > 0 && len(items) > limit {
		items = items[:limit]
	}
	return items
}

// HotSearchItem 热门搜索条目。
type HotSearchItem struct {
	Term  string `json:"term"`
	Count int    `json:"count"`
}

// TopHotSearches 返回热度最高的 N 个搜索词。
func (s *Service) TopHotSearches(limit int) []HotSearchItem {
	all := s.store.ListHotSearches()
	sort.Slice(all, func(i, j int) bool { return all[i].Count > all[j].Count })
	items := make([]HotSearchItem, 0, len(all))
	for _, h := range all {
		items = append(items, HotSearchItem{Term: h.Term, Count: h.Count})
	}
	if limit > 0 && len(items) > limit {
		items = items[:limit]
	}
	return items
}

// CategoryStatItem 分类统计条目。
type CategoryStatItem struct {
	Name     string `json:"name"`
	DocCount int    `json:"doc_count"`
}

// CategoryDocStats 按分类统计文档数。
func (s *Service) CategoryDocStats() []CategoryStatItem {
	items := make([]CategoryStatItem, 0)
	for _, c := range s.store.ListCategories() {
		items = append(items, CategoryStatItem{Name: c.Name, DocCount: c.DocCount})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].DocCount > items[j].DocCount })
	return items
}

// SourceDocStats 按来源统计文档数。
func (s *Service) SourceDocStats() map[string]int {
	result := make(map[string]int)
	for _, d := range s.store.ListDocuments() {
		if d.Source == "" {
			result["未知"]++
		} else {
			result[d.Source]++
		}
	}
	return result
}
