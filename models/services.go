package models

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	host     = "localhost"
	port     = 5433
	user     = "postgres"
	password = "mealplannerpgadmin"
	dbname   = "mealplanner_db"
)

type Services struct {
	db      *sql.DB
	Users   UserService
	Recipes RecipeService
	Meals   MealService
	iLog    *log.Logger
	eLog    *log.Logger
}

func NewServices(iLog, eLog *log.Logger) (*Services, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable TimeZone=UTC", host, port, user, password, dbname)
	fmt.Println(psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println("Pinging DB...")
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("Successfully connected to database")

	rs := NewRecipeService(db, iLog, eLog)
	ms := NewMealService(db, iLog, eLog)
	us, err := NewUserService(db, iLog, eLog)
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
