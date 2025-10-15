package pkg

import (
	"regexp"
	"strings"
)

// MatchEndOfSentence 检查文本是否以句子结尾
// 实现了Python版本的句子结尾检测逻辑，但适配了Go的RE2语法限制
func MatchEndOfSentence(text string) bool {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return false
	}

	// 检查是否以句子结束标点符号结尾
	endPattern := regexp.MustCompile(`[\,，\.。\?？\!！:：]$`)
	if !endPattern.MatchString(trimmed) {
		return false
	}

	// 获取最后一个单词来检查特殊情况
	words := strings.Fields(trimmed)
	if len(words) == 0 {
		return false
	}

	lastWord := words[len(words)-1]

	// 检查是否是时间格式（数字后跟冒号和空格以及am/pm）
	timePattern1 := regexp.MustCompile(`\d+:\d{2}\s*[ap]\.?m\.?\.?$`)
	timePattern2 := regexp.MustCompile(`\d+\s*[ap]\.?m\.?\.?$`)
	if timePattern1.MatchString(trimmed) || timePattern1.MatchString(lastWord) ||
		timePattern2.MatchString(trimmed) || timePattern2.MatchString(lastWord) {
		return false
	}

	// 检查常见的缩写
	abbreviations := []string{"Mr.", "Ms.", "Dr.", "Mrs.", "Prof."}
	for _, abbr := range abbreviations {
		if strings.HasSuffix(trimmed, abbr) || strings.HasSuffix(lastWord, abbr) {
			return false
		}
	}

	// 检查全大写字母缩写（如U.S.A.）
	upperAbbrPattern := regexp.MustCompile(`^[A-Z]\.([A-Z]\.)*$`)
	if upperAbbrPattern.MatchString(lastWord) {
		return false
	}

	// 检查数字后跟点（列表项）
	numberedListPattern := regexp.MustCompile(`^\d+\.$`)
	if numberedListPattern.MatchString(lastWord) {
		return false
	}

	// 检查是否是大写字母后跟句号（可能是缩写的一部分）
	upperLetterPattern := regexp.MustCompile(`^[A-Z]\.$`)
	if upperLetterPattern.MatchString(lastWord) && len(words) > 1 {
		// 检查前一个词是否也是类似格式
		prevWord := words[len(words)-2]
		if upperLetterPattern.MatchString(prevWord) {
			return false
		}
	}

	return true
}
