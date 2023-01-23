package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"

	"github.com/R-Mckenzie/meal-planner/controllers"
	"github.com/R-Mckenzie/meal-planner/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "mealplannerpgadmin"
	dbname   = "mealplanner_dev"
)

func main() {
	infLog := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable TimeZone=UTC", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	infLog.Println("Successfully connected to database")

	userService := models.NewUserService(db, errLog, infLog)

	staticCtrl := controllers.NewStatic()
	recipeCtrl := controllers.NewRecipes(*errLog, *infLog)
	userCtrl := controllers.NewUsers(*userService, *errLog, *infLog)

	r := chi.NewRouter()
	r.Use(secureHeaders) // Set recommended security headers
	r.Use(middleware.Logger)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.Handle("/", staticCtrl.Home)
	r.Get("/recipes/create", recipeCtrl.CreatePage)
	r.Post("/recipes/create", recipeCtrl.Create)
	r.Get("/signup", userCtrl.SignupPage)
	r.Post("/signup", userCtrl.Signup)
	r.Get("/login", userCtrl.LoginPage)
	r.Post("/login", userCtrl.Login)
	r.Get("/cookietest", userCtrl.CookieTest)

	srv := &http.Server{
		Addr:     ":4000",
		ErrorLog: errLog,
		Handler:  r,
	}

	infLog.Println("Starting server on :4000")
	err = srv.ListenAndServe()
	errLog.Fatal(err)
}
