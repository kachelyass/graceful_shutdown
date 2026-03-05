package main

import (
	"log"
	"net/http"

	"example.com/m/v2/internal/database"
	"example.com/m/v2/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func main() {
	db, err := database.Connect("postgresql://user:pass@localhost:5432/db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := database.NewTaskStore(db)
	taskHandler := handlers.NewTaskHandler(store)

	r := chi.NewRouter()
	r.Mount("/tasks", taskHandler.Routes())
	http.ListenAndServe(":8080", r)
}
