package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"

	_ "github.com/lib/pq"

	"github.com/R-Mckenzie/meal-planner/controllers"
	"github.com/R-Mckenzie/meal-planner/models"
)

type application struct {
	services models.Services
	ctrl     controllers.Controllers
}

func main() {
	iLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	eLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)

	services, err := models.NewServices(iLog, eLog)
	if err != nil {
		panic(err)
	}
	defer services.CloseDB()

	controllers := controllers.NewControllers(*services, eLog, iLog)

	app := &application{
		services: *services,
		ctrl:     *controllers,
	}

	r := chi.NewRouter()
	r.Use(secureHeaders)     // Set recommended security headers
	r.Use(middleware.Logger) // Log all requests made to server
	r.Use(app.authenticate)

	// Routes
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.Get("/", app.ctrl.Static.Home)
	r.Get("/dashboard", app.authorise(app.ctrl.Dashboard.Dashboard))
	r.Post("/dashboard", app.authorise(app.ctrl.Dashboard.SaveMeals))
	r.Get("/recipes", app.authorise(app.ctrl.Recipes.ListPage))
	r.Get("/recipes/create", app.authorise(app.ctrl.Recipes.CreatePage))
	r.Post("/recipes", app.authorise(app.ctrl.Recipes.Create))
	r.Get("/recipes/{recipeID}", app.authorise(app.ctrl.Recipes.UpdatePage))
	r.Post("/recipes/{recipeID}", app.authorise(app.ctrl.Recipes.Update))
	r.Delete("/recipes", app.authorise(app.ctrl.Recipes.Delete))
	r.Get("/signup", app.ctrl.Users.SignupPage)
	r.Post("/signup", app.ctrl.Users.Signup)
	r.Get("/login", app.ctrl.Users.LoginPage)
	r.Post("/login", app.ctrl.Users.Login)
	r.Get("/logout", app.authorise(app.ctrl.Users.Logout))

	r.Get("/cookietest", app.ctrl.Users.CookieTest)

	csrf := nosurf.New(r)
	csrf.ExemptPath("/dashboard")
	csrf.ExemptPath("/recipes")

	srv := &http.Server{
		Addr:    ":4000",
		Handler: csrf,
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
