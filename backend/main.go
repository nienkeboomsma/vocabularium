package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"

	"github.com/nienkeboomsma/collatinus/textprocessor/infrastructure/collatinus"
	"github.com/nienkeboomsma/collatinus/textprocessor/ports/driven"
)

func handler(tp driven.TextProcessor) http.HandlerFunc {
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

		words, err := tp.Process(uploadedData)

		f, err := os.Create("/data/output.json")
		if err != nil {
			http.Error(w, "Failed to create output file", http.StatusBadRequest)
			return
		}
		defer f.Close()

		encoder := json.NewEncoder(f)

		err = encoder.Encode(words)
		if err != nil {
			http.Error(w, "Failed to encode data into output.json", http.StatusBadRequest)
			return
		}

		w.Write([]byte("Done!"))
	}
}

func main() {
	language := cmp.Or(os.Getenv("COLLATINUS_LANGUAGE"), "en")
	validLanguages := []string{"ca", "de", "en", "es", "eu", "fr", "gl", "it", "nl", "pt"}

	if !slices.Contains(validLanguages, language) {
		log.Fatal(`COLLATINUS_LANGUAGE must be one of "ca", "de", "en", "es", "eu", "fr", "gl", "it", "nl" or "pt"`)
	}

	tp := collatinus.NewTextProcessor(language)

	http.HandleFunc("/lemmatise", handler(tp))

	fmt.Println("Listening at :6666")

	err := http.ListenAndServe(":6666", nil)
	if err != nil {
		log.Fatal("Failed to start server: " + err.Error())
	}

}
