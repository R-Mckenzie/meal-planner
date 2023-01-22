package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/R-Mckenzie/meal-planner/models"
	"github.com/R-Mckenzie/meal-planner/views"
)

func NewRecipes(errLog, infLog log.Logger) *Recipe {
	return &Recipe {
		CreateView: views.NewView("root", "views/recipes/create.html"),
		errLog: &errLog,
		infLog: &infLog,
	}
}

type Recipe struct {
	CreateView *views.View	
	rs *models.RecipeService
	errLog *log.Logger
	infLog *log.Logger
}

// shows the create recipe page
func (re *Recipe) CreatePage(w http.ResponseWriter, r *http.Request) {
	if err := re.CreateView.Render(w, nil); err != nil {
		panic(err)
	}
}

// Adds a new recipe to the db
func (re *Recipe) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		re.errLog.Println("Error parsing recipe form: ", err)
	}

	formData := RecipeForm {
		Title: r.PostForm["title"][0],
		Ingredients: parseIngredients(r.PostForm["ingredients"][0]),
		Method: r.PostForm["method"][0],
	}

	fmt.Fprintln(w, formData.Title)
	fmt.Fprintln(w, formData.Ingredients)
	fmt.Fprintln(w, formData.Method)
	re.infLog.Println(formData)
}

func parseIngredients(text string) []string {
	split := strings.Split(text, "\n")
	var cleaned []string
	for _, i := range split {
		cleaned = append(cleaned, strings.TrimSpace(i))
	}
	return cleaned
}

type RecipeForm struct {
	Title string
	Ingredients []string
	Method string
}
