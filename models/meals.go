package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Meal struct {
	ID        int
	OwnerID   int
	RecipeID  int
	Date      time.Time
	createdAt time.Time
	updatedAt time.Time
	deletedAt sql.NullTime
}

type MealService interface {
	// Methods to get meals
	ByDateRange(id int, start, end time.Time) (*[]Meal, error) // All meals from a user in a date range

	// Methods to create meals
	Create(ownerID, recipeID int, date time.Time) error

	DeleteInRange(ownerID int, start, end time.Time) error
}

type mealService struct {
	db   *sql.DB
	iLog *log.Logger
	eLog *log.Logger
}

func NewMealService(db *sql.DB, iLog, eLog *log.Logger) MealService {
	return &mealService{
		db:   db,
		iLog: iLog,
		eLog: eLog,
	}
}

func (ms *mealService) Create(ownerID int, recipeID int, date time.Time) error {
	meal := &Meal{
		OwnerID:   ownerID,
		RecipeID:  recipeID,
		Date:      date,
		createdAt: time.Now().UTC(),
		updatedAt: time.Now().UTC(),
	}

	err := ms.db.QueryRow("INSERT INTO meals (owner_id, recipe_id, date, created_at, updated_at) VALUES($1, $2, $3, $4, $5) RETURNING id", meal.OwnerID, meal.RecipeID, meal.Date.Format("2006-01-02"), meal.createdAt.Format(time.RFC3339), meal.updatedAt.Format(time.RFC3339)).Scan(&meal.ID)
	if err != nil {
		return fmt.Errorf("in MealService.Create: %w", err)
	}
	ms.iLog.Printf("User %d created new meal on date %s with recipe %q", ownerID, date.Format("2006-01-02"), recipeID)
	return nil
}

func (ms *mealService) DeleteInRange(ownerID int, start, end time.Time) error {
	s := start.Format("2006-01-02")
	e := end.Format("2006-01-02")

	_, err := ms.db.Exec("DELETE FROM meals WHERE date >= $1 AND date <= $2 AND owner_id = $3", s, e, ownerID)
	if err != nil {
		return fmt.Errorf("in MealService.DeleteInRange: %w", err)
	}
	ms.iLog.Printf("User %d deleted meals between %s and %s", ownerID, s, e)
	return nil
}

func (ms *mealService) ByDateRange(ownerID int, start, end time.Time) (*[]Meal, error) {
	s := start.Format("2006-01-02")
	e := end.Format("2006-01-02")
	rows, err := ms.db.Query("SELECT * FROM meals WHERE date >= $1 AND date <= $2 AND owner_id = $3", s, e, ownerID)
	if err != nil {
		return nil, fmt.Errorf("in MealService.ByDateRange query: %w", err)
	}
	defer rows.Close()

	meals := make([]Meal, 0)
	for rows.Next() {
		m := Meal{}
		err := rows.Scan(&m.ID, &m.OwnerID, &m.RecipeID, &m.Date, &m.createdAt, &m.updatedAt, &m.deletedAt)
		if err != nil {
			return nil, fmt.Errorf("in MealService.ByDateRange scanning rows: %w", err)
		}
		meals = append(meals, m)
	}
	ms.iLog.Printf("User %d retrieved meals between %s and %s", ownerID, s, e)
	return &meals, nil
}
