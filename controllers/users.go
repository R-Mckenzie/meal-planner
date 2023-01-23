package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/R-Mckenzie/meal-planner/models"
	"github.com/R-Mckenzie/meal-planner/views"
)

type User struct {
	SignupView *views.View
	LoginView  *views.View
	us         *models.UserService
	errLog     *log.Logger
	infLog     *log.Logger
}

type SignupForm struct {
	email    string
	password string
}

func NewUsers(us models.UserService, errLog, infLog log.Logger) *User {
	return &User{
		SignupView: views.NewView("root", "views/users/signup.html"),
		LoginView:  views.NewView("root", "views/users/login.html"),
		us:         &us,
		errLog:     &errLog,
		infLog:     &infLog,
	}
}

func (u *User) SignupPage(w http.ResponseWriter, r *http.Request) {
	a := &views.Alert{Type: views.AlertSuccess, Message: "Successfully created user"}
	if r.URL.Query().Get("success") != "true" {
		a.Message = ""
	}
	err := u.SignupView.Render(w, a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (u *User) Signup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		u.errLog.Println("Error parsing signup form: ", err)
		w.WriteHeader(http.StatusBadRequest)
		if err := u.SignupView.Render(w, views.Alert{Type: views.AlertError, Message: "There was a problem with your input"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	email := r.PostForm["email"][0]
	password := r.PostForm["password"][0]
	if _, err := u.us.Create(email, password); err != nil {
		u.errLog.Println("Error creating new user: ", err)
		w.WriteHeader(http.StatusConflict)
		if err := u.SignupView.Render(w, views.Alert{Type: views.AlertError, Message: err.Error()}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	http.Redirect(w, r, "/signup?success=true", http.StatusSeeOther)
}

func (u *User) LoginPage(w http.ResponseWriter, r *http.Request) {
	err := u.LoginView.Render(w, nil)
	if err != nil {
		panic(err)
	}
}

func (u *User) Login(w http.ResponseWriter, r *http.Request) {
	u.infLog.Println("login attempt")
	if err := r.ParseForm(); err != nil {
		u.errLog.Println("Error parsing login form: ", err)
		a := views.Alert{Type: views.AlertError, Message: "There was a problem with your input"}
		w.WriteHeader(http.StatusBadRequest)
		if err := u.LoginView.Render(w, a); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	email := r.PostForm["email"][0]
	password := r.PostForm["password"][0]
	user, err := u.us.Authenticate(email, password)
	if err != nil {
		u.errLog.Printf("There was a problem authenticating user %q. %s\n", email, err.Error())
		a := views.Alert{Type: views.AlertError, Message: err.Error()}
		w.WriteHeader(http.StatusUnauthorized)
		if err := u.LoginView.Render(w, a); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	u.infLog.Printf("Successfully logged in user %q\n", user.Email)

	// Create a new remember token for the user on every login
	rToken, err := r.Cookie("remember_token")
	if err != nil {
		u.errLog.Println(err)
		user.Remember = ""
	} else {
		user.Remember = rToken.Value
	}
	err = u.us.GenerateRemember(user)
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	if err != nil {
		u.errLog.Println("Problem generating remember token: ", err)
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (u *User) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		u.errLog.Println("Problem finding remember token cookie")
		http.Error(w, "There was a problem identifying yo,", http.StatusInternalServerError)
		return
	}

	user, err := u.us.GetByRemember(cookie.Value)
	if err != nil {
		u.errLog.Printf("Problem finding user by remember token %q. %v\n", cookie.Value, err)
	}
	w.Write([]byte(fmt.Sprintf("%+v", user)))

}
