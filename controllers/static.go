package controllers

import (
	"github.com/InSystem/lenslocked/views"
)

// NewStatic used for generate  Static View
func NewStatic() *Static{
	return &Static{
		Home: views.NewView("bootstrap", "static/home"),
		Contact: views.NewView("bootstrap", "static/contact"),
	}
}

// Static type
type Static struct {
	Home *views.View
	Contact *views.View
}	