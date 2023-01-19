package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/R-Mckenzie/meal-planner/views"
)

func NewRecipes() *Recipe {
	return &Recipe {
		CreateView: views.NewView("root", "views/recipes/create.html"),
	}
}

type Recipe struct {
	CreateView *views.View	
}

func (re *Recipe) Create(w http.ResponseWriter, r *http.Request) {
	if err := re.CreateView.Render(w, nil); err != nil {
		panic(err)
	}
}

func (re *Recipe) Add(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	formData := RecipeForm {
		Title: r.PostForm["title"][0],
		Ingredients: parseIngredients(r.PostForm["ingredients"][0]),
		Method: r.PostForm["method"][0],
	}

	fmt.Fprintln(w, formData.Title)
	fmt.Fprintln(w, formData.Ingredients)
	fmt.Fprintln(w, formData.Method)
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
