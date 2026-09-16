package service

import (
	"sort"
	"strings"
	"time"

	"searchengine/internal/model"
	"searchengine/pkg/idgen"
)

// CreateSpellCorrection 创建拼写纠错词对。
func (s *Service) CreateSpellCorrection(input model.SpellCorrection) (*model.SpellCorrection, error) {
	input.ID = idgen.Hex()
	now := time.Now()
	input.CreatedAt = now
	input.UpdatedAt = now
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.CreateSpellCorrection(&input); err != nil {
		return nil, err
	}
	return &input, nil
}

// GetSpellCorrection 获取纠错词对。
func (s *Service) GetSpellCorrection(id string) (*model.SpellCorrection, error) {
	return s.store.GetSpellCorrection(id)
}

// ListSpellCorrections 分页列出纠错词对。
func (s *Service) ListSpellCorrections(page, size int) ([]*model.SpellCorrection, int, error) {
	all := s.store.ListSpellCorrections()
	sort.Slice(all, func(i, j int) bool { return all[i].Word < all[j].Word })
	total := len(all)
	start := (page - 1) * size
	if start >= total {
		return []*model.SpellCorrection{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return all[start:end], total, nil
}

// UpdateSpellCorrection 更新纠错词对。
func (s *Service) UpdateSpellCorrection(id string, input model.SpellCorrection) (*model.SpellCorrection, error) {
	existing, err := s.store.GetSpellCorrection(id)
	if err != nil {
		return nil, err
	}
	existing.Word = input.Word
	existing.Correct = input.Correct
	existing.Status = input.Status
	existing.UpdatedAt = time.Now()
	if err := existing.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateSpellCorrection(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// DeleteSpellCorrection 删除纠错词对。
func (s *Service) DeleteSpellCorrection(id string) error {
	return s.store.DeleteSpellCorrection(id)
}

// Correct 对查询词做拼写纠错：先查词典精确匹配，再用编辑距离匹配最近词条。
func (s *Service) Correct(word string) (string, bool) {
	word = strings.ToLower(strings.TrimSpace(word))
	if word == "" {
		return "", false
	}
	// 精确匹配纠错词对
	if sc, err := s.store.GetSpellCorrectionByWord(word); err == nil && sc.Status == model.SpellCorrectionActive {
		return sc.Correct, true
	}
	// 编辑距离匹配词典词条
	best := ""
	bestDist := 3 // 最大容忍编辑距离
	for _, t := range s.store.ListTerms() {
		d := editDistance(word, t.Term)
		if d < bestDist {
			bestDist = d
			best = t.Term
		}
	}
	if best != "" && best != word {
		return best, true
	}
	return "", false
}

// editDistance 计算两个字符串的莱文斯坦编辑距离。
func editDistance(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}
	dp := make([][]int, la+1)
	for i := range dp {
		dp[i] = make([]int, lb+1)
		dp[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		dp[0][j] = j
	}
	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			dp[i][j] = min3(dp[i-1][j]+1, dp[i][j-1]+1, dp[i-1][j-1]+cost)
		}
	}
	return dp[la][lb]
}

func min3(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
