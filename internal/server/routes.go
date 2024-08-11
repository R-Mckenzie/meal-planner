package server

import (
	"encoding/json"
	"log/slog"
	"net/http"

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

	// HOMEPAGE
	mux.Handle("/", s.auth.Sessions(s.noSurf(http.HandlerFunc(s.indexHandler))))
	mux.Handle("GET /dashboard", s.auth.Sessions(s.auth.RequireAuthentication(s.noSurf(http.HandlerFunc(s.dashboardPageHandler)))))
	// AUTH
	mux.Handle("GET /signup", s.auth.Sessions(s.noSurf(http.HandlerFunc(s.signupPageHandler))))
	mux.Handle("GET /login", s.auth.Sessions(s.noSurf(http.HandlerFunc(s.loginPageHandler))))
	mux.Handle("POST /signup", s.auth.Sessions(s.noSurf(http.HandlerFunc(s.auth.UserSignup))))
	mux.Handle("POST /login", s.auth.Sessions(s.noSurf(http.HandlerFunc(s.auth.UserLogin))))
	mux.Handle("POST /logout", s.auth.Sessions(s.noSurf(http.HandlerFunc(s.auth.UserLogout))))

	return mux
}

func (s *Server) dashboardPageHandler(w http.ResponseWriter, r *http.Request) {
	pages.DashboardPage(nosurf.Token(r)).Render(r.Context(), w)
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

func (s *Server) noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}
