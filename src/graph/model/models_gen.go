// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"LinkKrec/graph/scalar"
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

type CreateCompanyInput struct {
	Name     string               `json:"name"`
	Email    string               `json:"email"`
	Location *CreateLocationInput `json:"location"`
}

type CreateLocationInput struct {
	Country     string `json:"country"`
	City        string `json:"city"`
	Street      string `json:"street"`
	HouseNumber string `json:"houseNumber"`
}

type CreateVacancyInput struct {
	Title               string                  `json:"title"`
	Description         string                  `json:"description"`
	Location            *CreateLocationInput    `json:"location"`
	StartDate           string                  `json:"startDate"`
	EndDate             string                  `json:"endDate"`
	Status              bool                    `json:"status"`
	RequiredDegreeType  DegreeType              `json:"requiredDegreeType"`
	RequiredDegreeField DegreeField             `json:"requiredDegreeField"`
	RequiredSkills      []*string               `json:"requiredSkills"`
	RequiredExperience  []*ExperienceEntryInput `json:"requiredExperience"`
}

type DateInterval struct {
	Start scalar.Date `json:"start"`
	End   scalar.Date `json:"end"`
}

type EducationEntry struct {
	ID          string      `json:"id"`
	Institution string      `json:"institution"`
	ExtraInfo   *string     `json:"extra_info,omitempty"`
	From        scalar.Date `json:"from"`
	Till        scalar.Date `json:"till"`
	Degree      DegreeType  `json:"degree"`
	Field       DegreeField `json:"field"`
}

type EducationEntryInput struct {
	Institution string      `json:"institution"`
	ExtraInfo   *string     `json:"extra_info,omitempty"`
	From        scalar.Date `json:"from"`
	Till        scalar.Date `json:"till"`
	Degree      DegreeType  `json:"degree"`
	Field       DegreeField `json:"field"`
}

type EmployeeIds struct {
	Ids []string `json:"Ids,omitempty"`
}

type Experience struct {
	ID               string `json:"id"`
	Label            string `json:"label"`
	DurationInMonths int    `json:"durationInMonths"`
}

type ExperienceEntryInput struct {
	ID               string `json:"id"`
	DurationInMonths int    `json:"durationInMonths"`
}

type Location struct {
	ID          string `json:"id"`
	Country     string `json:"country"`
	City        string `json:"city"`
	Street      string `json:"street"`
	HouseNumber string `json:"houseNumber"`
}

type LocationFilter struct {
	Country string `json:"country"`
	City    string `json:"city"`
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
	IsEmployer    bool                `json:"isEmployer"`
	ProfileUpdate *UpdateProfileInput `json:"profileUpdate,omitempty"`
}

type Skill struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type Subscription struct {
}

type UpdateCompanyInput struct {
	ID       string               `json:"id"`
	Name     *string              `json:"name,omitempty"`
	Email    *string              `json:"email,omitempty"`
	Location *CreateLocationInput `json:"location,omitempty"`
}

type UpdateProfileInput struct {
	Education                 []*EducationEntryInput  `json:"education,omitempty"`
	Experience                []*ExperienceEntryInput `json:"experience,omitempty"`
	Skills                    []string                `json:"skills,omitempty"`
	IsLookingForOpportunities bool                    `json:"isLookingForOpportunities"`
	Country                   string                  `json:"country"`
	City                      string                  `json:"city"`
	Streetname                string                  `json:"streetname"`
	Housenumber               string                  `json:"housenumber"`
}

type UpdateVacancyInput struct {
	Title               *string                 `json:"title,omitempty"`
	Description         *string                 `json:"description,omitempty"`
	Location            *CreateLocationInput    `json:"location,omitempty"`
	StartDate           *string                 `json:"startDate,omitempty"`
	EndDate             *string                 `json:"endDate,omitempty"`
	Status              *bool                   `json:"status,omitempty"`
	RequiredDegreeType  *DegreeType             `json:"requiredDegreeType,omitempty"`
	RequiredDegreeField *DegreeField            `json:"requiredDegreeField,omitempty"`
	RequiredSkills      []*string               `json:"requiredSkills,omitempty"`
	RequiredExperience  []*ExperienceEntryInput `json:"requiredExperience,omitempty"`
}

