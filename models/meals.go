package models

import (
	"database/sql"
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
	ByID(id int) (*Meal, error)                                // Get a meal by it's id
	ByOwner(id int) (*[]Meal, error)                           // All meals from a given user
	ByDate(id int, date time.Time) (*[]Meal, error)            // All meals from a user on a given date
	ByDateRange(id int, start, end time.Time) (*[]Meal, error) // All meals from a user in a date range

	// Methods to create meals
	Create(ownerID, recipeID int, date time.Time) error

	DeleteInRange(ownerID int, start, end time.Time) error
}

type mealService struct {
	db *sql.DB
}

func NewMealService(db *sql.DB) MealService {
	return &mealService{
		db: db,
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
		log.Println(err)
		return err
	}
	return nil
}

func (ms *mealService) ByDate(id int, date time.Time) (*[]Meal, error) {
	return nil, nil
}

func (ms *mealService) DeleteInRange(ownerID int, start, end time.Time) error {
	s := start.Format("2006-01-02")
	e := end.Format("2006-01-02")

	log.Printf("Start: %q, end: %q\n", s, e)

	_, err := ms.db.Exec("DELETE FROM meals WHERE date >= $1 AND date <= $2 AND owner_id = $3", s, e, ownerID)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (ms *mealService) ByDateRange(id int, start, end time.Time) (*[]Meal, error) {
	s := start.Format("2006-01-02")
	e := end.Format("2006-01-02")
	rows, err := ms.db.Query("SELECT * FROM meals WHERE date >= $1 AND date <= $2 AND owner_id = $3", s, e, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	meals := make([]Meal, 0)
	for rows.Next() {
		m := Meal{}
		err := rows.Scan(&m.ID, &m.OwnerID, &m.RecipeID, &m.Date, &m.createdAt, &m.updatedAt, &m.deletedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		meals = append(meals, m)
	}
	return &meals, nil
}

func (ms *mealService) ByOwner(id int) (*[]Meal, error) {
	return nil, nil
}

func (ms *mealService) ByID(id int) (*Meal, error) {
	return nil, nil
}
