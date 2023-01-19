package controllers

import (
	"github.com/R-Mckenzie/meal-planner/views"
)

func NewStatic() *Static {
	return &Static{
		Home: views.NewView("root", "views/static/home.html"),
	}
}

type Static struct {
	Home *views.View
}