type User struct {
	ID                      string            `json:"id"`
	Name                    string            `json:"name"`
	Email                   string            `json:"email"`
	Location                *Location         `json:"location"`
	Connections             []*User           `json:"connections,omitempty"`
	Education               []*EducationEntry `json:"education,omitempty"`
	Skills                  []*Skill          `json:"skills,omitempty"`
	Experiences             []*Experience     `json:"experiences,omitempty"`
	LookingForOpportunities bool              `json:"lookingForOpportunities"`
	IsProfileComplete       *bool             `json:"isProfileComplete,omitempty"`
	Companies               []*Company        `json:"companies,omitempty"`
}

type Vacancy struct {
	ID                  string        `json:"id"`
	Title               string        `json:"title"`
	Description         string        `json:"description"`
	Location            *Location     `json:"location"`
	PostedBy            *Company      `json:"postedBy"`
	StartDate           string        `json:"startDate"`
	EndDate             string        `json:"endDate"`
	Status              bool          `json:"status"`
	RequiredDegreeType  *DegreeType   `json:"requiredDegreeType,omitempty"`
	RequiredDegreeField *DegreeField  `json:"requiredDegreeField,omitempty"`
	RequiredExperience  []*Experience `json:"requiredExperience,omitempty"`
	RequiredSkills      []*Skill      `json:"requiredSkills,omitempty"`
}

type DegreeField string

const (
	DegreeFieldComputerScience                 DegreeField = "ComputerScience"
	DegreeFieldSoftwareEngineering             DegreeField = "SoftwareEngineering"
	DegreeFieldArtificialIntelligence          DegreeField = "ArtificialIntelligence"
	DegreeFieldCyberSecurity                   DegreeField = "CyberSecurity"
	DegreeFieldDataScience                     DegreeField = "DataScience"
	DegreeFieldEngineering                     DegreeField = "Engineering"
	DegreeFieldMechanicalEngineering           DegreeField = "MechanicalEngineering"
	DegreeFieldElectricalEngineering           DegreeField = "ElectricalEngineering"
	DegreeFieldCivilEngineering                DegreeField = "CivilEngineering"
	DegreeFieldChemicalEngineering             DegreeField = "ChemicalEngineering"
	DegreeFieldBusiness                        DegreeField = "Business"
	DegreeFieldMarketing                       DegreeField = "Marketing"
	DegreeFieldFinance                         DegreeField = "Finance"
	DegreeFieldManagement                      DegreeField = "Management"
	DegreeFieldEntrepreneurship                DegreeField = "Entrepreneurship"
	DegreeFieldMathematics                     DegreeField = "Mathematics"
	DegreeFieldPureMathematics                 DegreeField = "PureMathematics"
	DegreeFieldAppliedMathematics              DegreeField = "AppliedMathematics"
	DegreeFieldStatistics                      DegreeField = "Statistics"
	DegreeFieldMathematicalModeling            DegreeField = "MathematicalModeling"
	DegreeFieldPhysics                         DegreeField = "Physics"
	DegreeFieldTheoreticalPhysics              DegreeField = "TheoreticalPhysics"
	DegreeFieldQuantumPhysics                  DegreeField = "QuantumPhysics"
	DegreeFieldAstrophysicsAndAstronomy        DegreeField = "AstrophysicsAndAstronomy"
	DegreeFieldNuclearPhysics                  DegreeField = "NuclearPhysics"
	DegreeFieldChemistry                       DegreeField = "Chemistry"
	DegreeFieldOrganicChemistry                DegreeField = "OrganicChemistry"
	DegreeFieldInorganicChemistry              DegreeField = "InorganicChemistry"
	DegreeFieldPhysicalChemistry               DegreeField = "PhysicalChemistry"
	DegreeFieldBiochemistryAndMolecularBiology DegreeField = "BiochemistryAndMolecularBiology"
	DegreeFieldMedicine                        DegreeField = "Medicine"
	DegreeFieldGeneralMedicine                 DegreeField = "GeneralMedicine"
	DegreeFieldSurgery                         DegreeField = "Surgery"
	DegreeFieldPediatrics                      DegreeField = "Pediatrics"
	DegreeFieldPsychiatry                      DegreeField = "Psychiatry"
	DegreeFieldLaw                             DegreeField = "Law"
	DegreeFieldCorporateLaw                    DegreeField = "CorporateLaw"
	DegreeFieldCriminalLaw                     DegreeField = "CriminalLaw"
	DegreeFieldInternationalLaw                DegreeField = "InternationalLaw"
	DegreeFieldConstitutionalLaw               DegreeField = "ConstitutionalLaw"
	DegreeFieldSocialScience                   DegreeField = "SocialScience"
	DegreeFieldSociology                       DegreeField = "Sociology"
	DegreeFieldPoliticalScience                DegreeField = "PoliticalScience"
	DegreeFieldEconomics                       DegreeField = "Economics"
	DegreeFieldAnthropology                    DegreeField = "Anthropology"
	DegreeFieldHumanities                      DegreeField = "Humanities"
	DegreeFieldLiterature                      DegreeField = "Literature"
	DegreeFieldPhilosophy                      DegreeField = "Philosophy"
	DegreeFieldHistory                         DegreeField = "History"
	DegreeFieldLinguistics                     DegreeField = "Linguistics"
	DegreeFieldArt                             DegreeField = "Art"
	DegreeFieldPainting                        DegreeField = "Painting"
	DegreeFieldSculpture                       DegreeField = "Sculpture"
	DegreeFieldGraphicDesign                   DegreeField = "GraphicDesign"
	DegreeFieldPhotography                     DegreeField = "Photography"
	DegreeFieldMusic                           DegreeField = "Music"
	DegreeFieldComposition                     DegreeField = "Composition"
	DegreeFieldPerformance                     DegreeField = "Performance"
	DegreeFieldMusicTheory                     DegreeField = "MusicTheory"
	DegreeFieldConducting                      DegreeField = "Conducting"
	DegreeFieldSport                           DegreeField = "Sport"
	DegreeFieldSportsScience                   DegreeField = "SportsScience"
	DegreeFieldSportsManagement                DegreeField = "SportsManagement"
	DegreeFieldPhysicalEducation               DegreeField = "PhysicalEducation"
	DegreeFieldSportsTherapy                   DegreeField = "SportsTherapy"
	DegreeFieldEducation                       DegreeField = "Education"
	DegreeFieldElementaryEducation             DegreeField = "ElementaryEducation"
	DegreeFieldSecondaryEducation              DegreeField = "SecondaryEducation"
	DegreeFieldSpecialEducation                DegreeField = "SpecialEducation"
	DegreeFieldEducationalLeadership           DegreeField = "EducationalLeadership"
	DegreeFieldPsychology                      DegreeField = "Psychology"
	DegreeFieldClinicalPsychology              DegreeField = "ClinicalPsychology"
	DegreeFieldCognitivePsychology             DegreeField = "CognitivePsychology"
	DegreeFieldDevelopmentalPsychology         DegreeField = "DevelopmentalPsychology"
	DegreeFieldIndustrialPsychology            DegreeField = "IndustrialPsychology"
)

