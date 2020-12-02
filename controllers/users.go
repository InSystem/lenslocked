package controllers

import (
	"fmt"
	"net/http"

	"github.com/InSystem/lenslocked/views"
)

// NewUsers create signup view
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}
}

type Users struct {
	NewView *views.View
}

// This is used to create a form where user can create  anew account
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

// This is usef to process the signup form
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "we are created  new user!")
}
