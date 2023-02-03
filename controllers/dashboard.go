package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/R-Mckenzie/meal-planner/models"
	"github.com/R-Mckenzie/meal-planner/views"
	"github.com/justinas/nosurf"
)

type Dashboard struct {
	DashView *views.View
	rs       models.RecipeService
	ms       models.MealService
	us       models.UserService
	week     time.Time
}

func NewDashboard(rs models.RecipeService, ms models.MealService, us models.UserService) *Dashboard {
	return &Dashboard{
		DashView: views.NewView("root", "views/dashboard/dashboard.html"),
		rs:       rs,
		ms:       ms,
		us:       us,
		week:     time.Now(),
	}
}

func shiftWeekStartMonday(day time.Weekday) int {
	adjusted := day - 1
	if adjusted == -1 {
		adjusted = 6
	}
	return int(adjusted)
}

func getMonday(date time.Time) time.Time {
	monday := date.AddDate(0, 0, -shiftWeekStartMonday(date.Weekday()))
	return monday
}

func weekBoundaries(date time.Time) (time.Time, time.Time) {
	start := getMonday(date)
	end := start.AddDate(0, 0, 6)
	return start, end
}

func (d *Dashboard) Dashboard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("mealplanner_current_user").(int)
	d.DashView.Data.User = userID >= 0
	d.DashView.Data.CSRFtoken = nosurf.Token(r)

	recipes, err := d.rs.GetByUser(userID)
	if err != nil {
		log.Println(err)
		return
	}

	monday, sunday := weekBoundaries(time.Now())
	meals, err := d.ms.ByDateRange(userID, monday, sunday)
	if err != nil {
		log.Println(err)
		return
	}

	mealMap := make(map[int][]views.MealNode)
	log.Printf("mealmap: %+v", mealMap)

	for _, m := range *meals {
		t, err := d.rs.GetTitle(m.RecipeID)
		if err != nil {
			log.Println(err)
			return
		}
		node := views.MealNode{
			Title:    t,
			RecipeID: m.RecipeID,
		}

		mealMap[shiftWeekStartMonday(m.Date.Weekday())] = append(mealMap[shiftWeekStartMonday(m.Date.Weekday())], node)
	}

	d.DashView.Data.MealMap = mealMap
	d.DashView.Data.Recipes = recipes
	if err := d.DashView.Render(w); err != nil {
		panic(err)
	}
}

func (d *Dashboard) SaveMeals(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("mealplanner_current_user").(int)
	if !ok {
		log.Println("user not authenticated")
		return
	}

	type reqBody struct {
		CSRF  string `json:"csrf"`
		Meals []struct {
			RecipeID int       `json:"recipeID"`
			Date     time.Time `json:"date"`
		} `json:"meals"`
	}

	data := reqBody{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
	}

	ok = nosurf.VerifyToken(nosurf.Token(r), data.CSRF)
	if !ok {
		log.Println("csrf fail")
	}

	monday, sunday := weekBoundaries(time.Now())
	d.ms.DeleteInRange(userID, monday, sunday)

	log.Println(data)
	for _, m := range data.Meals {
		err := d.ms.Create(userID, m.RecipeID, m.Date)
		if err != nil {
			log.Println(err)
			return
		}
	}

	//success message
	log.Println("success")
}
