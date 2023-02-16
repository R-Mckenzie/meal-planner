package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/R-Mckenzie/meal-planner/models"
	"github.com/R-Mckenzie/meal-planner/validation"
	"github.com/R-Mckenzie/meal-planner/views"
	"github.com/justinas/nosurf"
)

type User struct {
	SignupView *views.View
	LoginView  *views.View
	us         models.UserService
	iLog       *log.Logger
	eLog       *log.Logger
}

type SignupForm struct {
	email    string
	password string
}

func NewUsers(us models.UserService, iLog, eLog *log.Logger) *User {
	return &User{
		SignupView: views.NewView("root", "views/users/signup.html"),
		LoginView:  views.NewView("root", "views/users/login.html"),
		us:         us,
		iLog:       iLog,
		eLog:       eLog,
	}
}

func (u *User) SignupPage(w http.ResponseWriter, r *http.Request) {
	u.SignupView.Data.User = r.Context().Value("mealplanner_current_user").(int) >= 0
	u.SignupView.SetAlert("", views.Success)
	u.SignupView.Data.CSRFtoken = nosurf.Token(r)

	m, t, err := getAlertData(w, r)
	if err != nil {
		u.eLog.Println("in SignupPage getting alert data: ", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
	u.SignupView.SetAlert(m, t)

	if err = u.SignupView.Render(w); err != nil {
		u.eLog.Println("in SignupPage rendering view: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (u *User) Signup(w http.ResponseWriter, r *http.Request) {
	u.SignupView.Data.User = r.Context().Value("mealplanner_current_user").(int) >= 0

	email, pass, err := parseUserForm(r)
	if err != nil {
		u.eLog.Println("in SignupPage parsing form: ", err)
		w.WriteHeader(http.StatusBadRequest)
		u.SignupView.SetAlert("There was a problem with your input", views.Error)
		if err := u.SignupView.Render(w); err != nil {
			u.eLog.Println("in SignupPage rendering view: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	// Validate the data
	valid, faults := validation.PasswordCheck(pass)
	validEmail := validation.IsEmail(email)
	if !validEmail {
		faults = append(faults, "Must be a valid email")
	}
	message := strings.Join(faults, "\n")

	if !valid || !validEmail {
		u.SignupView.SetAlert(message, views.Error)
		if err := u.SignupView.Render(w); err != nil {
			u.eLog.Println("in SignupPage rendering view: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Create the user
	if err := u.us.Create(email, pass); err != nil {
		u.eLog.Println("in SignupPage creating new user: ", err)
		w.WriteHeader(http.StatusConflict)
		u.SignupView.Data.Alert = views.Alert{Type: views.Error, Message: err.Error()}
		if err := u.SignupView.Render(w); err != nil {
			u.eLog.Println("in SignupPage rendering view: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	u.iLog.Printf("Successfully signed up user %q", email)
	setAlertData(w, fmt.Sprintf("Successfully added user %q", email), views.Success)
	http.Redirect(w, r, "/signup", http.StatusSeeOther)
}

func (u *User) LoginPage(w http.ResponseWriter, r *http.Request) {
	u.LoginView.Data.User = r.Context().Value("mealplanner_current_user").(int) >= 0
	u.LoginView.Data.CSRFtoken = nosurf.Token(r)
	u.LoginView.SetAlert("", views.Success)
	if err := u.LoginView.Render(w); err != nil {
		http.Error(w, "There was a problem...", http.StatusInternalServerError)
	}
}

func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	// Parse form
	email, pass, err := parseUserForm(r)
	if err != nil {
		u.LoginView.SetAlert("There was a problem with your input", views.Error)
		w.WriteHeader(http.StatusBadRequest)
		if err := u.LoginView.Render(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Authenticate user
	user, err := u.us.Authenticate(email, pass)
	if err != nil {
		log.Println(err)
		u.LoginView.SetAlert(err.Error(), views.Error)
		w.WriteHeader(http.StatusUnauthorized)
		if err := u.LoginView.Render(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Create a new remember token for the user on every login
	err = u.updateRememberToken(w, user)
	if err != nil {
		log.Printf(err.Error())
		http.Error(w, "Problem updated remember cookie", http.StatusInternalServerError)
		return
	}

	u.iLog.Printf("User %q logged in", email)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (u *User) Logout(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("mealplanner_current_user").(int)
	if !ok {
		return
	}

	user, err := u.us.ByID(userID)
	if err != nil {
		u.eLog.Println("in logout, could not get user by ID: ", err)
	}

	clearRememberCookie(w)
	u.iLog.Printf("User %q logged out", user.Email)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func parseUserForm(r *http.Request) (string, string, error) {
	if err := r.ParseForm(); err != nil {
		return "", "", err
	}

	email := r.PostForm["email"][0]
	password := r.PostForm["password"][0]
	email = strings.ToLower(email)

	return email, password, nil
}

// ===== HELPERS =====

func (uc *User) updateRememberToken(w http.ResponseWriter, u *models.User) error {
	if err := uc.us.GenerateRemember(u); err != nil {
		log.Println("Problem generating remember token: ", err)
		return err
	}
	cookie := &http.Cookie{
		Name:     "remember_token",
		Value:    u.Remember,
		HttpOnly: true,
	}

	http.SetCookie(w, cookie)
	return nil
}

func clearRememberCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "remember_token",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
	})
}
