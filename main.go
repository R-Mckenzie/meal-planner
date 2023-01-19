package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/R-Mckenzie/meal-planner/controllers"
)

func main() {
	staticCtrl := controllers.NewStatic()
	recipeCtrl := controllers.NewRecipes()

	r := chi.NewRouter()

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.Handle("/", staticCtrl.Home)
	r.Get("/recipes/create", recipeCtrl.Create)
	r.Post("/recipes/create", recipeCtrl.Add)

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", r)
	log.Fatal(err)
}
