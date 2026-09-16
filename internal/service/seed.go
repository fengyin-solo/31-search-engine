package service

import (
	"fmt"

	"searchengine/internal/model"
)

// SeedResult 初始化演示数据的结果统计。
type SeedResult struct {
	Documents   int `json:"documents"`
	Indexes     int `json:"indexes"`
	Analyzers   int `json:"analyzers"`
	StopWords   int `json:"stop_words"`
	Synonyms    int `json:"synonyms"`
	Terms       int `json:"terms"`
	Postings    int `json:"postings"`
}

// SeedDemoData 初始化一套演示数据，并完成索引构建。
func (s *Service) SeedDemoData() (*SeedResult, error) {
	res := &SeedResult{}

	// 停用词
	for _, w := range []string{"的", "了", "和", "是", "在", "the", "a", "an", "of", "to", "and"} {
		if _, err := s.CreateStopWord(model.StopWord{Word: w}); err != nil {
			return nil, fmt.Errorf("seed stopword: %w", err)
		}
		res.StopWords++
	}

	// 分词器
	analyzer, err := s.CreateAnalyzer(model.Analyzer{
		Name: "standard", Type: model.AnalyzerTypeStandard,
		StopWords: []string{"的", "了", "和", "是", "在", "the", "a", "an", "of", "to", "and"},
	})
	if err != nil {
		return nil, fmt.Errorf("seed analyzer: %w", err)
	}
	res.Analyzers++

	// 同义词
	synonyms := []model.Synonym{
		{Word: "汽车", Synonyms: []string{"轿车", "车辆"}},
		{Word: "电脑", Synonyms: []string{"计算机", "笔记本"}},
		{Word: "手机", Synonyms: []string{"电话", "移动设备"}},
	}
	for _, sy := range synonyms {
		if _, err := s.CreateSynonym(sy); err != nil {
			return nil, err
		}
		res.Synonyms++
	}

	// 分类
	categories := []model.Category{
		{Name: "技术", Description: "编程与系统设计"},
		{Name: "生活", Description: "生活百科与健康"},
		{Name: "购物", Description: "选购指南"},
	}
	catByName := make(map[string]string)
	for _, c := range categories {
		created, err := s.CreateCategory(c)
		if err != nil {
			return nil, err
		}
		catByName[c.Name] = created.ID
	}

	// 索引
	idx, err := s.CreateIndex(model.Index{
		Name: "docs", Description: "全站文档索引",
		Fields: []string{"title", "body"}, AnalyzerID: analyzer.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("seed index: %w", err)
	}
	if _, err := s.ActivateIndex(idx.ID); err != nil {
		return nil, err
	}
	res.Indexes++

	// 文档
	docs := []model.Document{
		{Title: "Go 语言并发编程指南", Body: "goroutine 和 channel 是 Go 并发编程的核心，本指南介绍并发模型与最佳实践", Source: "技术博客"},
		{Title: "Python 数据分析入门", Body: "使用 pandas 进行数据分析，处理表格数据与可视化", Source: "技术博客"},
		{Title: "汽车保养小常识", Body: "定期更换机油与滤芯，保持轿车发动机良好状态", Source: "生活百科"},
		{Title: "笔记本电脑选购指南", Body: "如何根据需求选择计算机，介绍处理器与内存的搭配", Source: "购物指南"},
		{Title: "智能手机拍照技巧", Body: "掌握构图与光线，用手机也能拍出好照片", Source: "生活百科"},
		{Title: "分布式系统设计模式", Body: "介绍微服务、消息队列与缓存的一致性设计模式", Source: "技术博客"},
		{Title: "健康饮食搭配", Body: "均衡膳食，多吃蔬菜水果，减少高糖高油食物摄入", Source: "健康"},
		{Title: "机器学习基础", Body: "从线性回归到神经网络，掌握机器学习核心算法", Source: "技术博客"},
	}
	for _, d := range docs {
		created, err := s.CreateDocument(d)
		if err != nil {
			return nil, fmt.Errorf("seed document: %w", err)
		}
		res.Documents++
		// 建立索引
		if err := s.IndexDocument(idx.ID, created.ID); err != nil {
			return nil, fmt.Errorf("index document: %w", err)
		}
		// 按来源归类
		if catID, ok := catByName[sourceToCategory(created.Source)]; ok {
			if _, err := s.SetDocumentCategory(created.ID, catID); err != nil {
				return nil, err
			}
		}
	}

	// 加权规则
	boostRules := []model.BoostRule{
		{Field: model.BoostFieldTitle, Term: "并发", Boost: 2.0},
		{Field: model.BoostFieldTitle, Term: "汽车", Boost: 1.5},
	}
	for _, b := range boostRules {
		if _, err := s.CreateBoostRule(b); err != nil {
			return nil, err
		}
	}

	// 搜索建议词
	suggestions := []model.Suggestion{
		{Term: "并发", Weight: 100},
		{Term: "数据分析", Weight: 80},
		{Term: "汽车保养", Weight: 60},
		{Term: "机器学习", Weight: 90},
	}
	for _, sug := range suggestions {
		if _, err := s.CreateSuggestion(sug); err != nil {
			return nil, err
		}
	}

	// 拼写纠错词对
	spellCorrections := []model.SpellCorrection{
		{Word: "golang", Correct: "go"},
		{Word: "pyhton", Correct: "python"},
	}
	for _, sc := range spellCorrections {
		if _, err := s.CreateSpellCorrection(sc); err != nil {
			return nil, err
		}
	}

	res.Terms = len(s.store.ListTerms())
	res.Postings = len(s.store.ListPostings())
	return res, nil
}

// sourceToCategory 将文档来源映射到分类名。
func sourceToCategory(source string) string {
	switch source {
	case "技术博客":
		return "技术"
	case "生活百科", "健康":
		return "生活"
	case "购物指南":
		return "购物"
	default:
		return "生活"
	}
}
