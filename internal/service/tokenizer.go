package service

import (
	"strings"
	"unicode"
)

// Tokenize 将文本分词：统一小写、按非字母数字切分、剔除停用词、中文按单字切分。
func Tokenize(text string, stopwords map[string]bool) []string {
	text = strings.ToLower(text)
	tokens := make([]string, 0)
	var current []rune

	flush := func() {
		if len(current) > 0 {
			w := string(current)
			if !stopwords[w] {
				tokens = append(tokens, w)
			}
			current = current[:0]
		}
	}

	for _, r := range text {
		if unicode.Is(unicode.Han, r) {
			// 中文逐字切分
			flush()
			w := string(r)
			if !stopwords[w] {
				tokens = append(tokens, w)
			}
			continue
		}
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current = append(current, r)
			continue
		}
		flush()
	}
	flush()
	return tokens
}
