package graph

import (
	"github.com/gorilla/sessions"
	"github.com/knakk/sparql"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Repo       *sparql.Repo
	UpdateRepo *sparql.Repo
	Store      *sessions.CookieStore
}
