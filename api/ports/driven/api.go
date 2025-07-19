package driven

import "net/http"

type API interface {
	DeleteWork() http.HandlerFunc
	GetFrequencyList() http.HandlerFunc
	GetFrequencyListByWork() http.HandlerFunc
	GetFrequencyListByAuthor() http.HandlerFunc
	GetGlossaryByWork() http.HandlerFunc
	GetWorks() http.HandlerFunc
	Lemmatise() http.HandlerFunc
	ToggleKnownStatus() http.HandlerFunc
	Upload() http.HandlerFunc
}
