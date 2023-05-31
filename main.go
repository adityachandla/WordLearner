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
	lineIdx      int
	line         string
	matches      []string
	phrases      []string
	wordMeanings []WordMeaning
}

func getWordMeanings(lines []string) []WordMeaning {
	processor := LineProcessor{}
	numberRegex := regexp.MustCompile("\\d+\\.")
	for i, line := range lines {
		processor.lineIdx = i
		processor.line = line
		processor.matches = numberRegex.FindAllString(line, -1)
		processor.phrases = getPhrases(line)
		if len(processor.matches) == 0 {
			handleZeroMatches(&processor)
		} else if len(processor.matches) == 1 {
			handleSingleMatch(&processor)
		} else {
			handleTwoMatches(&processor)
		}
	}
	return processor.wordMeanings
}

func handleZeroMatches(processor *LineProcessor) {
	if len(processor.phrases) == 1 {
		if isRightColumn(processor.line) {
			processor.right.meaning += " " + processor.phrases[0]
		} else {
			processor.left.meaning += " " + processor.phrases[0]
		}
	} else if len(processor.phrases) == 2 {
		processor.left.meaning += " " + processor.phrases[0]
		processor.right.meaning += " " + processor.phrases[1]
	}
}

func handleSingleMatch(processor *LineProcessor) {
	matchedNum, _ := strconv.Atoi(strings.TrimRight(processor.matches[0], "."))
	if matchedNum == processor.left.number+1 {
		processor.wordMeanings = append(processor.wordMeanings, processor.left)
		processor.left = WordMeaning{
			number:  matchedNum,
			word:    processor.phrases[1],
			meaning: processor.phrases[2],
		}
		if len(processor.phrases) == 4 {
			processor.right.meaning += " " + processor.phrases[3]
		}
	} else if matchedNum == processor.right.number+1 {
		processor.wordMeanings = append(processor.wordMeanings, processor.right)
		processor.right = WordMeaning{number: matchedNum}
		if len(processor.phrases) == 4 {
			processor.right.word = processor.phrases[2]
			processor.right.meaning = processor.phrases[3]
			processor.left.meaning += " " + processor.phrases[0]
		} else if len(processor.phrases) == 3 {
			processor.right.word = processor.phrases[1]
			processor.right.meaning = processor.phrases[2]
		}
	}
}

func handleTwoMatches(processor *LineProcessor) {
	if processor.lineIdx != 0 {
		processor.wordMeanings = append(processor.wordMeanings, processor.right)
		processor.wordMeanings = append(processor.wordMeanings, processor.left)
	}
	processor.left.number, _ = strconv.Atoi(strings.TrimRight(processor.matches[0], "."))
	processor.right.number, _ = strconv.Atoi(strings.TrimRight(processor.matches[1], "."))
	processor.left.word = processor.phrases[1]
	processor.left.meaning = processor.phrases[2]

	processor.right.word = processor.phrases[4]
	processor.right.meaning = processor.phrases[5]
}

func isRightColumn(line string) bool {
	if len(line) < 50 {
		return false
	}
	for _, b := range line[:50] {
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
