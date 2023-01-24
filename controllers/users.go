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
	us         models.UserService
}

type SignupForm struct {
	email    string
	password string
}

func NewUsers(us models.UserService) *User {
	return &User{
		SignupView: views.NewView("root", "views/users/signup.html"),
		LoginView:  views.NewView("root", "views/users/login.html"),
		us:         us,
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
		w.WriteHeader(http.StatusBadRequest)
		if err := u.SignupView.Render(w, views.Alert{Type: views.AlertError, Message: "There was a problem with your input"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	email := r.PostForm["email"][0]
	password := r.PostForm["password"][0]
	if err := u.us.Create(email, password); err != nil {
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
	if err := r.ParseForm(); err != nil {
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
		log.Println(err)
		a := views.Alert{Type: views.AlertError, Message: err.Error()}
		w.WriteHeader(http.StatusUnauthorized)
		if err := u.LoginView.Render(w, a); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Create a new remember token for the user on every login
	rToken, err := r.Cookie("remember_token")
	if err != nil {
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
		log.Println("Problem generating remember token: ", err)
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (u *User) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, "There was a problem identifying yo,", http.StatusInternalServerError)
		return
	}

	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		log.Printf("Problem finding user by remember token %q. %v\n", cookie.Value, err)
	}
	w.Write([]byte(fmt.Sprintf("%+v", user)))
}
