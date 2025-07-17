package main

import (
	"cmp"
	"fmt"
	"log"
	"net/http"
	"os"

	api "github.com/nienkeboomsma/collatinus/api/infrastructure"
	"github.com/nienkeboomsma/collatinus/database"
	repositories "github.com/nienkeboomsma/collatinus/repositories/infrastructure/postgres"
	"github.com/nienkeboomsma/collatinus/textprocessor/infrastructure/collatinus"
	"github.com/nienkeboomsma/collatinus/workpersister/infrastructure/postgres"
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

	authorRepository := repositories.NewAuthorRepository()
	workRepository := repositories.NewWorkRepository()
	wordRepository := repositories.NewWordRepository()
	workWordRepository := repositories.NewWorkWordRepository()

	wp := postgres.NewWorkPersister(db, authorRepository, workRepository, wordRepository, workWordRepository)

	api := api.NewAPI(tp, wp)

	http.HandleFunc("/lemmatise", api.Lemmatise())

	fmt.Println("Listening at :6666")

	err = http.ListenAndServe(":6666", nil)
	if err != nil {
		log.Fatal("Failed to start server: " + err.Error())
	}

}
