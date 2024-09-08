package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/R-Mckenzie/mealplanner/cmd/web"
	"github.com/R-Mckenzie/mealplanner/cmd/web/pages"
	"github.com/R-Mckenzie/mealplanner/internal/database"

	"github.com/justinas/nosurf"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/health", s.healthHandler)

	fileServer := http.FileServer(http.FS(web.Files))
	mux.Handle("/assets/", fileServer)
	mux.HandleFunc("/favicon.ico", favicon)

	// HOMEPAGE
	mux.Handle("/", s.auth.Sessions(http.HandlerFunc(s.indexHandler)))
	mux.Handle("GET /dashboard", s.auth.Sessions(s.auth.RequireAuthentication(http.HandlerFunc(s.dashboardPageHandler))))

	// RECIPES
	mux.Handle("GET /recipes", s.auth.Sessions(s.auth.RequireAuthentication(http.HandlerFunc(s.recipesPageHandler))))
	mux.Handle("POST /recipes/create", s.auth.Sessions(s.auth.RequireAuthentication(http.HandlerFunc(s.recipes.CreateRecipeHandler))))
	mux.Handle("GET /recipes/{id}", s.auth.Sessions(s.auth.RequireAuthentication(http.HandlerFunc(s.recipes.GetRecipeHandler))))
	mux.Handle("PUT /recipes/{id}", s.auth.Sessions(s.auth.RequireAuthentication(http.HandlerFunc(s.recipes.UpdateRecipeHandler))))
	mux.Handle("DELETE /recipes/delete/{id}", s.auth.Sessions(s.auth.RequireAuthentication(http.HandlerFunc(s.recipes.DeleteRecipeHandler))))

	// LOGIN/SIGNUP
	mux.Handle("GET /signup", s.auth.Sessions(http.HandlerFunc(s.signupPageHandler)))
	mux.Handle("GET /login", s.auth.Sessions(http.HandlerFunc(s.loginPageHandler)))

	// AUTH
	mux.Handle("POST /signup", s.auth.Sessions(http.HandlerFunc(s.auth.UserSignup)))
	mux.Handle("POST /login", s.auth.Sessions(http.HandlerFunc(s.auth.UserLogin)))
	mux.Handle("POST /logout", s.auth.Sessions(s.auth.RequireAuthentication(http.HandlerFunc(s.auth.UserLogout))))

	// MEALS
	mux.Handle("POST /meals", s.auth.Sessions(s.auth.RequireAuthentication(http.HandlerFunc(s.meals.PostMealsHandler))))
	mux.Handle("DELETE /meals/delete/{id}", s.auth.Sessions(s.auth.RequireAuthentication(http.HandlerFunc(s.meals.DeleteMealsHandler))))

	csrfHandler := nosurf.New(mux)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	csrfHandler.ExemptGlob("/meals/delete/*")
	csrfHandler.ExemptGlob("/recipes/delete/*")
	return logRequest(csrfHandler)
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, web.Files, "assets/favicon.ico")
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info(fmt.Sprintf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL))
		handler.ServeHTTP(w, r)
	})
}

func (s *Server) recipesPageHandler(w http.ResponseWriter, r *http.Request) {
	recipeData := s.recipes.BuildUserRecipeData(s.auth.AuthorisedUser(r.Context()), r.Context())
	pages.RecipesPage(recipeData, nosurf.Token(r)).Render(r.Context(), w)
}

func (s *Server) dashboardPageHandler(w http.ResponseWriter, r *http.Request) {
	weekStart := r.URL.Query().Get("week_start")
	var targetDate time.Time
	if weekStart != "" {
		date, err := time.Parse(time.DateOnly, weekStart)
		targetDate = date
		if err != nil {
			slog.Error(err.Error())
			return
		}

	} else {
		now := time.Now()
		dayOffset := now.Weekday() - time.Monday
		targetDate = now.AddDate(0, 0, -int(dayOffset))
		if targetDate.After(now) {
			targetDate = targetDate.AddDate(0, 0, -7)
		}
	}

	recipes, err := s.recipes.ByUser(s.auth.AuthorisedUser(r.Context()))
	if err != nil {
		slog.Error(err.Error())
		// can't get recipes
	}

	recipeDatas := []pages.MealData{}
	for _, r := range recipes {
		recipeDatas = append(recipeDatas, pages.MealData{RecipeName: r.Title, RecipeID: r.ID})
	}

	meals, err := s.meals.GetMealsByDateRange(targetDate, targetDate.Add(time.Hour*24*7))
	if err != nil {
		slog.Error(err.Error())
	}

	mealDatas := make(map[time.Weekday][]pages.MealData)
	for _, m := range meals {
		mealDatas[m.MealDate.Weekday()] = append(mealDatas[m.MealDate.Weekday()], pages.MealData{RecipeName: m.RecipeTitle, RecipeID: m.RecipeID, MealID: m.MealID})
	}

	data := pages.DashboardData{WeekStartDate: targetDate, Recipes: recipeDatas, Meals: mealDatas}

	if weekStart != "" {
		pages.Dashboard(data, nosurf.Token(r)).Render(r.Context(), w)
		return
	}
	pages.DashboardPage(data, nosurf.Token(r)).Render(r.Context(), w)
}

func (s *Server) signupPageHandler(w http.ResponseWriter, r *http.Request) {
	pages.SignupIndex(pages.SignupPageData{FormValues: pages.SignupFormValues{CSRFToken: nosurf.Token(r)}}).Render(r.Context(), w)
}

func (s *Server) loginPageHandler(w http.ResponseWriter, r *http.Request) {
	pages.LoginIndex(pages.LoginIndexPageData{FormValues: pages.LoginFormValues{CSRFToken: nosurf.Token(r)}}).Render(r.Context(), w)
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	authed := s.auth.IsAuthenticated(r)

	if r.URL.Path != "/" {
		pages.Error404(authed, nosurf.Token(r)).Render(r.Context(), w)
		return
	}

	pages.Homepage(authed, nosurf.Token(r)).Render(r.Context(), w)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, err := json.Marshal(database.Health(s.db))

	if err != nil {
		slog.Error("error handling JSON marshal", "Error", err)
	}

	_, _ = w.Write(jsonResp)
}
