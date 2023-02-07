package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/R-Mckenzie/meal-planner/models"
	"github.com/R-Mckenzie/meal-planner/views"
	"github.com/go-chi/chi/v5"
	"github.com/justinas/nosurf"
)

type RecipeForm struct {
	Title       string
	Ingredients []string
	Method      string
}

func NewRecipes(rs models.RecipeService) *Recipe {
	return &Recipe{
		CreateView: views.NewView("root", "views/recipes/create.html"),
		UpdateView: views.NewView("root", "views/recipes/update.html"),
		ListView:   views.NewView("root", "views/recipes/list.html"),
		rs:         rs,
	}
}

type Recipe struct {
	CreateView *views.View
	UpdateView *views.View
	ListView   *views.View
	rs         models.RecipeService
}

// shows the create recipe page
func (re *Recipe) CreatePage(w http.ResponseWriter, r *http.Request) {
	re.CreateView.Data.User = r.Context().Value("mealplanner_current_user").(int) >= 0
	re.CreateView.Data.CSRFtoken = nosurf.Token(r)
	re.CreateView.Data.Alert.Message = ""

	m, t, err := getAlertData(w, r)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	re.CreateView.SetAlert(m, t)

	if err := re.CreateView.Render(w); err != nil {
		panic(err)
	}
}

// Adds a new recipe to the db
func (re *Recipe) Create(w http.ResponseWriter, r *http.Request) {
	ok := r.Context().Value("mealplanner_current_user").(int) >= 0
	if !ok {
		log.Println("user not authenticated")
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println(err)
		re.CreateView.SetAlert("There was a problem adding your recipe", views.Error)
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
		re.CreateView.SetAlert("There was a problem adding your recipe", views.Error)
		w.WriteHeader(http.StatusInternalServerError)
		if err := re.CreateView.Render(w); err != nil {
			panic(err)
		}
		return
	}

	setAlertData(w, fmt.Sprintf("Successfuly added your %q recipe", formData.Title), views.Success)
	http.Redirect(w, r, "/recipes/create", http.StatusSeeOther)
}

func (re *Recipe) ListPage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("mealplanner_current_user").(int)
	if !ok {
		log.Println("user not authenticated")
		return
	}
	re.ListView.Data.User = userID >= 0
	re.ListView.Data.CSRFtoken = nosurf.Token(r)
	re.ListView.Data.Alert.Message = ""

	recipes, err := re.rs.GetByUser(userID)
	if err != nil {
		log.Println(err)
		return
	}

	re.ListView.Data.Recipes = recipes
	if err := re.ListView.Render(w); err != nil {
		http.Error(w, "There was an error getting the page", http.StatusInternalServerError)
		log.Println(err)
	}
}

func (re *Recipe) UpdatePage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("mealplanner_current_user").(int)
	if !ok {
		log.Println("user not authenticated")
		return
	}
	re.UpdateView.Data.User = userID >= 0
	re.UpdateView.Data.CSRFtoken = nosurf.Token(r)
	re.UpdateView.Data.Alert.Message = ""

	recipeID := chi.URLParam(r, "recipeID")
	if recipeID == "" {
		http.Error(w, "Could not find this recipe", http.StatusNotFound)
		return
	}

	id, err := strconv.Atoi(recipeID)
	if err != nil {
		http.Error(w, "There was a problem...", http.StatusInternalServerError)
		return
	}

	recipe, err := re.rs.GetByID(id, userID)
	if err != nil {
		log.Println(err)
		return
	}
	re.UpdateView.Data.Recipe = *recipe

	m, t, err := getAlertData(w, r)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
	re.UpdateView.SetAlert(m, t)
	if err := re.UpdateView.Render(w); err != nil {
		panic(err)
	}
}

func (re *Recipe) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("mealplanner_current_user").(int)
	if !ok || userID < 0 {
		log.Println("user not authenticated")
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println(err)
		re.UpdateView.SetAlert("There was a problem adding your recipe", views.Error)
		w.WriteHeader(http.StatusBadRequest)
		if err := re.UpdateView.Render(w); err != nil {
			panic(err)
		}
		return
	}

	formData := RecipeForm{
		Title:       r.PostForm["title"][0],
		Ingredients: parseIngredients(r.PostForm["ingredients"][0]),
		Method:      r.PostForm["method"][0],
	}

	recipe := re.UpdateView.Data.Recipe
	recipe.Title = formData.Title
	recipe.Method = formData.Method
	recipe.Ingredients = formData.Ingredients

	err := re.rs.Update(recipe)
	if err != nil {
		log.Println(err)
		re.UpdateView.SetAlert("There was a problem adding your recipe", views.Error)
		w.WriteHeader(http.StatusInternalServerError)
		if err := re.UpdateView.Render(w); err != nil {
			panic(err)
		}
		return
	}

	re.UpdateView.Data.Recipe = recipe

	setAlertData(w, fmt.Sprintf("Successfuly updated your %q recipe", formData.Title), views.Success)
	rURL := fmt.Sprintf("/recipes/%d", recipe.ID)
	http.Redirect(w, r, rURL, http.StatusSeeOther)
}

func parseIngredients(text string) []string {
	split := strings.Split(text, "\n")
	var cleaned []string
	for _, i := range split {
		cleaned = append(cleaned, strings.TrimSpace(i))
	}
	return cleaned
}

func (re *Recipe) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("mealplanner_current_user").(int)
	if userID < 0 {
		log.Println("user not authenticated")
		return
	}

	type body struct {
		ID   int    `json:"recipeID"`
		CSRF string `json:"csrf"`
	}

	b := body{}
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		log.Println(err)
		return
	}

	ok = nosurf.VerifyToken(nosurf.Token(r), b.CSRF)
	if !ok {
		log.Println("csrf fail")
	}

	re.rs.Delete(b.ID, userID)
}
