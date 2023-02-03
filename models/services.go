package models

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "mealplannerpgadmin"
	dbname   = "mealplanner_dev"
)

type Services struct {
	db      *sql.DB
	Users   UserService
	Recipes RecipeService
	Meals   MealService
}

func NewServices() (*Services, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable TimeZone=UTC", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("Successfully connected to database")

	rs := NewRecipeService(db)
	ms := NewMealService(db)
	us, err := NewUserService(db)
	if err != nil {
		return nil, err
	}

	return &Services{
		Users:   us,
		Recipes: rs,
		Meals:   ms,
	}, nil
}

func (s *Services) CloseDB() {
	s.db.Close()
}
