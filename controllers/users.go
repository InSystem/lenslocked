package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/InSystem/lenslocked/rand"

	"github.com/InSystem/lenslocked/models"

	"github.com/InSystem/lenslocked/views"
	_ "github.com/gorilla/schema"
)

// Users is type for Users View
type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        models.UserService
}

// NewUsers create signup view
func NewUsers(us models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

// New is used to create a form where user can create  anew account
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	// d := views.Data{
	// 	Alert: &views.Alert{
	// 		Level:   views.AlertLevelWarning,
	// 		Message: "something went wrong",
	// 	},
	// 	Yield: "Hello",
	// }

	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

// SignupForm is  New User sctruct
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Create is used to process the signup form
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data

	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		//nolint:errcheck
		u.NewView.Render(w, vd)
		return 
	}

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		//nolint:errcheck
		u.NewView.Render(w, vd)
		return
	}

	err := u.signIn(w, &user)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/cookietest", http.StatusFound)
	fmt.Fprintln(w, user)
}

// SignupForm is  New User sctruct
type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Login is used to verify the provided email and password
// and log in if they are correct
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	vd := views.Data{}
	var form LoginForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		//nolint:errcheck
		u.LoginView.Render(w, vd)
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Password)

	if err != nil {
		switch err {
		case models.ErrorNotFound:
			vd.AlertError("Invalid email adress")
		default:
			vd.SetAlert(err)
		}
		//nolint:errcheck
		u.LoginView.Render(w, vd)
	} else {
		err := u.signIn(w, user)
		if err != nil {
			vd.SetAlert(err)
			//nolint:errcheck
			u.LoginView.Render(w, vd)
			return
		}
		http.Redirect(w, r, "/cookietest", http.StatusFound)
	}
}

// CookieTest is used to display cookies to the current user
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, user)
}

func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}

		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}

	http.SetCookie(w, &cookie)
	return nil
}
