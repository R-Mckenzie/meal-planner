package auth

import (
	"context"
	"errors"
	"github.com/R-Mckenzie/mealplanner/cmd/web/components"
	"github.com/R-Mckenzie/mealplanner/cmd/web/pages"
	"github.com/R-Mckenzie/mealplanner/internal/database"
	"github.com/R-Mckenzie/mealplanner/internal/validation"
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	sessions *scs.SessionManager
	users    database.UserStore
}

func NewAuth(sessions *scs.SessionManager, userStore database.UserStore) *AuthService {
	return &AuthService{
		sessions: sessions,
		users:    userStore,
	}
}

func (a *AuthService) Sessions(next http.Handler) http.Handler {
	return a.sessions.LoadAndSave(next)
}

func (a *AuthService) Logout(ctx context.Context) error {
	err := a.sessions.RenewToken(ctx)
	if err != nil {
		return err
	}
	a.sessions.Remove(ctx, "authenticatedUserId")
	return nil
}

func (a *AuthService) RenewToken(ctx context.Context, userID int) error {
	err := a.sessions.RenewToken(ctx)
	if err != nil {
		return err
	}
	a.sessions.Put(ctx, "authenticatedUserId", userID)
	return nil
}

// Return true if the current request is from an authenticated user, otherwise return false
func (a *AuthService) IsAuthenticated(r *http.Request) bool {
	return a.sessions.Exists(r.Context(), "authenticatedUserId")
}

func (a *AuthService) AuthorisedUser(ctx context.Context) int {
	id, ok := a.sessions.Get(ctx, "authenticatedUserId").(int)
	if !ok {
		return -1
	}
	return id
}

func (a *AuthService) RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If not authenticated, redirect to login
		if !a.IsAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Otherwise set the "Cache-Control: no-store" header so that pages require
		// authentication are not stored in the users browser cache (or other intermediary cache).
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func (a *AuthService) UserSignup(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	passwordConfirmation := r.FormValue("password-confirmation")

	e := make(validation.Errors)
	// Validation
	validation.Validate("email", "Please enter a valid email", validation.IsEmail(email), e)
	validation.Validate("password", "Password must not be blank", validation.NotBlank(password), e)
	validation.Validate("password", "Password must be at least 8 characters", validation.MinChars(password, 8), e)
	validation.Validate("password-confirmation", "Passwords do not match", password == passwordConfirmation, e)

	if e.Any() {
		pages.SignupForm(pages.SignupFormValues{Email: email, Password: password, PasswordConfirmation: passwordConfirmation}, e).Render(r.Context(), w)
		return
	}

	// Save to DB
	err := a.users.Create(email, password, "template")
	if errors.Is(err, database.ErrUserExists) {
		e.Add("email", "Email already in use")
		pages.SignupForm(pages.SignupFormValues{CSRFToken: nosurf.Token(r), Email: email, Password: password, PasswordConfirmation: passwordConfirmation}, e).Render(r.Context(), w)
		return
	}
	if err != nil {
		components.Warning(err.Error()).Render(r.Context(), w)
		return
	}
	components.Success("Account created!").Render(r.Context(), w)
}

func (a *AuthService) UserLogin(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")
	e := make(validation.Errors)

	u, err := a.users.GetByEmail(email)
	if err != nil {
		slog.Info("Could not find user with email", "email", email)
	}
	err = bcrypt.CompareHashAndPassword(u.PassHash, []byte(password))
	if err != nil {
		validation.Validate("credentials", "Email and password do not match", false, e)
		slog.Info("Email and password do not match", "email", email)
		pages.LoginForm(pages.LoginFormValues{Email: email, Password: password, CSRFToken: nosurf.Token(r)}, e).Render(r.Context(), w)
		return
	}

	err = a.RenewToken(r.Context(), u.ID)
	if err != nil {
		slog.Info("There was a problem logging in", "email", email)
		components.Warning("There was a problem logginf in").Render(r.Context(), w)
	}

	w.Header().Add("HX-Redirect", "/")
}

func (a *AuthService) UserLogout(w http.ResponseWriter, r *http.Request) {
	err := a.Logout(r.Context())
	if err != nil {
		components.Warning("Error!").Render(r.Context(), w)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
