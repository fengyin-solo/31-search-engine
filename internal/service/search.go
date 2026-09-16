package service

import (
	"math"
	"sort"
	"strings"
	"time"

	"searchengine/internal/model"
	"searchengine/pkg/idgen"
)

// SearchResult 搜索结果条目。
type SearchResult struct {
	DocID          string  `json:"doc_id"`
	Title          string  `json:"title"`
	HighlightTitle string  `json:"highlight_title"`
	Body           string  `json:"body"`
	CategoryID     string  `json:"category_id"`
	Score          float64 `json:"score"`
}

// IndexDocument 将文档索引到指定索引（先移除旧索引再重建）。
func (s *Service) IndexDocument(indexID, docID string) error {
	idx, err := s.store.GetIndex(indexID)
	if err != nil {
		return err
	}
	if idx.Status != model.IndexReady && idx.Status != model.IndexCreated {
		return model.NewValidationError("index", "索引未就绪")
	}
	doc, err := s.store.GetDocument(docID)
	if err != nil {
		return err
	}
	if doc.Status != model.DocumentActive {
		return model.NewValidationError("document", "文档已被删除")
	}
	// 先移除该文档在该索引的旧倒排记录，避免重复计数
	_ = s.removeDocumentFromIndex(indexID, docID)

	stopwords := s.stopWordSet()
	if idx.AnalyzerID != "" {
		if a, err := s.store.GetAnalyzer(idx.AnalyzerID); err == nil {
			stopwords = a.StopWordSet()
		}
	}
	tokens := Tokenize(doc.Title+" "+doc.Body, stopwords)

	counts := make(map[string]int)
	for _, t := range tokens {
		counts[t]++
	}

	for term, tf := range counts {
		posting := &model.Posting{
			ID: idgen.Hex(), IndexID: indexID, Term: term, DocID: docID, TF: tf,
		}
		if err := s.store.CreatePosting(posting); err != nil {
			return err
		}
		termObj, err := s.store.GetTermByKey(indexID, term)
		if err != nil {
			termObj = &model.Term{ID: idgen.Hex(), IndexID: indexID, Term: term}
			if err := s.store.CreateTerm(termObj); err != nil {
				return err
			}
		}
		termObj.DocCount++
		termObj.TotalTF += tf
		if err := s.store.UpdateTerm(termObj); err != nil {
			return err
		}
	}
	idx.DocCount++
	idx.UpdatedAt = time.Now()
	return s.store.UpdateIndex(idx)
}

// removeDocumentFromIndex 从索引中移除指定文档的全部倒排记录。
func (s *Service) removeDocumentFromIndex(indexID, docID string) error {
	var toRemove []*model.Posting
	for _, p := range s.store.ListPostings() {
		if p.IndexID == indexID && p.DocID == docID {
			toRemove = append(toRemove, p)
		}
	}
	if len(toRemove) == 0 {
		return nil
	}
	for _, p := range toRemove {
		if t, err := s.store.GetTermByKey(indexID, p.Term); err == nil {
			t.DocCount--
			t.TotalTF -= p.TF
			if t.DocCount <= 0 {
				_ = s.store.DeleteTerm(t.ID)
			} else {
				_ = s.store.UpdateTerm(t)
			}
		}
		_ = s.store.DeletePosting(p.ID)
	}
	if idx, err := s.store.GetIndex(indexID); err == nil {
		if idx.DocCount > 0 {
			idx.DocCount--
		}
		idx.UpdatedAt = time.Now()
		_ = s.store.UpdateIndex(idx)
	}
	return nil
}

// Search 在索引中搜索查询词，按 TF·IDF 相关度排序返回结果。
func (s *Service) Search(indexID, query string, topK int) ([]SearchResult, int, error) {
	start := time.Now()
	idx, err := s.store.GetIndex(indexID)
	if err != nil {
		return nil, 0, err
	}
	stopwords := s.stopWordSet()
	if idx.AnalyzerID != "" {
		if a, err := s.store.GetAnalyzer(idx.AnalyzerID); err == nil {
			stopwords = a.StopWordSet()
		}
	}
	tokens := Tokenize(query, stopwords)
	tokens = s.expandSynonyms(tokens)

	tokenSet := make(map[string]bool)
	for _, t := range tokens {
		tokenSet[t] = true
	}

	totalDocs := len(s.store.ListDocuments())
	N := float64(totalDocs)
	if N == 0 {
		N = 1
	}
	scores := make(map[string]float64)
	for term := range tokenSet {
		for _, p := range s.store.ListPostings() {
			if p.IndexID != indexID || p.Term != term {
				continue
			}
			idf := 1.0
			if termObj, err := s.store.GetTermByKey(indexID, term); err == nil && termObj.DocCount > 0 {
				idf = math.Log(1 + N/(1+float64(termObj.DocCount)))
			}
			scores[p.DocID] += float64(p.TF) * idf
		}
	}

	results := make([]SearchResult, 0, len(scores))
	for docID, score := range scores {
		doc, err := s.store.GetDocument(docID)
		if err != nil || doc.Status != model.DocumentActive {
			continue
		}
		score = s.applyBoost(doc, score)
		results = append(results, SearchResult{
			DocID: docID, Title: doc.Title, HighlightTitle: highlight(doc.Title, tokenSet),
			Body: truncateText(doc.Body, 200), CategoryID: doc.CategoryID, Score: score,
		})
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		return results[i].DocID < results[j].DocID
	})
	total := len(results)
	if topK > 0 && len(results) > topK {
		results = results[:topK]
	}

	duration := int(time.Since(start).Milliseconds())
	s.logQuery(indexID, query, total, duration)
	s.bumpHotSearch(query)
	return results, total, nil
}

