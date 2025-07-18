package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	t "text/template"

	"github.com/google/uuid"
	"github.com/nienkeboomsma/collatinus/api/infrastructure/template"
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
		template.GetWordListTemplate("Frequency list"),
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
		template.GetWordListTemplate("Frequency list"),
		a.workRepository.GetByID,
		a.wordRepository.GetFrequencyListByWorkID,
	)
}

func (a *API) GetGlossaryByWork() http.HandlerFunc {
	return handleWordList(
		template.GetWordListTemplate("Glossary"),
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

		workListTemplate := template.GetWorkListTemplate()

		tmpl, err := t.New("works").Parse(workListTemplate)
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
			errorHTML := template.GetFailedWorkUploadTemplate("Failed to get uploaded file from the request", err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, errorHTML)
			return
		}
		defer file.Close()

		uploadedData, err := io.ReadAll(file)
		if err != nil {
			errorHTML := template.GetFailedWorkUploadTemplate("Failed to read the contents of the uploaded file", err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, errorHTML)
			return
		}

		author := domain.Author{
			ID:   database.StringToUUID(r.FormValue("author")),
			Name: r.FormValue("author"),
		}
		work := domain.Work{
			ID:    database.StringToUUID(fmt.Sprintf("%s_%s", r.FormValue("author"), r.FormValue("title"))),
			Title: r.FormValue("title"),
		}

		workWords, words, logs, err := a.textProcessor.Process(uploadedData)
		if err != nil {
			errorHTML := template.GetFailedWorkUploadTemplate("Failed to process the uploaded text", err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, errorHTML)
			return
		}

		err = a.workPersister.Persist(r.Context(), author, work, words, workWords)
		if err != nil {
			errorHTML := template.GetFailedWorkUploadTemplate("Failed to save the work in the database", err)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, errorHTML)
			return
		}

		successTemplate := template.GetSuccessfulWorkUploadTemplate()

		tmpl, err := t.New("works").Parse(successTemplate)
		if err != nil {
			http.Error(w, "Failed to allocate HTML template", http.StatusInternalServerError)
			return
		}

		for _, log := range logs {
			log = t.HTMLEscapeString(log)
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err = tmpl.Execute(w, template.UploadSuccessData{Logs: logs})
		if err != nil {
			http.Error(w, "Failed to render HTML template", http.StatusInternalServerError)
			return
		}
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

func (a *API) Upload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		html := template.GetUploadTemplate()
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
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

		pageData := template.WordListPageData{
			Title:  work.Title,
			Author: work.Author.Name,
			Words:  words,
		}

		tmpl, err := t.New("words").Parse(htmlTemplate)
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