var AllDegreeField = []DegreeField{
	DegreeFieldComputerScience,
	DegreeFieldSoftwareEngineering,
	DegreeFieldArtificialIntelligence,
	DegreeFieldCyberSecurity,
	DegreeFieldDataScience,
	DegreeFieldEngineering,
	DegreeFieldMechanicalEngineering,
	DegreeFieldElectricalEngineering,
	DegreeFieldCivilEngineering,
	DegreeFieldChemicalEngineering,
	DegreeFieldBusiness,
	DegreeFieldMarketing,
	DegreeFieldFinance,
	DegreeFieldManagement,
	DegreeFieldEntrepreneurship,
	DegreeFieldMathematics,
	DegreeFieldPureMathematics,
	DegreeFieldAppliedMathematics,
	DegreeFieldStatistics,
	DegreeFieldMathematicalModeling,
	DegreeFieldPhysics,
	DegreeFieldTheoreticalPhysics,
	DegreeFieldQuantumPhysics,
	DegreeFieldAstrophysicsAndAstronomy,
	DegreeFieldNuclearPhysics,
	DegreeFieldChemistry,
	DegreeFieldOrganicChemistry,
	DegreeFieldInorganicChemistry,
	DegreeFieldPhysicalChemistry,
	DegreeFieldBiochemistryAndMolecularBiology,
	DegreeFieldMedicine,
	DegreeFieldGeneralMedicine,
	DegreeFieldSurgery,
	DegreeFieldPediatrics,
	DegreeFieldPsychiatry,
	DegreeFieldLaw,
	DegreeFieldCorporateLaw,
	DegreeFieldCriminalLaw,
	DegreeFieldInternationalLaw,
	DegreeFieldConstitutionalLaw,
	DegreeFieldSocialScience,
	DegreeFieldSociology,
	DegreeFieldPoliticalScience,
	DegreeFieldEconomics,
	DegreeFieldAnthropology,
	DegreeFieldHumanities,
	DegreeFieldLiterature,
	DegreeFieldPhilosophy,
	DegreeFieldHistory,
	DegreeFieldLinguistics,
	DegreeFieldArt,
	DegreeFieldPainting,
	DegreeFieldSculpture,
	DegreeFieldGraphicDesign,
	DegreeFieldPhotography,
	DegreeFieldMusic,
	DegreeFieldComposition,
	DegreeFieldPerformance,
	DegreeFieldMusicTheory,
	DegreeFieldConducting,
	DegreeFieldSport,
	DegreeFieldSportsScience,
	DegreeFieldSportsManagement,
	DegreeFieldPhysicalEducation,
	DegreeFieldSportsTherapy,
	DegreeFieldEducation,
	DegreeFieldElementaryEducation,
	DegreeFieldSecondaryEducation,
	DegreeFieldSpecialEducation,
	DegreeFieldEducationalLeadership,
	DegreeFieldPsychology,
	DegreeFieldClinicalPsychology,
	DegreeFieldCognitivePsychology,
	DegreeFieldDevelopmentalPsychology,
	DegreeFieldIndustrialPsychology,
}

