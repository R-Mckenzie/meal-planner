package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"

	"github.com/R-Mckenzie/meal-planner/controllers"
	"github.com/R-Mckenzie/meal-planner/models"
)

type application struct {
	services models.Services
}

func main() {
	// Create servies
	services, err := models.NewServices()
	if err != nil {
		panic(err)
	}
	defer services.CloseDB()

	app := &application{
		services: *services,
	}

	// Create controllers
	staticCtrl := controllers.NewStatic()
	recipeCtrl := controllers.NewRecipes()
	userCtrl := controllers.NewUsers(services.Users)

	r := chi.NewRouter()
	r.Use(secureHeaders)     // Set recommended security headers
	r.Use(middleware.Logger) // Log all requests made to server
	r.Use(app.authenticate)

	// Routes
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.Get("/", staticCtrl.Home)
	r.Get("/recipes/create", app.authorise(recipeCtrl.CreatePage))
	r.Post("/recipes/create", app.authorise(recipeCtrl.Create))
	r.Get("/signup", userCtrl.SignupPage)
	r.Post("/signup", userCtrl.Signup)
	r.Get("/login", userCtrl.LoginPage)
	r.Post("/login", userCtrl.Login)
	r.Get("/logout", userCtrl.Logout)

	r.Get("/cookietest", userCtrl.CookieTest)

	srv := &http.Server{
		Addr:    ":4000",
		Handler: r,
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
