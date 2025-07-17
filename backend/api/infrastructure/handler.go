package api

import (
	"fmt"
	"io"
	"net/http"
	"text/template"

	"github.com/nienkeboomsma/collatinus/database"
	"github.com/nienkeboomsma/collatinus/domain"
	repositories "github.com/nienkeboomsma/collatinus/repositories/ports/driving"
	"github.com/nienkeboomsma/collatinus/textprocessor/ports/driven"
	"github.com/nienkeboomsma/collatinus/workpersister/ports/driving"
)

type API struct {
	textProcessor  driven.TextProcessor
	workPersister  driving.WorkPersister
	workRepository repositories.WorkRepository
}

func NewAPI(
	tp driven.TextProcessor,
	wp driving.WorkPersister,
	wr repositories.WorkRepository,
) *API {
	return &API{
		textProcessor:  tp,
		workPersister:  wp,
		workRepository: wr,
	}
}

func (a *API) GetFrequencyListByWork() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
} // for a workID, return list of words per work, grouped by wordID, ordered by count(wordID)

func (a *API) GetFrequencyListByAuthor() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
} // for an authorID, return list of words across all works, grouped by wordID, ordered by count(wordID)

func (a *API) GetGlossaryByWork() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
} // for a workID, return list of words per work, ordered by wordIndex

func (a *API) GetWorks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		works, err := a.workRepository.Get(r.Context())
		if err != nil {
			http.Error(w, "Failed to retrieve works", http.StatusBadRequest)
			return
		}

		workListTemplate := `
<!DOCTYPE html>
<html>
<head><title>Works</title></head>
<body>
	<h1>Works</h1>
	<table>
		<thead>
			<tr>
				<th colspan="2">Author</th>
				<th colspan="3">Title</th>
				<th>Type</th>
			</tr>
		</thead>
		<tbody>
		{{range .}}
			<tr>
				<td>{{.Author.Name}}</td>
				<td><a href="http://localhost:4321/frequency-list-author/{{.Author.ID}}">frequency list</a></td>
				<td>{{.Title}}</td>
				<td><a href="http://localhost:4321/frequency-list/{{.ID}}">frequency list</a></td>
				<td><a href="http://localhost:4321/glossary/{{.ID}}">glossary</a></td>
				<td>{{.Type}}</td>
			</tr>
		{{else}}
			<tr><td colspan="3">No works to display</td></tr>
		{{end}}
		</tbody>
	</table>
</body>
</html>`

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

		w.Write([]byte("Done!"))
	}
}

func (a *API) UpdateWordStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
} // for a wordID, update the known status to either 'true' or 'false'
