package controllers

import (
	"github.com/InSystem/lenslocked/views"
)

// NewStatic used for generate  Static View
func NewStatic() *Static{
	return &Static{
		Home: views.NewView("bootstrap", "views/static/home.gohtml"),
		Contact: views.NewView("bootstrap", "views/static/contact.gohtml"),
	}
}

// Static type
type Static struct {
	Home *views.View
	Contact *views.View
}	