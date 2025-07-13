package collatinus

import (
	"bytes"
	"io"
	"testing"

	"github.com/nienkeboomsma/collatinus/domain"
	"github.com/stretchr/testify/assert"
)

func TestMapper(t *testing.T) {
	tests := []struct {
		name     string
		input    io.Reader
		expected *[]domain.Word
	}{
		{
			name: "two sentences",
			input: bytes.NewReader([]byte(`
1	1	1	Pedicabo	v1 	pedico	pēdīco, as, are	4	to perform anal intercourse; to commit sodomy with;	pēdīcābō̆ future indicative active 1st singular
2	1	2	ego	p11	ego	ĕgō̆, mei, pron.	14846	moi, me	ĕgō̆ masculine nominative singular
3	1	3	vos	p32	vos	vōs, uestrum, pl. pron.	2402	you (pl.), ye;	vōs masculine accusative plural
1	1	1	et	d   (c  )	et	ĕt, conj. adv.	42726	and, and even; also, even; (et ... et = both ... and);	ĕt
2	1	2	irrumabo	v1 	inrumo	īnrŭmo, as, are	6	to extend the breast to, to give suck;	īrrŭmābō̆ future indicative active 1st singular
`)),
			expected: &[]domain.Word{
				{
					WordIndex:                 1,
					SentenceIndex:             1,
					WordIndexInSentence:       1,
					OriginalForm:              "Pedicabo",
					LemmaRaw:                  "pedico",
					LemmaRich:                 "pēdīco, as, are",
					Translation:               "to perform anal intercourse; to commit sodomy with;",
					FrequencyInLASLA:          4,
					Tag:                       "v1",
					MorphoSyntacticalAnalysis: "pēdīcābō̆ future indicative active 1st singular",
				},
				{
					WordIndex:                 2,
					SentenceIndex:             1,
					WordIndexInSentence:       2,
					OriginalForm:              "ego",
					LemmaRaw:                  "ego",
					LemmaRich:                 "ĕgō̆, mei, pron.",
					Translation:               "moi, me",
					FrequencyInLASLA:          14846,
					Tag:                       "p11",
					MorphoSyntacticalAnalysis: "ĕgō̆ masculine nominative singular",
				},
				{
					WordIndex:                 3,
					SentenceIndex:             1,
					WordIndexInSentence:       3,
					OriginalForm:              "vos",
					LemmaRaw:                  "vos",
					LemmaRich:                 "vōs, uestrum, pl. pron.",
					Translation:               "you (pl.), ye;",
					FrequencyInLASLA:          2402,
					Tag:                       "p32",
					MorphoSyntacticalAnalysis: "vōs masculine accusative plural",
				},
				{
					WordIndex:                 4,
					SentenceIndex:             2,
					WordIndexInSentence:       1,
					OriginalForm:              "et",
					LemmaRaw:                  "et",
					LemmaRich:                 "ĕt, conj. adv.",
					Translation:               "and, and even; also, even; (et ... et = both ... and);",
					FrequencyInLASLA:          42726,
					Tag:                       "d   (c  )",
					MorphoSyntacticalAnalysis: "ĕt",
				},
				{
					WordIndex:                 5,
					SentenceIndex:             2,
					WordIndexInSentence:       2,
					OriginalForm:              "irrumabo",
					LemmaRaw:                  "inrumo",
					LemmaRich:                 "īnrŭmo, as, are",
					Translation:               "to extend the breast to, to give suck;",
					FrequencyInLASLA:          6,
					Tag:                       "v1",
					MorphoSyntacticalAnalysis: "īrrŭmābō̆ future indicative active 1st singular",
				},
			},
		},
		{
			name: "malformed string",
			input: bytes.NewReader([]byte(`
1	1	1	Pedicabo	v1 	pedico	pēdīco, as, are	4	to perform anal intercourse; to commit sodomy with;	pēdīcābō̆ future indicative active 1st singular
2	1	2	ego	p11	ego	ĕgō̆, mei, pron.	14846	moi, me	ĕgō̆ masculine nominative singular
3	1	3	vos	p32	vos	vōs, uestrum, pl. pron.	2402	you (pl.), ye;	vōs masculine accusative plural
1	1	1	et	d   (c  )	et	ĕt, conj. adv.	42726	and, an
2	1	2	irrumabo	v1 	inrumo	īnrŭmo, as, are	6	to extend the breast to, to give suck;	īrrŭmābō̆ future indicative active 1st singular
`)),
			expected: &[]domain.Word{
				{
					WordIndex:                 1,
					SentenceIndex:             1,
					WordIndexInSentence:       1,
					OriginalForm:              "Pedicabo",
					LemmaRaw:                  "pedico",
					LemmaRich:                 "pēdīco, as, are",
					Translation:               "to perform anal intercourse; to commit sodomy with;",
					FrequencyInLASLA:          4,
					Tag:                       "v1",
					MorphoSyntacticalAnalysis: "pēdīcābō̆ future indicative active 1st singular",
				},
				{
					WordIndex:                 2,
					SentenceIndex:             1,
					WordIndexInSentence:       2,
					OriginalForm:              "ego",
					LemmaRaw:                  "ego",
					LemmaRich:                 "ĕgō̆, mei, pron.",
					Translation:               "moi, me",
					FrequencyInLASLA:          14846,
					Tag:                       "p11",
					MorphoSyntacticalAnalysis: "ĕgō̆ masculine nominative singular",
				},
				{
					WordIndex:                 3,
					SentenceIndex:             1,
					WordIndexInSentence:       3,
					OriginalForm:              "vos",
					LemmaRaw:                  "vos",
					LemmaRich:                 "vōs, uestrum, pl. pron.",
					Translation:               "you (pl.), ye;",
					FrequencyInLASLA:          2402,
					Tag:                       "p32",
					MorphoSyntacticalAnalysis: "vōs masculine accusative plural",
				},
				{
					WordIndex:                 5,
					SentenceIndex:             2,
					WordIndexInSentence:       2,
					OriginalForm:              "irrumabo",
					LemmaRaw:                  "inrumo",
					LemmaRich:                 "īnrŭmo, as, are",
					Translation:               "to extend the breast to, to give suck;",
					FrequencyInLASLA:          6,
					Tag:                       "v1",
					MorphoSyntacticalAnalysis: "īrrŭmābō̆ future indicative active 1st singular",
				},
			},
		},
		{
			name: "unknown translation",
			input: bytes.NewReader([]byte(`
1	1	1	Pedicabo	v1 	pedico	pēdīco, as, are	4	to perform anal intercourse; to commit sodomy with;	pēdīcābō̆ future indicative active 1st singular
2	1	2	ego	p11	ego	ĕgō̆, mei, pron.	14846	moi, me	ĕgō̆ masculine nominative singular
3	1	3	vos	p32	vos	vōs, uestrum, pl. pron.	2402	you (pl.), ye;	vōs masculine accusative plural
1	1	1	et		unknown
2	1	2	irrumabo	v1 	inrumo	īnrŭmo, as, are	6	to extend the breast to, to give suck;	īrrŭmābō̆ future indicative active 1st singular
`)),
			expected: &[]domain.Word{
				{
					WordIndex:                 1,
					SentenceIndex:             1,
					WordIndexInSentence:       1,
					OriginalForm:              "Pedicabo",
					LemmaRaw:                  "pedico",
					LemmaRich:                 "pēdīco, as, are",
					Translation:               "to perform anal intercourse; to commit sodomy with;",
					FrequencyInLASLA:          4,
					Tag:                       "v1",
					MorphoSyntacticalAnalysis: "pēdīcābō̆ future indicative active 1st singular",
				},
				{
					WordIndex:                 2,
					SentenceIndex:             1,
					WordIndexInSentence:       2,
					OriginalForm:              "ego",
					LemmaRaw:                  "ego",
					LemmaRich:                 "ĕgō̆, mei, pron.",
					Translation:               "moi, me",
					FrequencyInLASLA:          14846,
					Tag:                       "p11",
					MorphoSyntacticalAnalysis: "ĕgō̆ masculine nominative singular",
				},
				{
					WordIndex:                 3,
					SentenceIndex:             1,
					WordIndexInSentence:       3,
					OriginalForm:              "vos",
					LemmaRaw:                  "vos",
					LemmaRich:                 "vōs, uestrum, pl. pron.",
					Translation:               "you (pl.), ye;",
					FrequencyInLASLA:          2402,
					Tag:                       "p32",
					MorphoSyntacticalAnalysis: "vōs masculine accusative plural",
				},
				{
					WordIndex:           4,
					SentenceIndex:       2,
					WordIndexInSentence: 1,
					OriginalForm:        "et",
					Translation:         "unknown",
				},
				{
					WordIndex:                 5,
					SentenceIndex:             2,
					WordIndexInSentence:       2,
					OriginalForm:              "irrumabo",
					LemmaRaw:                  "inrumo",
					LemmaRich:                 "īnrŭmo, as, are",
					Translation:               "to extend the breast to, to give suck;",
					FrequencyInLASLA:          6,
					Tag:                       "v1",
					MorphoSyntacticalAnalysis: "īrrŭmābō̆ future indicative active 1st singular",
				},
			},
		},
		{
			name: "empty line in between",
			input: bytes.NewReader([]byte(`
1	1	1	Pedicabo	v1 	pedico	pēdīco, as, are	4	to perform anal intercourse; to commit sodomy with;	pēdīcābō̆ future indicative active 1st singular
2	1	2	ego	p11	ego	ĕgō̆, mei, pron.	14846	moi, me	ĕgō̆ masculine nominative singular

1	1	1	et	d   (c  )	et	ĕt, conj. adv.	42726	and, and even; also, even; (et ... et = both ... and);	ĕt
2	1	2	irrumabo	v1 	inrumo	īnrŭmo, as, are	6	to extend the breast to, to give suck;	īrrŭmābō̆ future indicative active 1st singular
`)),
			expected: &[]domain.Word{
				{
					WordIndex:                 1,
					SentenceIndex:             1,
					WordIndexInSentence:       1,
					OriginalForm:              "Pedicabo",
					LemmaRaw:                  "pedico",
					LemmaRich:                 "pēdīco, as, are",
					Translation:               "to perform anal intercourse; to commit sodomy with;",
					FrequencyInLASLA:          4,
					Tag:                       "v1",
					MorphoSyntacticalAnalysis: "pēdīcābō̆ future indicative active 1st singular",
				},
				{
					WordIndex:                 2,
					SentenceIndex:             1,
					WordIndexInSentence:       2,
					OriginalForm:              "ego",
					LemmaRaw:                  "ego",
					LemmaRich:                 "ĕgō̆, mei, pron.",
					Translation:               "moi, me",
					FrequencyInLASLA:          14846,
					Tag:                       "p11",
					MorphoSyntacticalAnalysis: "ĕgō̆ masculine nominative singular",
				},
				{
					WordIndex:                 3,
					SentenceIndex:             2,
					WordIndexInSentence:       1,
					OriginalForm:              "et",
					LemmaRaw:                  "et",
					LemmaRich:                 "ĕt, conj. adv.",
					Translation:               "and, and even; also, even; (et ... et = both ... and);",
					FrequencyInLASLA:          42726,
					Tag:                       "d   (c  )",
					MorphoSyntacticalAnalysis: "ĕt",
				},
				{
					WordIndex:                 4,
					SentenceIndex:             2,
					WordIndexInSentence:       2,
					OriginalForm:              "irrumabo",
					LemmaRaw:                  "inrumo",
					LemmaRich:                 "īnrŭmo, as, are",
					Translation:               "to extend the breast to, to give suck;",
					FrequencyInLASLA:          6,
					Tag:                       "v1",
					MorphoSyntacticalAnalysis: "īrrŭmābō̆ future indicative active 1st singular",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := mapToWords(test.input)
			assert.Equal(t, test.expected, output)
		})
	}
}
