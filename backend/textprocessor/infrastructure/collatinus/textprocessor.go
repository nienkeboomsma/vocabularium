package collatinus

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/nienkeboomsma/collatinus/domain"
)

type Language string

const (
	LanguageCA Language = "-tca"
	LanguageDE Language = "-tde"
	LanguageEN Language = "-ten"
	LanguageES Language = "-tes"
	LanguageEU Language = "-teu"
	LanguageFR Language = "-tfr"
	LanguageGL Language = "-tgl"
	LanguageIT Language = "-tit"
	LanguageNL Language = "-tnl"
	LanguagePT Language = "-tpt"
)

type TextProcessor struct {
	chunkSize int
	language  Language
}

func NewTextProcessor(language string) (*TextProcessor, error) {
	var validatedLanguage Language

	switch strings.ToLower(language) {
	case "ca":
		validatedLanguage = LanguageCA
	case "de":
		validatedLanguage = LanguageDE
	case "en":
		validatedLanguage = LanguageEN
	case "es":
		validatedLanguage = LanguageES
	case "eu":
		validatedLanguage = LanguageEU
	case "fr":
		validatedLanguage = LanguageFR
	case "gl":
		validatedLanguage = LanguageGL
	case "it":
		validatedLanguage = LanguageIT
	case "nl":
		validatedLanguage = LanguageNL
	case "pt":
		validatedLanguage = LanguagePT
	default:
		return &TextProcessor{}, errors.New(`COLLATINUS_LANGUAGE must be one of "ca", "de", "en", "es", "eu", "fr", "gl", "it", "nl" or "pt"`)
	}

	return &TextProcessor{
		language: validatedLanguage,
	}, nil
}

func (tp *TextProcessor) Process(input []byte) (*[]domain.WorkWord, *map[uuid.UUID]domain.Word, []string, error) {
	sanitised := sanitise(input)

	chunks := chunkBySentence(sanitised)

	var total bytes.Buffer

	err := lemmatise(chunks, &total, LanguageEN)
	if err != nil {
		return &[]domain.WorkWord{}, &map[uuid.UUID]domain.Word{}, []string{}, fmt.Errorf("failed to lemmatise: %w", err)
	}

	workWords, words, logs := mapToWords(&total)

	return workWords, words, logs, nil
}
