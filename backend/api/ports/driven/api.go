package driven

import "net/http"

type API interface {
	Lemmatise() http.HandlerFunc
}
