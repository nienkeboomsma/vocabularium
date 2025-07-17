package driven

import "net/http"

type API interface {
	GetFrequencyListByWork() http.HandlerFunc   // for a workID, return list of words per work, grouped by wordID, ordered by count(wordID)
	GetFrequencyListByAuthor() http.HandlerFunc // for an authorID, return list of words across all works, grouped by wordID, ordered by count(wordID)
	GetGlossaryByWork() http.HandlerFunc        // for a workID, return list of words per work, ordered by wordIndex
	GetWorks() http.HandlerFunc                 // return list of works and their authors
	Lemmatise() http.HandlerFunc
	UpdateWordStatus() http.HandlerFunc // for a wordID, update the known status to either 'true' or 'false'
}
