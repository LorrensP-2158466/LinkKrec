// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Company struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Location  *Location  `json:"location,omitempty"`
	Vacancies []*Vacancy `json:"vacancies"`
	Employees []*User    `json:"employees"`
}

type ConnectionRequest struct {
	ID              string `json:"id"`
	FromUser        *User  `json:"fromUser"`
	ConnectedToUser *User  `json:"connectedToUser"`
	Status          bool   `json:"status"`
}

type CreateLocationInput struct {
	Country     string `json:"country"`
	City        string `json:"city"`
	Street      string `json:"street"`
	HouseNumber string `json:"houseNumber"`
}

type CreateVacancyInput struct {
	Title                      string      `json:"title"`
	Description                string      `json:"description"`
	Location                   string      `json:"location"`
	StartDate                  string      `json:"startDate"`
	EndDate                    string      `json:"endDate"`
	Status                     bool        `json:"status"`
	RequiredDegreeType         DegreeType  `json:"requiredDegreeType"`
	RequiredDegreeField        DegreeField `json:"requiredDegreeField"`
	RequiredExperienceDuration int         `json:"requiredExperienceDuration"`
	RequiredSkills             []*string   `json:"requiredSkills"`
}

type EducationEntry struct {
	ID          string      `json:"id"`
	Institution string      `json:"institution"`
	Info        string      `json:"info"`
	Degree      DegreeType  `json:"degree"`
	Field       DegreeField `json:"field"`
}

type EducationEntryInput struct {
	Institution string `json:"institution"`
	Degree      string `json:"degree"`
}

type ExperienceEntry struct {
	ID             string         `json:"id"`
	Title          string         `json:"title"`
	Description    *string        `json:"description,omitempty"`
	ExperienceType ExperienceType `json:"experienceType"`
	StartDate      *string        `json:"startDate,omitempty"`
	EndDate        *string        `json:"endDate,omitempty"`
}

type ExperienceEntryInput struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	StartDate   *string `json:"startDate,omitempty"`
	EndDate     *string `json:"endDate,omitempty"`
}

type Location struct {
	ID          string `json:"id"`
	Country     string `json:"country"`
	City        string `json:"city"`
	Street      string `json:"street"`
	HouseNumber string `json:"houseNumber"`
}

type Mutation struct {
}

type Notification struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Message   string  `json:"message"`
	ForUser   *User   `json:"forUser"`
	CreatedAt *string `json:"createdAt,omitempty"`
}

type Query struct {
}

type RegisterUserInput struct {
	Name          string              `json:"name"`
	Email         string              `json:"email"`
	Password      string              `json:"password"`
	ProfileUpdate *UpdateProfileInput `json:"profileUpdate"`
}

type Subscription struct {
}

type UpdateProfileInput struct {
	Education                 []*EducationEntryInput `json:"education,omitempty"`
	Skills                    []*string              `json:"skills,omitempty"`
	IsLookingForOpportunities *bool                  `json:"isLookingForOpportunities,omitempty"`
}

type UpdateUserInput struct {
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Location *string `json:"location,omitempty"`
}

type UpdateVacancyInput struct {
	Title                      *string      `json:"title,omitempty"`
	Description                *string      `json:"description,omitempty"`
	Location                   *string      `json:"location,omitempty"`
	StartDate                  *string      `json:"startDate,omitempty"`
	EndDate                    *string      `json:"endDate,omitempty"`
	Status                     *bool        `json:"status,omitempty"`
	RequiredDegreeType         *DegreeType  `json:"requiredDegreeType,omitempty"`
	RequiredDegreeField        *DegreeField `json:"requiredDegreeField,omitempty"`
	RequiredExperienceDuration *int         `json:"requiredExperienceDuration,omitempty"`
	RequiredSkills             []*string    `json:"requiredSkills,omitempty"`
}

type User struct {
	ID                      string            `json:"id"`
	Name                    string            `json:"name"`
	Email                   string            `json:"email"`
	Location                *Location         `json:"location,omitempty"`
	Connections             []*User           `json:"connections,omitempty"`
	Education               []*EducationEntry `json:"education,omitempty"`
	Skills                  []*string         `json:"skills,omitempty"`
	LookingForOpportunities *bool             `json:"lookingForOpportunities,omitempty"`
	IsProfileComplete       *bool             `json:"isProfileComplete,omitempty"`
	Companies               []*Company        `json:"companies,omitempty"`
}

