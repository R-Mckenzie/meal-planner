package controllers

import (
	"net/http"

	"github.com/R-Mckenzie/meal-planner/views"
)

func NewStatic() *Static {
	return &Static{
		HomeView: views.NewView("root", "views/static/home.html"),
	}
}

func (s *Static) Home(w http.ResponseWriter, r *http.Request) {
	s.HomeView.Data.User = r.Context().Value("mealplanner_current_user").(bool)
	if err := s.HomeView.Render(w); err != nil {
		panic(err)
	}
}

type Static struct {
	HomeView *views.View
}
