package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"

	"github.com/R-Mckenzie/mealplanner/internal/auth"
	"github.com/R-Mckenzie/mealplanner/internal/database"
	"github.com/R-Mckenzie/mealplanner/internal/services"
)

type Server struct {
	port int

	db *sql.DB

	// Services
	auth    *auth.AuthService
	meals   *services.MealService
	recipes *services.RecipeService
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	host := os.Getenv("HOST")
	db := database.New()

	sessions := scs.New()
	sessions.Store = postgresstore.New(db)
	sessions.Lifetime = 12 * time.Hour

	auth := auth.NewAuth(sessions, database.NewUserStore(db))
	NewServer := &Server{
		port:    port,
		db:      db,
		auth:    auth,
		meals:   services.NewMealsService(database.NewMealStore(db), auth),
		recipes: services.NewRecipeService(database.NewRecipeStore(db), auth),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
