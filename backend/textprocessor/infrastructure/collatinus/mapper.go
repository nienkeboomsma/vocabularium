package collatinus

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/nienkeboomsma/collatinus/domain"
)

func mapToWords(input io.Reader) *[]domain.Word {
	var words []domain.Word

	scanner := bufio.NewScanner(input)

	wordCount := 0
	sentenceCount := 1

	previousWordIndexInSentence := 0

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		cols := strings.Split(line, "\t")

		wordCount++

		if len(cols) != 10 && !slices.Contains(cols, "unknown") {
			fmt.Printf("⚠️ SKIPPED %d: malformed line: %#v\n", wordCount, cols)
			continue
		}

		wordIndexInSentence, err := strconv.Atoi(cols[2])
		if err != nil {
			fmt.Printf("⚠️ SKIPPED %d: failed to convert word index in sentence string to integer: %#v\n", wordCount, cols[2])
			continue
		}

		if wordIndexInSentence <= previousWordIndexInSentence {
			sentenceCount++
		}

		previousWordIndexInSentence = wordIndexInSentence

		word := domain.Word{
			WordIndex:           wordCount,
			SentenceIndex:       sentenceCount,
			WordIndexInSentence: wordIndexInSentence,
			OriginalForm:        strings.TrimSpace(cols[3]),
			Translation:         "unknown",
		}

		if len(cols) != 10 {
			words = append(words, word)
			continue
		}

		word.LemmaRaw = strings.TrimSpace(cols[5])
		word.LemmaRich = strings.TrimSpace(cols[6])
		word.Translation = strings.TrimSpace(cols[8])
		word.Tag = strings.TrimSpace(cols[4])
		word.MorphoSyntacticalAnalysis = strings.TrimSpace(cols[9])

		if cols[7] != "" {
			frequencyInLASLA, err := strconv.Atoi(cols[7])
			if err != nil {
				fmt.Printf("⚠️ SKIPPED %d: failed to convert frequency in LASLA string to integer: %#v\n", wordCount, cols[7])
				continue
			}

			word.FrequencyInLASLA = frequencyInLASLA
		}

		words = append(words, word)
	}

	return &words
}
