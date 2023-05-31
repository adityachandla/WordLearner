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

func getWordMeanings(lines []string) []WordMeaning {
	finalWordMeanings := make([]WordMeaning, 0, 1000)
	left := WordMeaning{}
	right := WordMeaning{}
	numberRegex := regexp.MustCompile("\\d+\\.")
	for i, line := range lines {
		matches := numberRegex.FindAllString(line, -1)
		phrases := getPhrases(line)
		if len(matches) == 0 {
			//No new words put stuff in existing words
			if len(phrases) == 1 {
				if isRightColumn(line) {
					right.meaning += " " + phrases[0]
				} else {
					left.meaning += " " + phrases[0]
				}
			} else if len(phrases) == 2 {
				left.meaning += " " + phrases[0]
				right.meaning += " " + phrases[1]
			}
		} else if len(matches) == 1 {
			//One of the words has ended 3 or 4
			matchedNum, _ := strconv.Atoi(strings.TrimRight(matches[0], "."))
			if matchedNum == left.number+1 {
				finalWordMeanings = append(finalWordMeanings, left)
				left = WordMeaning{
					number:  matchedNum,
					word:    phrases[1],
					meaning: phrases[2],
				}
				if len(phrases) == 4 {
					right.meaning += " " + phrases[3]
				}
			} else if matchedNum == right.number+1 {
				finalWordMeanings = append(finalWordMeanings, right)
				right = WordMeaning{number: matchedNum}
				if len(phrases) == 4 {
					right.word = phrases[2]
					right.meaning = phrases[3]
					left.meaning += " " + phrases[0]
				} else if len(phrases) == 3 {
					right.word = phrases[1]
					right.meaning = phrases[2]
				}
			}
		} else if len(matches) == 2 {
			//Both words have ended
			if i != 0 {
				finalWordMeanings = append(finalWordMeanings, left)
				finalWordMeanings = append(finalWordMeanings, right)
				left = WordMeaning{}
				right = WordMeaning{}
			}
			left.number, _ = strconv.Atoi(strings.TrimRight(matches[0], "."))
			right.number, _ = strconv.Atoi(strings.TrimRight(matches[1], "."))
			if len(phrases) != 6 {
				panic("Invalid sequence")
			}
			left.word = phrases[1]
			left.meaning = phrases[2]

			right.word = phrases[4]
			right.meaning = phrases[5]
		}
	}
	return finalWordMeanings
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
