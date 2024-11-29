// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type AskedConnection struct {
	User        *User `json:"user"`
	ConnectedTo *User `json:"connectedTo"`
	Status      bool  `json:"status"`
}

type CreateVacancyInput struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Location    *string `json:"location,omitempty"`
	StartDate   *string `json:"startDate,omitempty"`
	EndDate     *string `json:"endDate,omitempty"`
}

type EducationEntry struct {
	Institution string      `json:"institution"`
	Info        string      `json:"info"`
	Degree      DegreeType  `json:"degree"`
	Field       DegreeField `json:"field"`
}

type EducationEntryInput struct {
	Institution string `json:"institution"`
	Degree      string `json:"degree"`
}

type Employer struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Location  *string    `json:"location,omitempty"`
	Vacancies []*Vacancy `json:"vacancies"`
	Employees []*User    `json:"employees"`
}

type ExperienceEntry struct {
	Title          string         `json:"title"`
	ExperienceType ExperienceType `json:"experienceType"`
	Description    *string        `json:"description,omitempty"`
	StartDate      *string        `json:"startDate,omitempty"`
	EndDate        *string        `json:"endDate,omitempty"`
}

type ExperienceEntryInput struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	StartDate   *string `json:"startDate,omitempty"`
	EndDate     *string `json:"endDate,omitempty"`
}

type Mutation struct {
}

type Notification struct {
	ID        string  `json:"id"`
	ForUser   *User   `json:"forUser"`
	Message   string  `json:"message"`
	CreatedAt *string `json:"createdAt,omitempty"`
}

type Query struct {
}

type RegisterUserInput struct {
	Name          string              `json:"name"`
	Email         string              `json:"email"`
	Password      string              `json:"password"`
	IsEmployer    *bool               `json:"isEmployer,omitempty"`
	ProfileUpdate *UpdateProfileInput `json:"profileUpdate"`
}

type Subscription struct {
}

type UpdateProfileInput struct {
	Education                 []*EducationEntryInput  `json:"education,omitempty"`
	Experience                []*ExperienceEntryInput `json:"experience,omitempty"`
	Skills                    []*string               `json:"skills,omitempty"`
	IsLookingForOpportunities *bool                   `json:"isLookingForOpportunities,omitempty"`
}

type UpdateUserInput struct {
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Location *string `json:"location,omitempty"`
}

type User struct {
	ID                        string             `json:"id"`
	Name                      string             `json:"name"`
	Email                     string             `json:"email"`
	Location                  *string            `json:"location,omitempty"`
	IsEmployer                *bool              `json:"isEmployer,omitempty"`
	Connections               []*User            `json:"connections,omitempty"`
	Education                 []*EducationEntry  `json:"education,omitempty"`
	Experience                []*ExperienceEntry `json:"experience,omitempty"`
	Skills                    []*string          `json:"skills,omitempty"`
	IsLookingForOpportunities *bool              `json:"isLookingForOpportunities,omitempty"`
}

type Vacancy struct {
	ID                          string           `json:"id"`
	Title                       string           `json:"title"`
	Description                 string           `json:"description"`
	RequiredEducation           DegreeType       `json:"requiredEducation"`
	RequiredExperiences         []ExperienceType `json:"requiredExperiences"`
	RequiredExperienceDurations []int            `json:"requiredExperienceDurations"`
	Location                    string           `json:"location"`
	PostedBy                    *Employer        `json:"postedBy"`
	StartDate                   *string          `json:"startDate,omitempty"`
	EndDate                     *string          `json:"endDate,omitempty"`
	Status                      *string          `json:"status,omitempty"`
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
