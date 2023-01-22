package models

import (
	"database/sql"
	"time"
)

type Recipe struct {
	id          uint
	ownerId     uint
	title       string
	ingredients []string
	method      string
	createdAt   time.Time
	updatedAt   time.Time
	deletedAt   time.Time
}

type RecipeService struct {
	db *sql.DB
}

func (rs *RecipeService) Create(ownerId uint, title, method string, ingredients []string) error {
	_, err := rs.db.Exec(`INSERT INTO recipes(owner_id, title, ingredients, method, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5, $5)`,
		ownerId, title, ingredients, method, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return err
	}
	return nil
}

func (rs *RecipeService) Update(id, ownerId uint, title, method string, ingredients []string) error {
	return nil
}

func (rs *RecipeService) Delete(id uint) error {
	return nil
}

// Retrieves both the user's own recipes and the public recipes (from user 0)
func (rs *RecipeService) GetAvailable() error {
	return nil
}

func (rs *RecipeService) GetByUser(ownerId uint) error {
	return nil
}

func (rs *RecipeService) GetById(id uint) error {
	return nil
}
