package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/R-Mckenzie/meal-planner/models"
	"github.com/R-Mckenzie/meal-planner/views"
)

func NewRecipes() *Recipe {
	return &Recipe{
		CreateView: views.NewView("root", "views/recipes/create.html"),
	}
}

type Recipe struct {
	CreateView *views.View
	rs         *models.RecipeService
}

// shows the create recipe page
func (re *Recipe) CreatePage(w http.ResponseWriter, r *http.Request) {
	re.CreateView.Data.User = r.Context().Value("mealplanner_current_user").(bool)
	if err := re.CreateView.Render(w); err != nil {
		panic(err)
	}
}

// Adds a new recipe to the db
func (re *Recipe) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "There was a problem understanding your input", http.StatusBadRequest)
		return
	}

	formData := RecipeForm{
		Title:       r.PostForm["title"][0],
		Ingredients: parseIngredients(r.PostForm["ingredients"][0]),
		Method:      r.PostForm["method"][0],
	}

	log.Println(formData)
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
	Title       string
	Ingredients []string
	Method      string
}
