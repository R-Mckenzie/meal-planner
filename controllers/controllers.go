package controllers

import (
	"log"

	"github.com/R-Mckenzie/meal-planner/models"
)

type Controllers struct {
	Static    Static
	Users     User
	Recipes   Recipe
	Dashboard Dashboard
}

func NewControllers(s models.Services, eLog, iLog *log.Logger) *Controllers {
	return &Controllers{
		Static:    *NewStatic(),
		Users:     *NewUsers(s.Users, iLog, eLog),
		Recipes:   *NewRecipes(s.Recipes, iLog, eLog),
		Dashboard: *NewDashboard(s.Recipes, s.Meals, s.Users, iLog, eLog),
	}
}
