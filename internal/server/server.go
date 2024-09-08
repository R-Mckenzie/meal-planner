package server

import (
	"database/sql"
	"fmt"
	"log/slog"
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
	newServer := &Server{
		port:    port,
		db:      db,
		auth:    auth,
		meals:   services.NewMealsService(database.NewMealStore(db), auth),
		recipes: services.NewRecipeService(database.NewRecipeStore(db), auth),
	}

	// recipes, err := newServer.recipes.ByUser(1)
	// if err != nil {
	// 	slog.Error(err.Error())
	// }
	// for _, r := range recipes {
	// 	slog.Info("recipe", "r", r)
	// }

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	slog.Info("Starting server:", "host", host, "port", port)

	return server
}
