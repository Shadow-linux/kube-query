package query

import (
	"github.com/c-bata/go-prompt"
	"k8s.io/utils/strings/slices"
	"strings"
)

func RuleHasPrefix(word string, prefix ...string) bool {
	for _, p := range prefix {
		if strings.HasPrefix(word, p) {
			return true
		}
	}
	return false
}

func RuleJudgeLineWordCount(line string, cnt int, cmp CompareAct) bool {
	lineCnt := len(FormatLineWithSpace(line))
	Logger("Line cnt: ", lineCnt)
	switch cmp {
	case Equal:
		return lineCnt == cnt
	case Greater:
		return lineCnt > cnt
	case Less:
		return lineCnt < cnt
	case GreaterEqual:
		return lineCnt >= cnt
	case LessEqual:
		return lineCnt <= cnt
	default:
		return false
	}
}

func RuleJudgeWordExists(word string, words ...string) bool {
	for _, w := range words {
		if word == w {
			return true
		}
	}
	return false
}

func RuleJudgeLineHasWords(line string, words ...string) bool {
	var res bool
	lines := FormatLineWithSpace(line)
	for _, w := range words {
		if !slices.Contains(lines, w) {
			res = false
			break
		}
		res = true
	}
	return res
}

func RuleJudgeLineHasSpace(line string) bool {
	return strings.Contains(line, Space)
}

func RuleJudgeLastWordIsOption(line string) bool {
	lines := strings.Split(line, " ")
	// last one
	word := lines[len(lines)-1]
	if strings.HasPrefix(word, "-") {
		return true
	}
	return false
}

func RuleCanRemind(d prompt.Document) bool {
	if strings.Contains(d.TextBeforeCursor(), "|") {
		return false
	}
	return true
}

func RuleCanRemindHelper(d prompt.Document) bool {
	if strings.Contains(d.TextBeforeCursor(), "|") {
		return false
	}
	lines := FormatLineWithSpace(d.TextBeforeCursor())
	word := lines[len(lines)-2]
	if !strings.HasPrefix(word, "-") {
		return true
	}
	return false
}

func RuleJudgeLabelSelectorMatch(selector, label map[string]string) bool {
	for k, v := range selector {
		lv, ok := label[k]
		if ok {
			if lv == v {
				return true
			}
		}
	}
	return false
}
