package services

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/R-Mckenzie/mealplanner/cmd/web/pages"
	"github.com/R-Mckenzie/mealplanner/internal/auth"
	"github.com/R-Mckenzie/mealplanner/internal/database"
	"github.com/gorilla/csrf"
)

type MealService struct {
	mealStore database.MealStore
	auth      *auth.AuthService
}

func NewMealsService(mealStore database.MealStore, auth *auth.AuthService) *MealService {
	return &MealService{
		mealStore: mealStore,
		auth:      auth,
	}
}

func (ms *MealService) PostMealsHandler(w http.ResponseWriter, r *http.Request) {
	dateString := r.FormValue("date")
	mealIDStrings := r.Form["meal_id"]
	recipeNames := r.Form["recipe_name"]
	recipeIDStrings := r.Form["recipe_id"]
	var recipeIDs []int
	for _, s := range recipeIDStrings {
		id, err := strconv.Atoi(s)
		if err != nil {
			slog.Error("Error parsing recipeID to int", "recipeID", s, "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		recipeIDs = append(recipeIDs, id)
	}
	var mealIDs []int
	for _, s := range mealIDStrings {
		id, err := strconv.Atoi(s)
		if err != nil {
			slog.Error("Error parsing mealID to int", "mealID", s, "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		mealIDs = append(mealIDs, id)
	}
	user := ms.auth.AuthorisedUser(r.Context())
	mealDate, err := time.Parse(time.DateOnly, dateString)
	if err != nil {
		slog.Error("Error parsing time", "time_string", dateString, "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	data, err := ms.CreateOrUpdateMeals(user, mealIDs, recipeNames, recipeIDs, mealDate)
	if err != nil {
		slog.Error("CreateOrUpdateMeals", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	slog.Info("data:", "data", data)

	date, err := time.Parse(time.DateOnly, dateString)
	pages.MealPlanItems(data, date, csrf.Token(r)).Render(r.Context(), w)
}

func (ms *MealService) DeleteMealsHandler(w http.ResponseWriter, r *http.Request) {
	csrf := r.URL.Query().Get("csrf_token")
	slog.Warn("not using CSRF", "csrf_token", csrf)

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ownerID := ms.auth.AuthorisedUser(r.Context())
	err = ms.mealStore.DeleteByID(ownerID, id)
	if err != nil {
		slog.Error("DeleteMealsHandler", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (ms *MealService) GetMealsByDateRange(start, end time.Time) ([]database.MealRecipe, error) {
	meals, err := ms.mealStore.GetByDateRangeWithRecipes(start, end)
	if err != nil {
		return nil, err
	}
	return meals, nil
}

func (ms *MealService) CreateOrUpdateMeals(userID int, mealIDs []int, recipes []string, recipeIDs []int, mealDate time.Time) ([]pages.MealData, error) {
	data := []pages.MealData{}
	for i, mealID := range mealIDs {
		if mealID == -1 {
			id, err := ms.mealStore.Create(userID, recipeIDs[i], mealDate, database.Breakfast)
			if err != nil {
				slog.Error(err.Error())
				return nil, err
			}
			data = append(data, pages.MealData{RecipeName: recipes[i], RecipeID: recipeIDs[i], MealID: id})
		} else {
			err := ms.mealStore.UpdateMealDate(mealID, mealDate, database.Breakfast)
			if err != nil {
				slog.Error(err.Error())
				return nil, err
			}
			data = append(data, pages.MealData{RecipeName: recipes[i], RecipeID: recipeIDs[i], MealID: mealIDs[i]})
		}
	}
	return data, nil
}
