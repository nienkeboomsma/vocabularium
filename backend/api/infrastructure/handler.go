package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/nienkeboomsma/collatinus/database"
	"github.com/nienkeboomsma/collatinus/domain"
	"github.com/nienkeboomsma/collatinus/textprocessor/ports/driven"
	"github.com/nienkeboomsma/collatinus/workpersister/ports/driving"
)

type API struct {
	textProcessor driven.TextProcessor
	workPersister driving.WorkPersister
}

func NewAPI(tp driven.TextProcessor, wp driving.WorkPersister) *API {
	return &API{
		textProcessor: tp,
		workPersister: wp,
	}
}

func (a *API) Lemmatise() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Missing file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		uploadedData, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read uploaded file", http.StatusBadRequest)
			return
		}

		author := domain.Author{
			ID:   database.StringToUUID(r.FormValue("author")),
			Name: r.FormValue("author"),
		}
		work := domain.Work{
			ID:    database.StringToUUID(fmt.Sprintf("%s_%s", r.FormValue("author"), r.FormValue("title"))),
			Title: r.FormValue("title"),
			Type:  domain.WorkType(r.FormValue("type")),
		}

		workWords, words, err := a.textProcessor.Process(uploadedData)

		err = a.workPersister.Persist(r.Context(), author, work, words, workWords)
		if err != nil {
			http.Error(w, "Failed to save work", http.StatusBadRequest)
			return
		}

		w.Write([]byte("Done!"))
	}
}