func (e DegreeField) IsValid() bool {
	switch e {
	case DegreeFieldComputerScience, DegreeFieldSoftwareEngineering, DegreeFieldArtificialIntelligence, DegreeFieldCyberSecurity, DegreeFieldDataScience, DegreeFieldEngineering, DegreeFieldMechanicalEngineering, DegreeFieldElectricalEngineering, DegreeFieldCivilEngineering, DegreeFieldChemicalEngineering, DegreeFieldBusiness, DegreeFieldMarketing, DegreeFieldFinance, DegreeFieldManagement, DegreeFieldEntrepreneurship, DegreeFieldMathematics, DegreeFieldPureMathematics, DegreeFieldAppliedMathematics, DegreeFieldStatistics, DegreeFieldMathematicalModeling, DegreeFieldPhysics, DegreeFieldTheoreticalPhysics, DegreeFieldQuantumPhysics, DegreeFieldAstrophysicsAndAstronomy, DegreeFieldNuclearPhysics, DegreeFieldChemistry, DegreeFieldOrganicChemistry, DegreeFieldInorganicChemistry, DegreeFieldPhysicalChemistry, DegreeFieldBiochemistryAndMolecularBiology, DegreeFieldMedicine, DegreeFieldGeneralMedicine, DegreeFieldSurgery, DegreeFieldPediatrics, DegreeFieldPsychiatry, DegreeFieldLaw, DegreeFieldCorporateLaw, DegreeFieldCriminalLaw, DegreeFieldInternationalLaw, DegreeFieldConstitutionalLaw, DegreeFieldSocialScience, DegreeFieldSociology, DegreeFieldPoliticalScience, DegreeFieldEconomics, DegreeFieldAnthropology, DegreeFieldHumanities, DegreeFieldLiterature, DegreeFieldPhilosophy, DegreeFieldHistory, DegreeFieldLinguistics, DegreeFieldArt, DegreeFieldPainting, DegreeFieldSculpture, DegreeFieldGraphicDesign, DegreeFieldPhotography, DegreeFieldMusic, DegreeFieldComposition, DegreeFieldPerformance, DegreeFieldMusicTheory, DegreeFieldConducting, DegreeFieldSport, DegreeFieldSportsScience, DegreeFieldSportsManagement, DegreeFieldPhysicalEducation, DegreeFieldSportsTherapy, DegreeFieldEducation, DegreeFieldElementaryEducation, DegreeFieldSecondaryEducation, DegreeFieldSpecialEducation, DegreeFieldEducationalLeadership, DegreeFieldPsychology, DegreeFieldClinicalPsychology, DegreeFieldCognitivePsychology, DegreeFieldDevelopmentalPsychology, DegreeFieldIndustrialPsychology:
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
