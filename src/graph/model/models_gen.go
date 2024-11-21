// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Mutation struct {
}

type Query struct {
}

type User struct {
	ID                      string    `json:"id"`
	Name                    string    `json:"name"`
	Email                   string    `json:"email"`
	Location                *string   `json:"location,omitempty"`
	IsEmployer              *bool     `json:"isEmployer,omitempty"`
	Skills                  []*string `json:"skills,omitempty"`
	LookingForOpportunities *bool     `json:"lookingForOpportunities,omitempty"`
}
