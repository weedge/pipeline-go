package pkg

import (
	"testing"
)

func TestMatchEndOfSentence(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{
			name:     "Simple sentence ending with period",
			text:     "This is a test.",
			expected: true,
		},
		{
			name:     "Sentence ending with exclamation mark",
			text:     "Wow!",
			expected: true,
		},
		{
			name:     "Sentence ending with question mark",
			text:     "How are you?",
			expected: true,
		},
		{
			name:     "Sentence ending with Chinese period",
			text:     "这是一个测试。",
			expected: true,
		},
		{
			name:     "Sentence ending with colon",
			text:     "Important note:",
			expected: true,
		},
		{
			name:     "Text not ending with sentence punctuation",
			text:     "This is not the end",
			expected: false,
		},
		{
			name:     "Complete sentence with abbreviation",
			text:     "I saw Mr. Smith yesterday.",
			expected: true, // 句子以"yesterday."结尾，是一个完整的句子
		},
		{
			name:     "Time format",
			text:     "The meeting is at 3:00 p.m.",
			expected: false, // 句子以"p.m."结尾，不是完整句子结尾
		},
		{
			name:     "Number followed by period",
			text:     "1. First item",
			expected: false, // 以"1."结尾，不是完整句子结尾
		},
		{
			name:     "Uppercase abbreviation",
			text:     "The U.S.A. is a country",
			expected: false, // 这个句子没有以句子结束标点结尾
		},
		{
			name:     "Uppercase abbreviation with ending",
			text:     "I visited the U.S.A.. It was fun.",
			expected: true, // 句子以"It was fun."结尾
		},
		{
			name:     "Text with trailing spaces",
			text:     "This is a test.   ",
			expected: true,
		},
		{
			name:     "Sentence with abbreviation in middle",
			text:     "I saw Mr. Smith yesterday. He was nice.",
			expected: true, // 句子以"nice."结尾
		},
		{
			name:     "Just an abbreviation",
			text:     "Mr.",
			expected: false, // 单独的缩写不构成完整句子结尾
		},
		{
			name:     "Empty string",
			text:     "",
			expected: false, // 空字符串
		},
		{
			name:     "Only spaces",
			text:     "   ",
			expected: false, // 只有空格
		},
		{
			name:     "Multiple sentences - ends with full sentence",
			text:     "First sentence. Second sentence!",
			expected: true, // 以感叹号结尾
		},
		{
			name:     "Doctor abbreviation",
			text:     "Dr. Jones said hello.",
			expected: true, // 句子以"hello."结尾
		},
		{
			name:     "Professor abbreviation",
			text:     "Prof. Smith teaches CS.",
			expected: true, // 句子以"CS."结尾
		},
		{
			name:     "List item with period",
			text:     "Item 2. Description",
			expected: false, // 以"2."结尾
		},
		{
			name:     "Time without minutes",
			text:     "Meet me at 3 p.m.",
			expected: false, // 以"p.m."结尾
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MatchEndOfSentence(tt.text)
			if result != tt.expected {
				t.Errorf("MatchEndOfSentence(%q) = %v, want %v", tt.text, result, tt.expected)
			}
		})
	}
}
