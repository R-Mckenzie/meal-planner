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
	iLog      *log.Logger
	eLog      *log.Logger
}

func NewControllers(s models.Services, eLog, iLog *log.Logger) *Controllers {
	return &Controllers{
		Static:    *NewStatic(),
		Users:     *NewUsers(s.Users),
		Recipes:   *NewRecipes(s.Recipes),
		Dashboard: *NewDashboard(s.Recipes, s.Meals, s.Users),
		iLog:      iLog,
		eLog:      eLog,
	}
}
