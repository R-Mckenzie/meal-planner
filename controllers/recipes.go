package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/R-Mckenzie/meal-planner/models"
	"github.com/R-Mckenzie/meal-planner/views"
	"github.com/justinas/nosurf"
)

func NewRecipes(rs models.RecipeService) *Recipe {
	return &Recipe{
		CreateView: views.NewView("root", "views/recipes/create.html"),
		rs:         rs,
	}
}

type Recipe struct {
	CreateView *views.View
	rs         models.RecipeService
}

// shows the create recipe page
func (re *Recipe) CreatePage(w http.ResponseWriter, r *http.Request) {
	re.CreateView.Data.User = r.Context().Value("mealplanner_current_user").(int) >= 0
	re.CreateView.Data.CSRFtoken = nosurf.Token(r)
	re.CreateView.Data.Alert.Message = ""
	if err := re.CreateView.Render(w); err != nil {
		panic(err)
	}
}

// Adds a new recipe to the db
func (re *Recipe) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		re.CreateView.Data.Alert.Message = "There was a problem adding your recipe"
		re.CreateView.Data.Alert.Type = views.Error
		w.WriteHeader(http.StatusBadRequest)
		if err := re.CreateView.Render(w); err != nil {
			panic(err)
		}
		return
	}

	formData := RecipeForm{
		Title:       r.PostForm["title"][0],
		Ingredients: parseIngredients(r.PostForm["ingredients"][0]),
		Method:      r.PostForm["method"][0],
	}

	userID := r.Context().Value("mealplanner_current_user").(int)
	err := re.rs.Create(userID, formData.Title, formData.Method, formData.Ingredients)
	if err != nil {
		log.Println(err)
		re.CreateView.Data.Alert.Message = "There was a problem adding your recipe"
		re.CreateView.Data.Alert.Type = views.Error
		w.WriteHeader(http.StatusInternalServerError)
		if err := re.CreateView.Render(w); err != nil {
			panic(err)
		}
		return
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
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
