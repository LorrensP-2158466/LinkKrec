package model

type User struct {
	ID                      string             `json:"id"`
	Name                    string             `json:"name"`
	Email                   string             `json:"email"`
	Location                *string            `json:"location,omitempty"`
	IsEmployer              *bool              `json:"isEmployer,omitempty"`
	Connections             []*User            `json:"connections,omitempty"`
	ConnectionIds           []string           `json:"-"`
	Education               []*EducationEntry  `json:"education,omitempty"`
	Experience              []*ExperienceEntry `json:"experience,omitempty"`
	Skills                  []*string          `json:"skills,omitempty"`
	LookingForOpportunities *bool              `json:"lookingForOpportunities,omitempty"`
}
