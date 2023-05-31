package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type WordMeaning struct {
	number  int
	word    string
	meaning string
}

func main() {
	data, err := os.ReadFile("./manhattan_prep_1000_gre_words_.txt")
	check(err)
	lines := strings.Split(string(data), "\n")
	wordMeanings := getWordMeanings(lines)
	sort.Slice(wordMeanings, func(i, j int) bool {
		return wordMeanings[i].number < wordMeanings[j].number
	})
	for _, wm := range wordMeanings[800:810] {
		fmt.Println(wm)
	}
}

type LineProcessor struct {
	left         WordMeaning
	right        WordMeaning
	wordMeanings []WordMeaning
}

type LineInfo struct {
	line    string
	lineIdx int
	matches []string
	phrases []string
}

func getWordMeanings(lines []string) []WordMeaning {
	processor := LineProcessor{}
	numberRegex := regexp.MustCompile("\\d+\\.")
	for i, line := range lines {
		lineInfo := LineInfo{
			lineIdx: i,
			line:    line,
			matches: numberRegex.FindAllString(line, -1),
			phrases: getPhrases(line),
		}
		if len(lineInfo.matches) == 0 {
			processor.handleZeroMatches(&lineInfo)
		} else if len(lineInfo.matches) == 1 {
			processor.handleSingleMatch(&lineInfo)
		} else {
			processor.handleTwoMatches(&lineInfo)
		}
	}
	return processor.wordMeanings
}

func (processor *LineProcessor) handleZeroMatches(lineInfo *LineInfo) {
	if len(lineInfo.phrases) == 1 {
		if lineInfo.isRightColumn() {
			processor.right.meaning += " " + lineInfo.phrases[0]
		} else {
			processor.left.meaning += " " + lineInfo.phrases[0]
		}
	} else if len(lineInfo.phrases) == 2 {
		processor.left.meaning += " " + lineInfo.phrases[0]
		processor.right.meaning += " " + lineInfo.phrases[1]
	}
}

func (processor *LineProcessor) handleSingleMatch(lineInfo *LineInfo) {
	matchedNum, _ := strconv.Atoi(strings.TrimRight(lineInfo.matches[0], "."))
	if matchedNum == processor.left.number+1 {
		processor.wordMeanings = append(processor.wordMeanings, processor.left)
		processor.left = WordMeaning{
			number:  matchedNum,
			word:    lineInfo.phrases[1],
			meaning: lineInfo.phrases[2],
		}
		if len(lineInfo.phrases) == 4 {
			processor.right.meaning += " " + lineInfo.phrases[3]
		}
	} else if matchedNum == processor.right.number+1 {
		processor.wordMeanings = append(processor.wordMeanings, processor.right)
		processor.right = WordMeaning{number: matchedNum}
		if len(lineInfo.phrases) == 4 {
			processor.right.word = lineInfo.phrases[2]
			processor.right.meaning = lineInfo.phrases[3]
			processor.left.meaning += " " + lineInfo.phrases[0]
		} else if len(lineInfo.phrases) == 3 {
			processor.right.word = lineInfo.phrases[1]
			processor.right.meaning = lineInfo.phrases[2]
		}
	}
}

func (processor *LineProcessor) handleTwoMatches(lineInfo *LineInfo) {
	if lineInfo.lineIdx != 0 {
		processor.wordMeanings = append(processor.wordMeanings, processor.right)
		processor.wordMeanings = append(processor.wordMeanings, processor.left)
	}
	processor.left.number, _ = strconv.Atoi(strings.TrimRight(lineInfo.matches[0], "."))
	processor.right.number, _ = strconv.Atoi(strings.TrimRight(lineInfo.matches[1], "."))
	processor.left.word = lineInfo.phrases[1]
	processor.left.meaning = lineInfo.phrases[2]

	processor.right.word = lineInfo.phrases[4]
	processor.right.meaning = lineInfo.phrases[5]
}

func (lineInfo *LineInfo) isRightColumn() bool {
	if len(lineInfo.line) < 50 {
		return false
	}
	for _, b := range lineInfo.line[:50] {
		if b != ' ' {
			return false
		}
	}
	return true
}

func getPhrases(line string) []string {
	line = strings.TrimLeft(line, " ")
	lineBytes := []byte(line)
	phrases := make([]string, 0, 4)
	var sb strings.Builder
	whitespaceCount := 0
	for _, b := range lineBytes {
		if b == ' ' {
			whitespaceCount++
		} else {
			whitespaceCount = 0
		}
		if whitespaceCount < 2 {
			sb.WriteByte(b)
		}
		if whitespaceCount == 2 {
			phrases = append(phrases, sb.String())
			sb.Reset()
		}
	}
	if sb.Len() > 0 {
		phrases = append(phrases, strings.Trim(sb.String(), "\n\r"))
	}
	return phrases
}
