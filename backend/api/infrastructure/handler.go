package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/google/uuid"
	"github.com/nienkeboomsma/collatinus/database"
	"github.com/nienkeboomsma/collatinus/domain"
	repositories "github.com/nienkeboomsma/collatinus/repositories/ports/driving"
	"github.com/nienkeboomsma/collatinus/textprocessor/ports/driven"
	"github.com/nienkeboomsma/collatinus/workpersister/ports/driving"
)

type API struct {
	textProcessor    driven.TextProcessor
	workPersister    driving.WorkPersister
	authorRepository repositories.AuthorRepository
	wordRepository   repositories.WordRepository
	workRepository   repositories.WorkRepository
}

func NewAPI(
	tp driven.TextProcessor,
	wp driving.WorkPersister,
	authorRepository repositories.AuthorRepository,
	wordRepository repositories.WordRepository,
	workRepository repositories.WorkRepository,
) *API {
	return &API{
		textProcessor:    tp,
		workPersister:    wp,
		authorRepository: authorRepository,
		wordRepository:   wordRepository,
		workRepository:   workRepository,
	}
}

func (a *API) GetFrequencyListByAuthor() http.HandlerFunc {
	return handleWordList(
		getWordListTemplate("Frequency list"),
		func(ctx context.Context, id uuid.UUID) (domain.Work, error) {
			author, err := a.authorRepository.GetByID(ctx, id)
			if err != nil {
				return domain.Work{}, err
			}

			return domain.Work{Author: domain.Author{
				Name: author.Name,
			}}, nil
		},
		a.wordRepository.GetFrequencyListByAuthorID,
	)
}

func (a *API) GetFrequencyListByWork() http.HandlerFunc {
	return handleWordList(
		getWordListTemplate("Frequency list"),
		a.workRepository.GetByID,
		a.wordRepository.GetFrequencyListByWorkID,
	)
}

func (a *API) GetGlossaryByWork() http.HandlerFunc {
	return handleWordList(
		getWordListTemplate("Glossary"),
		a.workRepository.GetByID,
		a.wordRepository.GetGlossaryByWorkID,
	)
}

func (a *API) GetWorks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		works, err := a.workRepository.Get(r.Context())
		if err != nil {
			http.Error(w, "Failed to retrieve works", http.StatusBadRequest)
			return
		}

		workListTemplate := getWorkListTemplate()

		tmpl, err := template.New("works").Parse(workListTemplate)
		if err != nil {
			http.Error(w, "Failed to allocate HTML template", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = tmpl.Execute(w, works)
		if err != nil {
			http.Error(w, "Failed to render HTML template", http.StatusInternalServerError)
			return
		}
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

		w.WriteHeader(http.StatusOK)
	}
}

func (a *API) ToggleKnownStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		_, err = a.wordRepository.ToggleKnownStatus(r.Context(), id)
		if err != nil {
			http.Error(w, "Failed to save updated known status", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func handleWordList(
	htmlTemplate string,
	workCallback func(ctx context.Context, id uuid.UUID) (domain.Work, error),
	wordCallback func(ctx context.Context, id uuid.UUID) (*[]domain.WordInWork, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		work, err := workCallback(r.Context(), id)
		if err != nil {
			http.Error(w, "Failed to retrieve work", http.StatusBadRequest)
			return
		}

		words, err := wordCallback(r.Context(), id)
		if err != nil {
			http.Error(w, "Failed to retrieve words", http.StatusBadRequest)
			return
		}

		if r.PathValue("skipKnown") == "true" {
			filteredWords := []domain.WordInWork{}

			for _, word := range *words {
				if word.Known {
					continue
				}

				filteredWords = append(filteredWords, word)
			}

			words = &filteredWords
		}

		pageData := WordListPageData{
			Title:  work.Title,
			Author: work.Author.Name,
			Words:  words,
		}

		tmpl, err := template.New("words").Parse(htmlTemplate)
		if err != nil {
			http.Error(w, "Failed to allocate HTML template", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = tmpl.Execute(w, pageData)
		if err != nil {
			http.Error(w, "Failed to render HTML template", http.StatusInternalServerError)
			return
		}
	}
}
