package collatinus

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/nienkeboomsma/vocabularium/database"
	"github.com/nienkeboomsma/vocabularium/domain"
)

func mapToWords(input io.Reader) (*[]domain.WorkWord, *map[uuid.UUID]domain.Word, []string) {
	workWords := []domain.WorkWord{}
	words := make(map[uuid.UUID]domain.Word)
	logs := []string{}

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

		if slices.Contains(cols, "unknown") {
			logs = append(logs, fmt.Sprintf("⚠️ SKIPPED word #%d: translation unknown: %#v", wordCount, cols))
			continue
		}

		if len(cols) != 10 {
			logs = append(logs, fmt.Sprintf("⚠️ SKIPPED word #%d: malformed line: %#v", wordCount, cols))
			continue
		}

		wordIndexInSentence, err := strconv.Atoi(cols[2])
		if err != nil {
			logs = append(logs, fmt.Sprintf("⚠️ SKIPPED word #%d: failed to convert word index in sentence string to integer: %#v", wordCount, cols[2]))
			continue
		}

		if wordIndexInSentence <= previousWordIndexInSentence {
			sentenceCount++
		}

		previousWordIndexInSentence = wordIndexInSentence

		workWord := domain.WorkWord{
			WordIndex:     wordCount,
			SentenceIndex: sentenceCount,
			OriginalForm:  strings.TrimSpace(cols[3]),
		}

		if len(cols) != 10 {
			workWords = append(workWords, workWord)
			continue
		}

		word := domain.Word{
			LemmaRaw:    strings.TrimSpace(cols[5]),
			LemmaRich:   strings.TrimSpace(cols[6]),
			Translation: strings.TrimSpace(cols[8]),
		}

		word.ID = database.StringToUUID(fmt.Sprintf("%s_%s", word.LemmaRaw, word.LemmaRich))

		workWord.WordID = word.ID
		workWord.Tag = strings.TrimSpace(cols[4])
		workWord.MorphoSyntacticalAnalysis = strings.TrimSpace(cols[9])

		if cols[7] != "" {
			frequencyInLASLA, err := strconv.Atoi(cols[7])
			if err != nil {
				logs = append(logs, fmt.Sprintf("⚠️ SKIPPED word #%d: failed to convert frequency in LASLA string to integer: %#v", wordCount, cols[7]))
				continue
			}

			word.FrequencyInLASLA = frequencyInLASLA
		}

		workWords = append(workWords, workWord)
		words[word.ID] = word
	}

	return &workWords, &words, logs
}
