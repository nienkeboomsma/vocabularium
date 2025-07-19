package api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	t "text/template"

	"github.com/google/uuid"
	"github.com/nienkeboomsma/vocabularium/api/infrastructure/template"
	"github.com/nienkeboomsma/vocabularium/database"
	"github.com/nienkeboomsma/vocabularium/domain"
	repositories "github.com/nienkeboomsma/vocabularium/repositories/ports/driving"
	"github.com/nienkeboomsma/vocabularium/textprocessor/ports/driven"
	"github.com/nienkeboomsma/vocabularium/workpersister/ports/driving"
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

func (a *API) DeleteWork() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			fmt.Println(id)
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		err = a.workRepository.Delete(r.Context(), id)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Failed to delete work", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (a *API) GetFrequencyList() http.HandlerFunc {
	return handleWordList(
		template.GetWordListTemplate("Frequency list", "📈"),
		func(r *http.Request) (uuid.UUID, error) { return uuid.UUID{}, nil },
		func(ctx context.Context, id uuid.UUID) (domain.Work, error) { return domain.Work{}, nil },
		func(ctx context.Context, id uuid.UUID) (*[]domain.WordInWork, error) {
			return a.wordRepository.GetFrequencyList(ctx)
		},
	)
}

func (a *API) GetFrequencyListByAuthor() http.HandlerFunc {
	return handleWordList(
		template.GetWordListTemplate("Frequency list", "📈"),
		func(r *http.Request) (uuid.UUID, error) { return uuid.Parse(r.PathValue("id")) },
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
		template.GetWordListTemplate("Frequency list", "📈"),
		func(r *http.Request) (uuid.UUID, error) { return uuid.Parse(r.PathValue("id")) },
		a.workRepository.GetByID,
		a.wordRepository.GetFrequencyListByWorkID,
	)
}

func (a *API) GetGlossaryByWork() http.HandlerFunc {
	return handleWordList(
		template.GetWordListTemplate("Glossary", "📖"),
		func(r *http.Request) (uuid.UUID, error) { return uuid.Parse(r.PathValue("id")) },
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

		useTemplate(w, template.GetWorkListTemplate(), works)
	}
}

func (a *API) Lemmatise() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workID := database.StringToUUID(fmt.Sprintf("%s_%s", r.FormValue("author"), r.FormValue("title")))
		work, err := a.workRepository.GetByID(r.Context(), workID)
		if err == nil {
			w.WriteHeader(http.StatusBadRequest)
			useTemplate(w, template.GetFailedWorkUploadTemplate(), template.UploadFailedData{
				Message: "A work with this author and title already exists. Please remove it and try again.",
				Error:   "",
			})
			return
		}

		file, _, err := r.FormFile("file")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			useTemplate(w, template.GetFailedWorkUploadTemplate(), template.UploadFailedData{
				Message: "Failed to get uploaded file from the request",
				Error:   err.Error(),
			})
			return
		}
		defer file.Close()

		uploadedData, err := io.ReadAll(file)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			useTemplate(w, template.GetFailedWorkUploadTemplate(), template.UploadFailedData{
				Message: "Failed to read the contents of the uploaded file",
				Error:   err.Error(),
			})
			return
		}

		author := domain.Author{
			ID:   database.StringToUUID(r.FormValue("author")),
			Name: r.FormValue("author"),
		}
		work = domain.Work{
			ID:    workID,
			Title: r.FormValue("title"),
		}

		workWords, words, logs, err := a.textProcessor.Process(uploadedData)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			useTemplate(w, template.GetFailedWorkUploadTemplate(), template.UploadFailedData{
				Message: "Failed to process the uploaded text",
				Error:   err.Error(),
			})
			return
		}

		err = a.workPersister.Persist(r.Context(), author, work, words, workWords)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			useTemplate(w, template.GetFailedWorkUploadTemplate(), template.UploadFailedData{
				Message: "Failed to save the work in the database",
				Error:   err.Error(),
			})
			return
		}

		for i, log := range logs {
			logs[i] = t.HTMLEscapeString(log)
		}

		w.WriteHeader(http.StatusBadRequest)
		useTemplate(w, template.GetSuccessfulWorkUploadTemplate(), template.UploadSuccessData{Logs: logs})
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
	idCallback func(r *http.Request) (uuid.UUID, error),
	workCallback func(ctx context.Context, id uuid.UUID) (domain.Work, error),
	wordCallback func(ctx context.Context, id uuid.UUID) (*[]domain.WordInWork, error),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := idCallback(r)
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

		useTemplate(w, htmlTemplate, template.WordListPageData{
			Title:  work.Title,
			Author: work.Author.Name,
			Words:  words,
		})
	}
}

func useTemplate(w http.ResponseWriter, template string, data any) {
	tmpl, err := t.New("works").Parse(template)
	if err != nil {
		http.Error(w, "Failed to allocate HTML template", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		http.Error(w, "Failed to render HTML template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	buf.WriteTo(w)
}
