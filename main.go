package main

import (
	"cmp"
	"fmt"
	"log"
	"net/http"
	"os"

	api "github.com/nienkeboomsma/vocabularium/api/infrastructure"
	"github.com/nienkeboomsma/vocabularium/database"
	repositories "github.com/nienkeboomsma/vocabularium/repositories/infrastructure/postgres"
	"github.com/nienkeboomsma/vocabularium/textprocessor/infrastructure/collatinus"
	"github.com/nienkeboomsma/vocabularium/workpersister/infrastructure/postgres"
)

func main() {
	language := cmp.Or(os.Getenv("COLLATINUS_LANGUAGE"), "en")

	tp, err := collatinus.NewTextProcessor(language)
	if err != nil {
		log.Fatal("Failed to create text processor: " + err.Error())
	}

	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatal("DB_URL must be provided in .env file")
	}

	db, err := database.New(dbURL)

	err = db.RunMigrations()
	if err != nil {
		log.Fatal("Failed to run database migrations: " + err.Error())
	}

	authorRepository := repositories.NewAuthorRepository(db)
	workRepository := repositories.NewWorkRepository(db)
	wordRepository := repositories.NewWordRepository(db)
	workWordRepository := repositories.NewWorkWordRepository()

	wp := postgres.NewWorkPersister(db, authorRepository, workRepository, wordRepository, workWordRepository)

	api := api.NewAPI(tp, wp, authorRepository, wordRepository, workRepository)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /frequency-list/{id}/{skipKnown}", api.GetFrequencyListByWork())
	mux.HandleFunc("GET /frequency-list-author/{id}/{skipKnown}", api.GetFrequencyListByAuthor())
	mux.HandleFunc("GET /glossary/{id}/{skipKnown}", api.GetGlossaryByWork())
	mux.HandleFunc("GET /upload", api.Upload())
	mux.HandleFunc("GET /", api.GetWorks())

	mux.HandleFunc("POST /lemmatise", api.Lemmatise())
	mux.HandleFunc("POST /delete/{id}", api.DeleteWork())
	mux.HandleFunc("POST /toggle-known-status/{id}", api.ToggleKnownStatus())

	// TODO: set ports via .env
	fmt.Println("Listening at :4321")

	err = http.ListenAndServe(":4321", mux)
	if err != nil {
		log.Fatal("Failed to start server: " + err.Error())
	}

}
