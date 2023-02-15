package models

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"
)

type Recipe struct {
	ID          int
	OwnerId     int
	Title       string
	Ingredients pq.StringArray
	Method      string
	createdAt   time.Time
	updatedAt   time.Time
	deletedAt   sql.NullTime
}

type recipeService struct {
	db   *sql.DB
	iLog *log.Logger
	eLog *log.Logger
}

type RecipeService interface {
	Create(ownerID int, title, method string, ingredients []string) error
	Update(recipe Recipe) error
	Delete(recipeID, ownderID int) error
	GetTitle(id int) (string, error)
	GetAvailable() ([]Recipe, error)
	GetByUser(ownerID int) ([]Recipe, error)
	GetByID(id, ownerID int) (*Recipe, error)
}

func NewRecipeService(db *sql.DB, iLog, eLog *log.Logger) RecipeService {
	return &recipeService{
		db:   db,
		iLog: iLog,
		eLog: eLog,
	}
}

func (rs *recipeService) Create(ownerID int, title, method string, ingredients []string) error {
	_, err := rs.db.Exec(`INSERT INTO recipes(owner_id, title, ingredients, method, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $5)`,
		ownerID, title, pq.Array(ingredients), method, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return fmt.Errorf("in RecipeService.Create: %w", err)
	}
	return nil
}

func (rs *recipeService) Update(recipe Recipe) error {
	_, err := rs.db.Exec("UPDATE recipes SET title=$2, ingredients=$3, method=$4, updated_at=$5 WHERE id=$1",
		recipe.ID, recipe.Title, recipe.Ingredients, recipe.Method, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		rs.eLog.Printf("in recipeService.Update: %v", err)
		return fmt.Errorf("in RecipeService.Update: %w", err)
	}
	return nil
}

func (rs *recipeService) Delete(recipeID, ownerID int) error {
	_, err := rs.db.Exec("DELETE FROM recipes WHERE id=$1 AND owner_id=$2", recipeID, ownerID)
	if err != nil {
		return fmt.Errorf("in RecipeService.Delete: %w", err)
	}
	return nil
}

func (rs *recipeService) GetTitle(id int) (string, error) {
	row := rs.db.QueryRow("SELECT (title) FROM recipes WHERE id=$1", id)
	var title string
	err := row.Scan(&title)
	if err != nil {
		return "", fmt.Errorf("in RecipeService.GetTitle: %w", err)
	}
	return title, nil
}

// Retrieves both the user's own recipes and the public recipes (from user 0)
func (rs *recipeService) GetAvailable() ([]Recipe, error) {
	return nil, nil
}

func (rs *recipeService) GetByUser(ownerId int) ([]Recipe, error) {
	rows, err := rs.db.Query("SELECT * FROM recipes WHERE owner_id=$1", ownerId)
	if err != nil {
		return nil, fmt.Errorf("in ReciperService.GetByUser: %w", err)
	}
	defer rows.Close()

	recipes := make([]Recipe, 0)
	for rows.Next() {
		r := Recipe{}
		err := rows.Scan(&r.ID, &r.OwnerId, &r.Title, &r.Ingredients, &r.Method, &r.createdAt, &r.updatedAt, &r.deletedAt)
		if err != nil {
			return nil, fmt.Errorf("in ReciperService.GetByUser: %w", err)
		}
		recipes = append(recipes, r)
	}
	return recipes, nil
}

func (rs *recipeService) GetByID(id, ownerID int) (*Recipe, error) {
	r := &Recipe{}
	row := rs.db.QueryRow("SELECT * FROM recipes WHERE id=$1 AND owner_id=$2", id, ownerID)
	err := row.Scan(&r.ID, &r.OwnerId, &r.Title, &r.Ingredients, &r.Method, &r.createdAt, &r.updatedAt, &r.deletedAt)
	if err != nil {
		return nil, fmt.Errorf("in RecipeService.GetByID: %w", err)
	}
	return r, nil
}
