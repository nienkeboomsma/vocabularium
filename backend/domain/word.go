package domain

type Word struct {
	WordIndex                 int    `json:"wordIndex,omitzero"`
	SentenceIndex             int    `json:"sentenceIndex,omitzero"`
	WordIndexInSentence       int    `json:"wordIndexInSentence,omitzero"`
	OriginalForm              string `json:"originalForm,omitzero"`
	LemmaRaw                  string `json:"lemmaRaw,omitzero"`
	LemmaRich                 string `json:"lemmaRich,omitzero"`
	Translation               string `json:"translations,omitzero"`
	FrequencyInLASLA          int    `json:"frequencyInLASLA,omitzero"`
	Tag                       string `json:"tag,omitzero"`
	MorphoSyntacticalAnalysis string `json:"morphoSyntacticalAnalisys,omitzero"`
}