type Vacancy struct {
	ID                         string       `json:"id"`
	Title                      string       `json:"title"`
	Description                string       `json:"description"`
	Location                   string       `json:"location"`
	PostedBy                   *Company     `json:"postedBy"`
	StartDate                  *string      `json:"startDate,omitempty"`
	EndDate                    *string      `json:"endDate,omitempty"`
	Status                     *bool        `json:"status,omitempty"`
	RequiredDegreeType         *DegreeType  `json:"requiredDegreeType,omitempty"`
	RequiredDegreeField        *DegreeField `json:"requiredDegreeField,omitempty"`
	RequiredExperienceDuration *int         `json:"requiredExperienceDuration,omitempty"`
	RequiredSkills             []*string    `json:"requiredSkills,omitempty"`
}

type DegreeField string

const (
	DegreeFieldComputerScience DegreeField = "ComputerScience"
	DegreeFieldEngineering     DegreeField = "Engineering"
	DegreeFieldBusiness        DegreeField = "Business"
	DegreeFieldEconomics       DegreeField = "Economics"
	DegreeFieldMarketing       DegreeField = "Marketing"
	DegreeFieldFinance         DegreeField = "Finance"
	DegreeFieldMedicine        DegreeField = "Medicine"
	DegreeFieldLaw             DegreeField = "Law"
	DegreeFieldPsychology      DegreeField = "Psychology"
)

var AllDegreeField = []DegreeField{
	DegreeFieldComputerScience,
	DegreeFieldEngineering,
	DegreeFieldBusiness,
	DegreeFieldEconomics,
	DegreeFieldMarketing,
	DegreeFieldFinance,
	DegreeFieldMedicine,
	DegreeFieldLaw,
	DegreeFieldPsychology,
}

func (e DegreeField) IsValid() bool {
	switch e {
	case DegreeFieldComputerScience, DegreeFieldEngineering, DegreeFieldBusiness, DegreeFieldEconomics, DegreeFieldMarketing, DegreeFieldFinance, DegreeFieldMedicine, DegreeFieldLaw, DegreeFieldPsychology:
		return true
	}
	return false
}

func (e DegreeField) String() string {
	return string(e)
}

func (e *DegreeField) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DegreeField(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DegreeField", str)
	}
	return nil
}

func (e DegreeField) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type DegreeType string

const (
	DegreeTypeNone         DegreeType = "None"
	DegreeTypeProfBachelor DegreeType = "ProfBachelor"
	DegreeTypeAcBachelor   DegreeType = "AcBachelor"
	DegreeTypeMaster       DegreeType = "Master"
	DegreeTypePhD          DegreeType = "PhD"
)

var AllDegreeType = []DegreeType{
	DegreeTypeNone,
	DegreeTypeProfBachelor,
	DegreeTypeAcBachelor,
	DegreeTypeMaster,
	DegreeTypePhD,
}

func (e DegreeType) IsValid() bool {
	switch e {
	case DegreeTypeNone, DegreeTypeProfBachelor, DegreeTypeAcBachelor, DegreeTypeMaster, DegreeTypePhD:
		return true
	}
	return false
}

func (e DegreeType) String() string {
	return string(e)
}

func (e *DegreeType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DegreeType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DegreeType", str)
	}
	return nil
}

func (e DegreeType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type ExperienceType string

const (
	ExperienceTypeIt          ExperienceType = "IT"
	ExperienceTypeEngineering ExperienceType = "Engineering"
	ExperienceTypeSales       ExperienceType = "Sales"
	ExperienceTypeHr          ExperienceType = "HR"
	ExperienceTypeConsultancy ExperienceType = "Consultancy"
	ExperienceTypeResearch    ExperienceType = "Research"
	ExperienceTypeMarketing   ExperienceType = "Marketing"
	ExperienceTypeFinance     ExperienceType = "Finance"
	ExperienceTypeCustomer    ExperienceType = "Customer"
	ExperienceTypeSupport     ExperienceType = "Support"
	ExperienceTypeOperation   ExperienceType = "Operation"
)

var AllExperienceType = []ExperienceType{
	ExperienceTypeIt,
	ExperienceTypeEngineering,
	ExperienceTypeSales,
	ExperienceTypeHr,
	ExperienceTypeConsultancy,
	ExperienceTypeResearch,
	ExperienceTypeMarketing,
	ExperienceTypeFinance,
	ExperienceTypeCustomer,
	ExperienceTypeSupport,
	ExperienceTypeOperation,
}

func (e ExperienceType) IsValid() bool {
	switch e {
	case ExperienceTypeIt, ExperienceTypeEngineering, ExperienceTypeSales, ExperienceTypeHr, ExperienceTypeConsultancy, ExperienceTypeResearch, ExperienceTypeMarketing, ExperienceTypeFinance, ExperienceTypeCustomer, ExperienceTypeSupport, ExperienceTypeOperation:
		return true
	}
	return false
}

func (e ExperienceType) String() string {
	return string(e)
}

func (e *ExperienceType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = ExperienceType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid ExperienceType", str)
	}
	return nil
}

func (e ExperienceType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
