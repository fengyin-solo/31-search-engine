package store

import (
	"testing"

	"searchengine/internal/model"
)

func TestDocumentCRUD(t *testing.T) {
	s := NewMemoryStore()
	d := &model.Document{ID: "d1", Title: "标题", Body: "正文", Status: model.DocumentActive}
	if err := s.CreateDocument(d); err != nil {
		t.Fatalf("create: %v", err)
	}
	got, err := s.GetDocument("d1")
	if err != nil || got.Title != "标题" {
		t.Fatalf("get failed: %v", err)
	}
	got.Status = model.DocumentDeleted
	if err := s.UpdateDocument(got); err != nil {
		t.Fatalf("update: %v", err)
	}
	if len(s.ListDocuments()) != 1 {
		t.Fatalf("expect 1 document")
	}
	if err := s.DeleteDocument("d1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestIndexUniqueness(t *testing.T) {
	s := NewMemoryStore()
	i := &model.Index{ID: "i1", Name: "docs", Status: model.IndexCreated}
	if err := s.CreateIndex(i); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.CreateIndex(&model.Index{ID: "i2", Name: "docs"}); err != ErrConflict {
		t.Fatalf("expect conflict, got %v", err)
	}
	if _, err := s.GetIndexByName("docs"); err != nil {
		t.Fatalf("get by name: %v", err)
	}
	if err := s.DeleteIndex("i1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestAnalyzerCRUD(t *testing.T) {
	s := NewMemoryStore()
	a := &model.Analyzer{ID: "a1", Name: "standard", Type: model.AnalyzerTypeStandard}
	if err := s.CreateAnalyzer(a); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := s.GetAnalyzerByName("standard"); err != nil {
		t.Fatalf("get by name: %v", err)
	}
	if err := s.DeleteAnalyzer("a1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestStopWordCRUD(t *testing.T) {
	s := NewMemoryStore()
	w := &model.StopWord{ID: "w1", Word: "the"}
	if err := s.CreateStopWord(w); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.CreateStopWord(&model.StopWord{ID: "w2", Word: "the"}); err != ErrConflict {
		t.Fatalf("expect conflict, got %v", err)
	}
	if _, err := s.GetStopWordByWord("the"); err != nil {
		t.Fatalf("get by word: %v", err)
	}
	if err := s.DeleteStopWord("w1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestTermCRUD(t *testing.T) {
	s := NewMemoryStore()
	tm := &model.Term{ID: "t1", IndexID: "i1", Term: "go", DocCount: 1, TotalTF: 3}
	if err := s.CreateTerm(tm); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := s.GetTermByKey("i1", "go"); err != nil {
		t.Fatalf("get by key: %v", err)
	}
	tm.DocCount = 2
	if err := s.UpdateTerm(tm); err != nil {
		t.Fatalf("update: %v", err)
	}
	if len(s.ListTerms()) != 1 {
		t.Fatalf("expect 1 term")
	}
	if err := s.DeleteTerm("t1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestPostingCRUD(t *testing.T) {
	s := NewMemoryStore()
	p := &model.Posting{ID: "p1", IndexID: "i1", Term: "go", DocID: "d1", TF: 2}
	if err := s.CreatePosting(p); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := s.GetPostingByKey("i1", "go", "d1"); err != nil {
		t.Fatalf("get by key: %v", err)
	}
	if len(s.ListPostings()) != 1 {
		t.Fatalf("expect 1 posting")
	}
	if err := s.DeletePosting("p1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestSynonymCRUD(t *testing.T) {
	s := NewMemoryStore()
	sy := &model.Synonym{ID: "s1", Word: "汽车", Synonyms: []string{"轿车"}}
	if err := s.CreateSynonym(sy); err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := s.CreateSynonym(&model.Synonym{ID: "s2", Word: "汽车"}); err != ErrConflict {
		t.Fatalf("expect conflict, got %v", err)
	}
	if _, err := s.GetSynonymByWord("汽车"); err != nil {
		t.Fatalf("get by word: %v", err)
	}
	if err := s.DeleteSynonym("s1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestQueryLogCRUD(t *testing.T) {
	s := NewMemoryStore()
	q := &model.QueryLog{ID: "q1", Query: "go", ResultCount: 3}
	if err := s.CreateQueryLog(q); err != nil {
		t.Fatalf("create: %v", err)
	}
	if len(s.ListQueryLogs()) != 1 {
		t.Fatalf("expect 1 query log")
	}
	if err := s.DeleteQueryLog("q1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestHotSearchCRUD(t *testing.T) {
	s := NewMemoryStore()
	h := &model.HotSearch{ID: "h1", Term: "go", Count: 5}
	if err := s.CreateHotSearch(h); err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := s.GetHotSearchByTerm("go"); err != nil {
		t.Fatalf("get by term: %v", err)
	}
	h.Count = 6
	if err := s.UpdateHotSearch(h); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := s.DeleteHotSearch("h1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}