// expandSynonyms 用同义词扩展查询词。
func (s *Service) expandSynonyms(tokens []string) []string {
	m := s.synonymMap()
	if len(m) == 0 {
		return tokens
	}
	result := make([]string, 0, len(tokens)*2)
	for _, t := range tokens {
		result = append(result, t)
		if syns, ok := m[t]; ok {
			result = append(result, syns...)
		}
	}
	return result
}

// logQuery 记录一次查询。
func (s *Service) logQuery(indexID, query string, resultCount, duration int) {
	q := &model.QueryLog{
		ID: idgen.Hex(), IndexID: indexID, Query: query,
		ResultCount: resultCount, DurationMs: duration, CreatedAt: time.Now(),
	}
	_ = s.store.CreateQueryLog(q)
}

// bumpHotSearch 增加搜索词热度。
func (s *Service) bumpHotSearch(term string) {
	if term == "" {
		return
	}
	h, err := s.store.GetHotSearchByTerm(term)
	if err != nil {
		h = &model.HotSearch{ID: idgen.Hex(), Term: term, UpdatedAt: time.Now()}
		_ = s.store.CreateHotSearch(h)
	}
	h.Count++
	h.UpdatedAt = time.Now()
	_ = s.store.UpdateHotSearch(h)
}

// truncateText 截断文本到指定长度。
func truncateText(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "..."
}

// applyBoost 根据加权规则调整文档相关度得分。
func (s *Service) applyBoost(doc *model.Document, score float64) float64 {
	for _, rule := range s.activeBoostRules() {
		var haystack string
		if rule.Field == model.BoostFieldTitle {
			haystack = doc.Title
		} else {
			haystack = doc.Body
		}
		if strings.Contains(strings.ToLower(haystack), rule.Term) {
			score *= rule.Boost
		}
	}
	return score
}

// highlight 将标题中命中的查询词用 <em> 包裹。
func highlight(text string, terms map[string]bool) string {
	lower := strings.ToLower(text)
	for term := range terms {
		if term == "" {
			continue
		}
		idx := 0
		for {
			pos := strings.Index(lower[idx:], term)
			if pos < 0 {
				break
			}
			abs := idx + pos
			text = text[:abs] + "<em>" + text[abs:abs+len(term)] + "</em>" + text[abs+len(term):]
			lower = strings.ToLower(text)
			idx = abs + len(term) + len("<em></em>")
		}
	}
	return text
}

// FacetItem 分面（按分类）统计条目。
type FacetItem struct {
	CategoryID string `json:"category_id"`
	Name       string `json:"name"`
	Count      int    `json:"count"`
}

// SearchFaceted 搜索并返回按分类分组的计数。
func (s *Service) SearchFaceted(indexID, query string) ([]FacetItem, int, error) {
	results, _, err := s.Search(indexID, query, 0)
	if err != nil {
		return nil, 0, err
	}
	categoryNames := make(map[string]string)
	for _, c := range s.store.ListCategories() {
		categoryNames[c.ID] = c.Name
	}
	counts := make(map[string]int)
	for _, r := range results {
		counts[r.CategoryID]++
	}
	items := make([]FacetItem, 0, len(counts))
	for catID, count := range counts {
		name := categoryNames[catID]
		if name == "" {
			name = "未分类"
		}
		items = append(items, FacetItem{CategoryID: catID, Name: name, Count: count})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Count > items[j].Count })
	return items, len(results), nil
}

// RebuildIndex 重建索引：清空该索引的倒排记录，重新索引全部活动文档。
func (s *Service) RebuildIndex(indexID string) (int, error) {
	idx, err := s.store.GetIndex(indexID)
	if err != nil {
		return 0, err
	}
	for _, p := range s.store.ListPostings() {
		if p.IndexID == indexID {
			_ = s.store.DeletePosting(p.ID)
		}
	}
	for _, t := range s.store.ListTerms() {
		if t.IndexID == indexID {
			_ = s.store.DeleteTerm(t.ID)
		}
	}
	idx.DocCount = 0
	idx.UpdatedAt = time.Now()
	if err := s.store.UpdateIndex(idx); err != nil {
		return 0, err
	}
	count := 0
	for _, d := range s.store.ListDocuments() {
		if d.Status != model.DocumentActive {
			continue
		}
		if err := s.IndexDocument(indexID, d.ID); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}

// IndexDocCount 返回索引中已索引的活动文档数。
func (s *Service) IndexDocCount(indexID string) (int, error) {
	idx, err := s.store.GetIndex(indexID)
	if err != nil {
		return 0, err
	}
	return idx.DocCount, nil
}
