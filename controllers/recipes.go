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
	returnAddr  string
}

func NewRecipes(rs models.RecipeService, iLog, eLog *log.Logger) *Recipe {
	return &Recipe{
		CreateView: views.NewView("root", "views/recipes/create.html"),
		UpdateView: views.NewView("root", "views/recipes/update.html"),
		ListView:   views.NewView("root", "views/recipes/list.html"),
		rs:         rs,
		iLog:       iLog,
		eLog:       eLog,
	}
}

type Recipe struct {
	CreateView *views.View
	UpdateView *views.View
	ListView   *views.View
	rs         models.RecipeService
	iLog       *log.Logger
	eLog       *log.Logger
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
	userID, ok := r.Context().Value("mealplanner_current_user").(int)
	if !ok {
		http.Error(w, "user not authenticated", http.StatusForbidden)
		return
	}

	formData, err := parseForm(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		re.CreateView.SetAlert("There was a problem adding your recipe", views.Error)
		if err := re.CreateView.Render(w); err != nil {
			panic(err)
		}
		return
	}

	err = re.rs.Create(userID, formData.Title, formData.Method, formData.Ingredients)
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
	switch formData.returnAddr {
	case "dashboard":
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	case "list":
		http.Redirect(w, r, "/recipes", http.StatusSeeOther)
	default:
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

func (re *Recipe) ListPage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("mealplanner_current_user").(int)
	if !ok {
		http.Error(w, "user not authenticated", http.StatusForbidden)
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

	m, t, err := getAlertData(w, r)
	if err != nil {
		panic(err)
	}

	re.ListView.SetAlert(m, t)
	re.ListView.Data.Recipes = recipes
	if err := re.ListView.Render(w); err != nil {
		http.Error(w, "There was an error getting the page", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (re *Recipe) UpdatePage(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("mealplanner_current_user").(int)
	if !ok {
		http.Error(w, "user not authenticated", http.StatusForbidden)
		return
	}
	re.UpdateView.Data.User = userID >= 0
	re.UpdateView.Data.CSRFtoken = nosurf.Token(r)
	re.UpdateView.SetAlert("", views.Success)

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
		return
	}

	re.UpdateView.SetAlert(m, t)
	if err := re.UpdateView.Render(w); err != nil {
		panic(err)
	}
}

func (re *Recipe) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("mealplanner_current_user").(int)
	if !ok || userID < 0 {
		http.Error(w, "user not authenticated", http.StatusForbidden)
		return
	}

	formData, err := parseForm(r)
	if err != nil {
		re.UpdateView.SetAlert("There was a problem adding your recipe", views.Error)
		if err := re.UpdateView.Render(w); err != nil {
			panic(err)
		}
	}
	recipe := re.UpdateView.Data.Recipe
	recipe.Title, recipe.Method, recipe.Ingredients = formData.Title, formData.Method, formData.Ingredients

	err = re.rs.Update(recipe)
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

func parseForm(r *http.Request) (*RecipeForm, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	return &RecipeForm{
		Title:       r.PostForm["title"][0],
		Ingredients: parseIngredients(r.PostForm["ingredients"][0]),
		Method:      r.PostForm["method"][0],
		returnAddr:  r.PostForm["return_url"][0],
	}, nil
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
		http.Error(w, "user not authenticated", http.StatusForbidden)
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
