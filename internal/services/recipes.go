package services

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/R-Mckenzie/mealplanner/cmd/web/pages"
	"github.com/R-Mckenzie/mealplanner/internal/auth"
	"github.com/R-Mckenzie/mealplanner/internal/database"
	"github.com/gorilla/csrf"
)

type RecipeService struct {
	recipeStore database.RecipeStore
	auth        *auth.AuthService
}

func NewRecipeService(recipeStore database.RecipeStore, auth *auth.AuthService) *RecipeService {
	return &RecipeService{
		recipeStore: recipeStore,
		auth:        auth,
	}
}

func (rs *RecipeService) CreateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	owner := rs.auth.AuthorisedUser(r.Context())
	err := r.ParseForm()
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	title := r.Form.Get("title")
	ingredients := r.Form.Get("ingredients")
	method := r.Form.Get("method")

	id, err := rs.recipeStore.Create(owner, title, method, strings.Split(ingredients, "\n"))
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	pages.RecipesCard(title, id).Render(r.Context(), w)
}

func (rs *RecipeService) BuildUserRecipeData(user int, ctx context.Context) []pages.RecipesData {
	recipes, err := rs.ByUser(rs.auth.AuthorisedUser(ctx))
	if err != nil {
		slog.Error(err.Error())
	}

	recipeData := []pages.RecipesData{}
	for _, r := range recipes {
		recipeData = append(recipeData, pages.RecipesData{RecipeTitle: r.Title, RecipeID: r.ID})
	}

	return recipeData
}

func (rs *RecipeService) UpdateRecipeHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = r.ParseForm()
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	title := r.Form.Get("title")
	ingredients := r.Form.Get("ingredients")
	method := r.Form.Get("method")

	recipe := database.Recipe{
		ID:          id,
		Title:       title,
		Ingredients: strings.Split(ingredients, "\n"),
		Method:      method,
	}

	err = rs.recipeStore.Update(recipe)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pages.RecipesCard(title, id).Render(r.Context(), w)

}

func (rs *RecipeService) DeleteRecipeHandler(w http.ResponseWriter, r *http.Request) {
	recipeID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	owner := rs.auth.AuthorisedUser(r.Context())
	err = rs.recipeStore.Delete(owner, recipeID)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pages.CreateRecipeForm(true, csrf.Token(r)).Render(r.Context(), w)

}

func (rs *RecipeService) GetRecipeHandler(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if id == -1 {
		pages.CreateRecipeForm(false, csrf.Token(r)).Render(r.Context(), w)
		return
	}

	recipe, err := rs.ByID(id)
	pages.RecipeForm(pages.RecipesData{RecipeID: recipe.ID, RecipeTitle: recipe.Title, RecipeIngredients: recipe.Ingredients, RecipeMethod: recipe.Method}, csrf.Token(r)).Render(r.Context(), w)
}

func (rs *RecipeService) ByUser(ownerID int) ([]database.Recipe, error) {
	recipes, err := rs.recipeStore.GetByUser(ownerID)
	if err != nil {
		return nil, fmt.Errorf("There was a problem getting user recipes: %w", err)
	}
	return recipes, nil
}

func (rs *RecipeService) ByID(recipeID int) (database.Recipe, error) {
	recipe, err := rs.recipeStore.GetByID(recipeID)
	if err != nil {
		return database.Recipe{}, fmt.Errorf("There was a problem getting user recipes: %w", err)
	}
	return recipe, nil
}
