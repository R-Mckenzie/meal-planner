package server

import (
	"encoding/json"
	"log"
	"net/http"

	"go-app-template/cmd/web"
	"go-app-template/cmd/web/pages"
	"go-app-template/internal/database"

	"github.com/justinas/nosurf"
)

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("/health", s.healthHandler)

	fileServer := http.FileServer(http.FS(web.Files))
	mux.Handle("/assets/", fileServer)

	// HOMEPAGE
	mux.Handle("/", s.auth.Sessions(s.noSurf(http.HandlerFunc(s.indexHandler))))

	// AUTH
	// -- Static Pages
	mux.Handle("GET /signup", s.auth.Sessions(s.noSurf(http.HandlerFunc(s.signupPageHandler))))
	mux.Handle("GET /login", s.auth.Sessions(s.noSurf(http.HandlerFunc(s.loginPageHandler))))

	mux.Handle("GET /locked", s.auth.Sessions(s.auth.RequireAuthentication(s.noSurf(http.HandlerFunc(s.authorisedPageHandler)))))

	// --Email/Password Handlers
	mux.Handle("POST /signup", s.auth.Sessions(s.noSurf(http.HandlerFunc(s.auth.UserSignup))))
	mux.Handle("POST /login", s.auth.Sessions(s.noSurf(http.HandlerFunc(s.auth.UserLogin))))
	mux.Handle("POST /logout", s.auth.Sessions(s.noSurf(http.HandlerFunc(s.auth.UserLogout))))

	return mux
}

func (s *Server) authorisedPageHandler(w http.ResponseWriter, r *http.Request) {
	pages.AuthorisedPage(nosurf.Token(r)).Render(r.Context(), w)
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
		log.Fatalf("error handling JSON marshal. Err: %v", err)
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
