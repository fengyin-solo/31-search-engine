package service

import (
	"testing"

	"searchengine/internal/config"
	"searchengine/internal/model"
	"searchengine/internal/store"
	"searchengine/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func TestTokenize(t *testing.T) {
	stop := map[string]bool{"的": true, "the": true}
	tokens := Tokenize("Go 语言 和 Channel 的并发", stop)
	if len(tokens) == 0 {
		t.Fatal("expect tokens")
	}
	// 应包含 go、语言、channel、并发（"的"被过滤）
	joined := ""
	for _, tk := range tokens {
		joined += tk + " "
	}
	for _, want := range []string{"go", "语", "言", "channel", "并", "发"} {
		if !contains(tokens, want) {
			t.Fatalf("expect token %q in %v", want, tokens)
		}
	}
	_ = joined
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func TestIndexAndSearch(t *testing.T) {
	s := newTestService()
	idx, _ := s.CreateIndex(model.Index{Name: "docs"})
	if _, err := s.ActivateIndex(idx.ID); err != nil {
		t.Fatalf("activate: %v", err)
	}
	doc, _ := s.CreateDocument(model.Document{Title: "Go 并发编程", Body: "goroutine 和 channel 是并发核心"})
	if err := s.IndexDocument(idx.ID, doc.ID); err != nil {
		t.Fatalf("index: %v", err)
	}
	results, total, err := s.Search(idx.ID, "并发", 10)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if total < 1 || len(results) < 1 {
		t.Fatalf("expect results, got %d", total)
	}
	if results[0].DocID != doc.ID {
		t.Fatalf("expect doc %s, got %s", doc.ID, results[0].DocID)
	}
}

func TestSearchSynonymExpansion(t *testing.T) {
	s := newTestService()
	if _, err := s.CreateSynonym(model.Synonym{Word: "汽车", Synonyms: []string{"轿车"}}); err != nil {
		t.Fatalf("synonym: %v", err)
	}
	idx, _ := s.CreateIndex(model.Index{Name: "docs"})
	_, _ = s.ActivateIndex(idx.ID)
	doc, _ := s.CreateDocument(model.Document{Title: "轿车保养", Body: "定期更换机油"})
	if err := s.IndexDocument(idx.ID, doc.ID); err != nil {
		t.Fatalf("index: %v", err)
	}
	// 用"汽车"查询应能命中含"轿车"的文档
	results, total, err := s.Search(idx.ID, "汽车", 10)
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if total < 1 {
		t.Fatal("expect synonym expansion hit")
	}
	_ = results
}

func TestRemoveDocumentFromIndex(t *testing.T) {
	s := newTestService()
	idx, _ := s.CreateIndex(model.Index{Name: "docs"})
	_, _ = s.ActivateIndex(idx.ID)
	doc, _ := s.CreateDocument(model.Document{Title: "待删除", Body: "内容"})
	if err := s.IndexDocument(idx.ID, doc.ID); err != nil {
		t.Fatalf("index: %v", err)
	}
	got, _ := s.GetIndex(idx.ID)
	if got.DocCount != 1 {
		t.Fatalf("expect doc count 1, got %d", got.DocCount)
	}
	if err := s.DeleteDocument(doc.ID); err != nil {
		t.Fatalf("delete document: %v", err)
	}
	got, _ = s.GetIndex(idx.ID)
	if got.DocCount != 0 {
		t.Fatalf("expect doc count 0 after delete, got %d", got.DocCount)
	}
	if len(s.store.ListPostings()) != 0 {
		t.Fatalf("expect 0 postings after delete")
	}
}

func TestIndexStateMachine(t *testing.T) {
	s := newTestService()
	idx, _ := s.CreateIndex(model.Index{Name: "docs"})
	if idx.Status != model.IndexCreated {
		t.Fatalf("expect created, got %s", idx.Status)
	}
	idx, _ = s.ActivateIndex(idx.ID)
	if idx.Status != model.IndexReady {
		t.Fatalf("expect ready, got %s", idx.Status)
	}
	idx, _ = s.DeactivateIndex(idx.ID)
	if idx.Status != model.IndexDeleted {
		t.Fatalf("expect deleted, got %s", idx.Status)
	}
	if _, err := s.ActivateIndex(idx.ID); err == nil {
		t.Fatal("expect conflict reactivating deleted index")
	}
}
