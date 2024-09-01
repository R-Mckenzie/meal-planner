package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/lib/pq"
)

type Recipe struct {
	ID          int
	OwnerID     int
	Title       string
	Ingredients pq.StringArray
	Method      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RecipeStore interface {
	Create(ownerID int, title, method string, ingredients []string) (int, error)
	GetByUser(ownerID int) ([]Recipe, error)
	GetByID(recipeID int) (Recipe, error)
	Update(recipe Recipe) error
	Delete(ownerID, recipeID int) error
}

func NewRecipeStore(db *sql.DB) RecipeStore {
	return &postgresRecipeStore{db}
}

type postgresRecipeStore struct {
	db *sql.DB
}

func (rs *postgresRecipeStore) Create(ownerID int, title, method string, ingredients []string) (int, error) {
	recipe := &Recipe{
		OwnerID:     ownerID,
		Title:       title,
		Method:      method,
		Ingredients: ingredients,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	err := rs.db.QueryRow("INSERT INTO recipes (owner_id, title, ingredients, method, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6) RETURNING id",
		recipe.OwnerID, recipe.Title, recipe.Ingredients, recipe.Method, recipe.CreatedAt.Format(time.RFC3339), recipe.UpdatedAt.Format(time.RFC3339)).Scan(&recipe.ID)

	if err != nil {
		slog.Error("Error creating recipe", "error", err)
		return -1, fmt.Errorf("There was a problem signing you up")
	}

	fmt.Println(recipe)
	return recipe.ID, nil
}

func (rs *postgresRecipeStore) GetByUser(ownerID int) ([]Recipe, error) {
	rows, err := rs.db.Query("SELECT * from recipes WHERE owner_id=$1 ORDER BY id DESC", ownerID)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("in GetByUser querying database: %w", err)
	}
	recipes := []Recipe{}
	for rows.Next() {
		var r = Recipe{}
		err := rows.Scan(&r.ID, &r.OwnerID, &r.Title, &r.Ingredients, &r.Method, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			slog.Error("Error getting recipes by user", "error", rows.Err(), "ownerID", ownerID)
			return nil, fmt.Errorf("in GetByuser scanning row: %w", err)
		}
		recipes = append(recipes, r)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("in GetByUser: %w", err)
	}
	return recipes, nil
}

func (rs *postgresRecipeStore) GetByID(recipeID int) (Recipe, error) {
	row := rs.db.QueryRow("SELECT * from recipes WHERE id=$1", recipeID)
	r := Recipe{}
	err := row.Scan(&r.ID, &r.OwnerID, &r.Title, &r.Ingredients, &r.Method, &r.CreatedAt, &r.UpdatedAt)

	if err != nil {
		slog.Error("Error getting recipes by ID", "error", row.Err(), "recipeID", recipeID)
		return Recipe{}, fmt.Errorf("in GetByID scanning row: %w", err)
	}
	return r, nil
}

func (rs *postgresRecipeStore) Update(recipe Recipe) error {
	now := time.Now().UTC().Format(time.RFC3339)
	row := rs.db.QueryRow("UPDATE recipes SET title=$1, ingredients=$2, method=$3, updated_at=$4 WHERE id=$5", recipe.Title, recipe.Ingredients, recipe.Method, now, recipe.ID)

	if row.Err() != nil {
		slog.Error("Error updating recipe", "error", row.Err(), "recipe", recipe)
		return row.Err()
	}

	return nil
}

func (rs *postgresRecipeStore) Delete(ownerID, recipeID int) error {
	_, err := rs.db.Exec("DELETE FROM recipes WHERE owner_id=$1 AND id=$2", ownerID, recipeID)
	if err != nil {
		return err
	}
	return nil
}
