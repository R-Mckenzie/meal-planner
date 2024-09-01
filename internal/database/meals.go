package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/jackc/pgx/v5"
)

type MealTime int

const (
	Breakfast MealTime = iota
	Lunch
	Dinner
)

type Meal struct {
	ID       int
	OwnerID  int
	RecipeID int
	Date     time.Time
	MealTime MealTime

	CreatedAt time.Time
	UpdatedAt time.Time
}

type MealRecipe struct {
	MealID      int
	MealDate    time.Time
	MealTime    MealTime
	RecipeID    int
	RecipeTitle string
}

type MealStore interface {
	Create(ownerID, recipeID int, date time.Time, mealTime MealTime) (int, error)
	GetByID(id int) (Meal, error)
	GetByDateRange(start, end time.Time) ([]Meal, error)
	GetByDateRangeWithRecipes(start, end time.Time) ([]MealRecipe, error)
	UpdateMealDate(id int, date time.Time, mealTime MealTime) error
	DeleteByID(userID, mealID int) error
}

type postgresMealStore struct {
	db *sql.DB
}

func NewMealStore(db *sql.DB) MealStore {
	return &postgresMealStore{db}
}

func (ms *postgresMealStore) Create(ownerID, recipeID int, date time.Time, mealTime MealTime) (int, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	var id int
	err := ms.db.QueryRow("INSERT INTO meals (owner_id, recipe_id, date, meal_time, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		ownerID, recipeID, date.UTC().Format(time.RFC3339), mealTime, now, now).Scan(&id)

	if err != nil {
		slog.Error("Error creating new meal", "error", err)
		return -1, errors.New("There was a problem creating this meal")
	}

	return id, nil
}

func (ms *postgresMealStore) GetByID(id int) (Meal, error) {
	row := ms.db.QueryRow("SELECT id, owner_id, recipe_id, date, meal_time, created_at, updated_at FROM meals where id=$1", id)
	var m = &Meal{}
	err := row.Scan(&m.ID, &m.OwnerID, &m.RecipeID, &m.Date, &m.MealTime, &m.CreatedAt, &m.UpdatedAt)
	if err != nil {
		slog.Error("Error getting meal by ID", "error", row.Err(), "mealID", id)
		return *m, fmt.Errorf("in GetByID: %w", err)
	}
	return *m, nil
}

func (ms *postgresMealStore) GetByDateRange(start, end time.Time) ([]Meal, error) {
	rows, err := ms.db.Query("SELECT id, owner_id, recipe_id, date, meal_time, created_at, updated_at FROM meals where date >= $1 AND date <= $2", start, end)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("in GetByDateRange querying database: %w", err)
	}

	meals := []Meal{}
	for rows.Next() {
		var m = Meal{}
		err := rows.Scan(&m.ID, &m.OwnerID, &m.RecipeID, &m.Date, &m.MealTime, &m.CreatedAt, &m.UpdatedAt)
		if err != nil {
			slog.Error("Error getting meal by date range", "error", rows.Err(), "startDate", start, "endDate", end)
			return nil, fmt.Errorf("in GetByDateRange scanning row: %w", err)
		}
		meals = append(meals, m)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("in GetByDateRange: %w", err)
	}

	return meals, nil
}

func (ms *postgresMealStore) GetByDateRangeWithRecipes(start, end time.Time) ([]MealRecipe, error) {
	rows, err := ms.db.Query("SELECT m.id AS meal_id, m.date, m.meal_time, m.recipe_id, r.title AS recipe_name FROM meals m INNER JOIN recipes r ON m.recipe_id = r.id where date >= $1 AND date <= $2", start, end)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("in GetByDateRange querying database: %w", err)
	}

	meals := []MealRecipe{}
	for rows.Next() {
		var m = MealRecipe{}
		err := rows.Scan(&m.MealID, &m.MealDate, &m.MealTime, &m.RecipeID, &m.RecipeTitle)
		if err != nil {
			slog.Error("Error getting meal by date range", "error", rows.Err(), "startDate", start, "endDate", end)
			return nil, fmt.Errorf("in GetByDateRange scanning row: %w", err)
		}
		meals = append(meals, m)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("in GetByDateRange: %w", err)
	}

	return meals, nil
}

func (ms *postgresMealStore) UpdateMealDate(id int, date time.Time, mealTime MealTime) error {
	now := time.Now().UTC().Format(time.RFC3339)
	row := ms.db.QueryRow("UPDATE meals SET date=$1, meal_time=$2, updated_at=$3 WHERE id=$4", date.UTC().Format(time.RFC3339), mealTime, now, id)
	if row.Err() != nil {
		return row.Err()
	}
	return nil
}

func (ms *postgresMealStore) DeleteByID(ownerID, mealId int) error {
	_, err := ms.db.Exec("DELETE FROM meals WHERE owner_id=$1 AND id=$2", ownerID, mealId)
	if err != nil {
		return err
	}
	return nil
}
