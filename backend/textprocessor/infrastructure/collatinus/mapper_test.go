package collatinus

import (
	"bytes"
	"io"
	"testing"

	"github.com/google/uuid"
	"github.com/nienkeboomsma/vocabularium/database"
	"github.com/nienkeboomsma/vocabularium/domain"
	"github.com/stretchr/testify/assert"
)

func TestMapper(t *testing.T) {
	tests := []struct {
		name              string
		input             io.Reader
		expectedWorkWords *[]domain.WorkWord
		expectedWords     *map[uuid.UUID]domain.Word
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
			expectedWorkWords: &[]domain.WorkWord{
				{
					WordID:                    database.StringToUUID("pēdīco, as, are to perform anal intercourse; to commit sodomy with;"),
					WordIndex:                 1,
					SentenceIndex:             1,
					OriginalForm:              "Pedicabo",
					Tag:                       "v1",
					MorphoSyntacticalAnalysis: "pēdīcābō̆ future indicative active 1st singular",
				},
				{
					WordID:                    database.StringToUUID("ĕgō̆, mei, pron. moi, me"),
					WordIndex:                 2,
					SentenceIndex:             1,
					OriginalForm:              "ego",
					Tag:                       "p11",
					MorphoSyntacticalAnalysis: "ĕgō̆ masculine nominative singular",
				},
				{
					WordID:                    database.StringToUUID("vōs, uestrum, pl. pron. you (pl.), ye;"),
					WordIndex:                 3,
					SentenceIndex:             1,
					OriginalForm:              "vos",
					Tag:                       "p32",
					MorphoSyntacticalAnalysis: "vōs masculine accusative plural",
				},
				{
					WordID:                    database.StringToUUID("ĕt, conj. adv. and, and even; also, even; (et ... et = both ... and);"),
					WordIndex:                 4,
					SentenceIndex:             2,
					OriginalForm:              "et",
					Tag:                       "d   (c  )",
					MorphoSyntacticalAnalysis: "ĕt",
				},
				{
					WordID:                    database.StringToUUID("īnrŭmo, as, are to extend the breast to, to give suck;"),
					WordIndex:                 5,
					SentenceIndex:             2,
					OriginalForm:              "irrumabo",
					Tag:                       "v1",
					MorphoSyntacticalAnalysis: "īrrŭmābō̆ future indicative active 1st singular",
				},
			},
			expectedWords: &map[uuid.UUID]domain.Word{
				database.StringToUUID("pēdīco, as, are to perform anal intercourse; to commit sodomy with;"): {
					ID:               database.StringToUUID("pēdīco, as, are to perform anal intercourse; to commit sodomy with;"),
					LemmaRaw:         "pedico",
					LemmaRich:        "pēdīco, as, are",
					Translation:      "to perform anal intercourse; to commit sodomy with;",
					FrequencyInLASLA: 4,
				},
				database.StringToUUID("ĕgō̆, mei, pron. moi, me"): {
					ID:               database.StringToUUID("ĕgō̆, mei, pron. moi, me"),
					LemmaRaw:         "ego",
					LemmaRich:        "ĕgō̆, mei, pron.",
					Translation:      "moi, me",
					FrequencyInLASLA: 14846,
				},
				database.StringToUUID("vōs, uestrum, pl. pron. you (pl.), ye;"): {
					ID:               database.StringToUUID("vōs, uestrum, pl. pron. you (pl.), ye;"),
					LemmaRaw:         "vos",
					LemmaRich:        "vōs, uestrum, pl. pron.",
					Translation:      "you (pl.), ye;",
					FrequencyInLASLA: 2402,
				},
				database.StringToUUID("ĕt, conj. adv. and, and even; also, even; (et ... et = both ... and);"): {
					ID:               database.StringToUUID("ĕt, conj. adv. and, and even; also, even; (et ... et = both ... and);"),
					LemmaRaw:         "et",
					LemmaRich:        "ĕt, conj. adv.",
					Translation:      "and, and even; also, even; (et ... et = both ... and);",
					FrequencyInLASLA: 42726,
				},
				database.StringToUUID("īnrŭmo, as, are to extend the breast to, to give suck;"): {
					ID:               database.StringToUUID("īnrŭmo, as, are to extend the breast to, to give suck;"),
					LemmaRaw:         "inrumo",
					LemmaRich:        "īnrŭmo, as, are",
					Translation:      "to extend the breast to, to give suck;",
					FrequencyInLASLA: 6,
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
			expectedWorkWords: &[]domain.WorkWord{
				{
					WordID:                    database.StringToUUID("pēdīco, as, are to perform anal intercourse; to commit sodomy with;"),
					WordIndex:                 1,
					SentenceIndex:             1,
					OriginalForm:              "Pedicabo",
					Tag:                       "v1",
					MorphoSyntacticalAnalysis: "pēdīcābō̆ future indicative active 1st singular",
				},
				{
					WordID:                    database.StringToUUID("ĕgō̆, mei, pron. moi, me"),
					WordIndex:                 2,
					SentenceIndex:             1,
					OriginalForm:              "ego",
					Tag:                       "p11",
					MorphoSyntacticalAnalysis: "ĕgō̆ masculine nominative singular",
				},
				{
					WordID:                    database.StringToUUID("vōs, uestrum, pl. pron. you (pl.), ye;"),
					WordIndex:                 3,
					SentenceIndex:             1,
					OriginalForm:              "vos",
					Tag:                       "p32",
					MorphoSyntacticalAnalysis: "vōs masculine accusative plural",
				},
				{
					WordID:                    database.StringToUUID("īnrŭmo, as, are to extend the breast to, to give suck;"),
					WordIndex:                 5,
					SentenceIndex:             2,
					OriginalForm:              "irrumabo",
					Tag:                       "v1",
					MorphoSyntacticalAnalysis: "īrrŭmābō̆ future indicative active 1st singular",
				},
			},
			expectedWords: &map[uuid.UUID]domain.Word{
				database.StringToUUID("pēdīco, as, are to perform anal intercourse; to commit sodomy with;"): {
					ID:               database.StringToUUID("pēdīco, as, are to perform anal intercourse; to commit sodomy with;"),
					LemmaRaw:         "pedico",
					LemmaRich:        "pēdīco, as, are",
					Translation:      "to perform anal intercourse; to commit sodomy with;",
					FrequencyInLASLA: 4,
				},
				database.StringToUUID("ĕgō̆, mei, pron. moi, me"): {
					ID:               database.StringToUUID("ĕgō̆, mei, pron. moi, me"),
					LemmaRaw:         "ego",
					LemmaRich:        "ĕgō̆, mei, pron.",
					Translation:      "moi, me",
					FrequencyInLASLA: 14846,
				},
				database.StringToUUID("vōs, uestrum, pl. pron. you (pl.), ye;"): {
					ID:               database.StringToUUID("vōs, uestrum, pl. pron. you (pl.), ye;"),
					LemmaRaw:         "vos",
					LemmaRich:        "vōs, uestrum, pl. pron.",
					Translation:      "you (pl.), ye;",
					FrequencyInLASLA: 2402,
				},
				database.StringToUUID("īnrŭmo, as, are to extend the breast to, to give suck;"): {
					ID:               database.StringToUUID("īnrŭmo, as, are to extend the breast to, to give suck;"),
					LemmaRaw:         "inrumo",
					LemmaRich:        "īnrŭmo, as, are",
					Translation:      "to extend the breast to, to give suck;",
					FrequencyInLASLA: 6,
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
			expectedWorkWords: &[]domain.WorkWord{

				{
					WordID:                    database.StringToUUID("pēdīco, as, are to perform anal intercourse; to commit sodomy with;"),
					WordIndex:                 1,
					SentenceIndex:             1,
					OriginalForm:              "Pedicabo",
					Tag:                       "v1",
					MorphoSyntacticalAnalysis: "pēdīcābō̆ future indicative active 1st singular",
				},
				{
					WordID:                    database.StringToUUID("ĕgō̆, mei, pron. moi, me"),
					WordIndex:                 2,
					SentenceIndex:             1,
					OriginalForm:              "ego",
					Tag:                       "p11",
					MorphoSyntacticalAnalysis: "ĕgō̆ masculine nominative singular",
				},
				{
					WordID:                    database.StringToUUID("vōs, uestrum, pl. pron. you (pl.), ye;"),
					WordIndex:                 3,
					SentenceIndex:             1,
					OriginalForm:              "vos",
					Tag:                       "p32",
					MorphoSyntacticalAnalysis: "vōs masculine accusative plural",
				},
				{
					WordIndex:     4,
					SentenceIndex: 2,
					OriginalForm:  "et",
				},
				{
					WordID:                    database.StringToUUID("īnrŭmo, as, are to extend the breast to, to give suck;"),
					WordIndex:                 5,
					SentenceIndex:             2,
					OriginalForm:              "irrumabo",
					Tag:                       "v1",
					MorphoSyntacticalAnalysis: "īrrŭmābō̆ future indicative active 1st singular",
				},
			},
			expectedWords: &map[uuid.UUID]domain.Word{
				database.StringToUUID("pēdīco, as, are to perform anal intercourse; to commit sodomy with;"): {
					ID:               database.StringToUUID("pēdīco, as, are to perform anal intercourse; to commit sodomy with;"),
					LemmaRaw:         "pedico",
					LemmaRich:        "pēdīco, as, are",
					Translation:      "to perform anal intercourse; to commit sodomy with;",
					FrequencyInLASLA: 4,
				},
				database.StringToUUID("ĕgō̆, mei, pron. moi, me"): {
					ID:               database.StringToUUID("ĕgō̆, mei, pron. moi, me"),
					LemmaRaw:         "ego",
					LemmaRich:        "ĕgō̆, mei, pron.",
					Translation:      "moi, me",
					FrequencyInLASLA: 14846,
				},
				database.StringToUUID("vōs, uestrum, pl. pron. you (pl.), ye;"): {
					ID:               database.StringToUUID("vōs, uestrum, pl. pron. you (pl.), ye;"),
					LemmaRaw:         "vos",
					LemmaRich:        "vōs, uestrum, pl. pron.",
					Translation:      "you (pl.), ye;",
					FrequencyInLASLA: 2402,
				},
				database.StringToUUID("īnrŭmo, as, are to extend the breast to, to give suck;"): {
					ID:               database.StringToUUID("īnrŭmo, as, are to extend the breast to, to give suck;"),
					LemmaRaw:         "inrumo",
					LemmaRich:        "īnrŭmo, as, are",
					Translation:      "to extend the breast to, to give suck;",
					FrequencyInLASLA: 6,
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
			expectedWorkWords: &[]domain.WorkWord{
				{
					WordID:                    database.StringToUUID("pēdīco, as, are to perform anal intercourse; to commit sodomy with;"),
					WordIndex:                 1,
					SentenceIndex:             1,
					OriginalForm:              "Pedicabo",
					Tag:                       "v1",
					MorphoSyntacticalAnalysis: "pēdīcābō̆ future indicative active 1st singular",
				},
				{
					WordID:                    database.StringToUUID("ĕgō̆, mei, pron. moi, me"),
					WordIndex:                 2,
					SentenceIndex:             1,
					OriginalForm:              "ego",
					Tag:                       "p11",
					MorphoSyntacticalAnalysis: "ĕgō̆ masculine nominative singular",
				},
				{
					WordID:                    database.StringToUUID("ĕt, conj. adv. and, and even; also, even; (et ... et = both ... and);"),
					WordIndex:                 3,
					SentenceIndex:             2,
					OriginalForm:              "et",
					Tag:                       "d   (c  )",
					MorphoSyntacticalAnalysis: "ĕt",
				},
				{
					WordID:                    database.StringToUUID("īnrŭmo, as, are to extend the breast to, to give suck;"),
					WordIndex:                 4,
					SentenceIndex:             2,
					OriginalForm:              "irrumabo",
					Tag:                       "v1",
					MorphoSyntacticalAnalysis: "īrrŭmābō̆ future indicative active 1st singular",
				},
			},
			expectedWords: &map[uuid.UUID]domain.Word{
				database.StringToUUID("pēdīco, as, are to perform anal intercourse; to commit sodomy with;"): {
					ID:               database.StringToUUID("pēdīco, as, are to perform anal intercourse; to commit sodomy with;"),
					LemmaRaw:         "pedico",
					LemmaRich:        "pēdīco, as, are",
					Translation:      "to perform anal intercourse; to commit sodomy with;",
					FrequencyInLASLA: 4,
				},
				database.StringToUUID("ĕgō̆, mei, pron. moi, me"): {
					ID:               database.StringToUUID("ĕgō̆, mei, pron. moi, me"),
					LemmaRaw:         "ego",
					LemmaRich:        "ĕgō̆, mei, pron.",
					Translation:      "moi, me",
					FrequencyInLASLA: 14846,
				},
				database.StringToUUID("ĕt, conj. adv. and, and even; also, even; (et ... et = both ... and);"): {
					ID:               database.StringToUUID("ĕt, conj. adv. and, and even; also, even; (et ... et = both ... and);"),
					LemmaRaw:         "et",
					LemmaRich:        "ĕt, conj. adv.",
					Translation:      "and, and even; also, even; (et ... et = both ... and);",
					FrequencyInLASLA: 42726,
				},
				database.StringToUUID("īnrŭmo, as, are to extend the breast to, to give suck;"): {
					ID:               database.StringToUUID("īnrŭmo, as, are to extend the breast to, to give suck;"),
					LemmaRaw:         "inrumo",
					LemmaRich:        "īnrŭmo, as, are",
					Translation:      "to extend the breast to, to give suck;",
					FrequencyInLASLA: 6,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			workWords, words, _ := mapToWords(test.input)
			assert.Equal(t, test.expectedWorkWords, workWords)
			assert.Equal(t, test.expectedWords, words)
		})
	}
}
