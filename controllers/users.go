package controllers

import (
	"fmt"
	"net/http"
	// _ "github.com/gorilla/schema" // 
	"github.com/InSystem/lenslocked/views"
)

// NewUsers create signup view
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}
}

// Users is type for Users View
type Users struct {
	NewView *views.View
}

// New is used to create a form where user can create  anew account
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

// SignupForm is  New User sctruct
type SignupForm struct {
	Email    string `schema: "email"`
	Password string `schema: "password"`
}

// Create is used to process the signup form
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, form)
}
