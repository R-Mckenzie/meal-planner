package models

import (
	"database/sql"
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
	db *sql.DB
}

type RecipeService interface {
	Create(ownerID int, title, method string, ingredients []string) error
	Update(recipe Recipe) error
	Delete(id int) error
	GetTitle(id int) (string, error)
	GetAvailable() ([]Recipe, error)
	GetByUser(ownerID int) ([]Recipe, error)
	GetByID(id int) (*Recipe, error)
}

func NewRecipeService(db *sql.DB) RecipeService {
	return &recipeService{
		db: db,
	}
}

func (rs *recipeService) Create(ownerID int, title, method string, ingredients []string) error {
	_, err := rs.db.Exec(`INSERT INTO recipes(owner_id, title, ingredients, method, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $5)`,
		ownerID, title, pq.Array(ingredients), method, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		return err
	}
	return nil
}

func (rs *recipeService) Update(recipe Recipe) error {
	return nil
}

func (rs *recipeService) Delete(id int) error {
	return nil
}

func (rs *recipeService) GetTitle(id int) (string, error) {
	row := rs.db.QueryRow("SELECT (title) FROM recipes WHERE id=$1", id)
	var title string
	err := row.Scan(&title)
	if err != nil {
		log.Println(err.Error() + "gettitle")
		return "", err
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
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	recipes := make([]Recipe, 0)
	for rows.Next() {
		r := Recipe{}
		err := rows.Scan(&r.ID, &r.OwnerId, &r.Title, &r.Ingredients, &r.Method, &r.createdAt, &r.updatedAt, &r.deletedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		recipes = append(recipes, r)
	}
	return recipes, nil
}

func (rs *recipeService) GetByID(id int) (*Recipe, error) {
	return nil, nil
}
