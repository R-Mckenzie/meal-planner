package controllers

import (
	"log"
	"net/http"

	"github.com/R-Mckenzie/meal-planner/models"
	"github.com/R-Mckenzie/meal-planner/views"
)

type User struct {
	SignupView *views.View
	LoginVuew  *views.View
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
		LoginVuew:  views.NewView("root", "views/users/login.html"),
		us:         &us,
		errLog:     &errLog,
		infLog:     &infLog,
	}
}

func (u *User) SignupPage(w http.ResponseWriter, r *http.Request) {
	a := &views.Alert{Type: views.AlertSuccess, Message: "Successfully created user"}
	// We use query params on the redirect when we successfully sign up
	if r.URL.Query().Get("success") != "true" {
		a.Message = ""
	}
	err := u.SignupView.Render(w, a)
	if err != nil {
		panic(err)
	}
}

func (u *User) Signup(w http.ResponseWriter, r *http.Request) {
	a := &views.Alert{Type: views.AlertSuccess, Message: "success"}

	if err := r.ParseForm(); err != nil {
		u.errLog.Println("Error parsing signup form: ", err)
		a.Message = "There was a problem with your input"
		a.Type = views.AlertError
		w.WriteHeader(http.StatusBadRequest)
		if err := u.SignupView.Render(w, a); err != nil {
			panic(err)
		}
		return
	}

	email := r.PostForm["email"][0]
	password := r.PostForm["password"][0]
	if err := u.us.Create(email, password); err != nil {
		u.errLog.Println("Error creating new user: ", err)
		a.Message = err.Error() // This error is client readable from the service
		a.Type = views.AlertError
		w.WriteHeader(http.StatusConflict)
		if err := u.SignupView.Render(w, a); err != nil {
			panic(err)
		}
		return
	}
	http.Redirect(w, r, "/signup?success=true", http.StatusSeeOther)
}

func (u *User) LoginPage(w http.ResponseWriter, r *http.Request) {
	err := u.LoginVuew.Render(w, nil)
	if err != nil {
		panic(err)
	}
}
