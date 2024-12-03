package util

import (
	"LinkKrec/graph/model"
	"strconv"
	"strings"

	"github.com/knakk/rdf"
)

func MapRdfUserToGQL(user map[string]rdf.Term) (*model.User, error) {
	userObj, err := MapPrimitiveBindingsToStruct[model.User](user)
	if err != nil {
		return nil, err
	}
	var connections = make([]*model.User, 0)
	for _, con := range strings.Split(user["connections"].String(), ", ") {
		connections = append(connections, &model.User{ID: con})
	}
	userObj.Connections = connections
	return &userObj, nil
}

func MapRdfEmployerToGQL(employer map[string]rdf.Term) (*model.Employer, error) {
	employerObj, err := MapPrimitiveBindingsToStruct[model.Employer](employer)
	if err != nil {
		return nil, err
	}

	var location = employer["location"].String()
	employerObj.Location = &location

	var vacancies = make([]*model.Vacancy, 0)
	for _, vac := range strings.Split(employer["vacancies"].String(), ", ") {
		vacancies = append(vacancies, &model.Vacancy{ID: vac})
	}
	employerObj.Vacancies = vacancies

	var employees = make([]*model.User, 0)
	for _, emp := range strings.Split(employer["employees"].String(), ", ") {
		employees = append(employees, &model.User{ID: emp})
	}
	employerObj.Employees = employees

	return &employerObj, nil
}

func MapRdfVacancyToGQL(vacancy map[string]rdf.Term) (*model.Vacancy, error) {
	vacancyObj, err := MapPrimitiveBindingsToStruct[model.Vacancy](vacancy)
	if err != nil {
		return nil, err
	}

	startDate := vacancy["startDate"].String()
	vacancyObj.StartDate = &startDate

	endDate := vacancy["endDate"].String()
	vacancyObj.EndDate = &endDate

	vacancyObj.PostedBy = &model.Employer{ID: vacancy["postedById"].String()}

	education := (vacancy["education"].String())
	var degree model.DegreeType
	for _, d := range model.AllDegreeType {
		if d.String() == education {
			degree = d
			break
		}
	}
	vacancyObj.RequiredEducation = degree

	experienceType := vacancy["experienceTypes"].String()
	experienceTypes := strings.Split(experienceType, ", ")
	var experiences []model.ExperienceType
	for _, e := range experienceTypes {
		for _, f := range model.AllExperienceType {
			if f.String() == e {
				experiences = append(experiences, f)
				break
			}
		}
	}
	vacancyObj.RequiredExperiences = experiences

	experienceDuration := vacancy["experienceDurations"].String()
	experienceDurations := strings.Split(experienceDuration, ", ")
	var durations []int
	for _, d := range experienceDurations {
		duration, _ := strconv.Atoi(d)
		durations = append(durations, duration)
	}
	vacancyObj.RequiredExperienceDurations = durations

	return &vacancyObj, nil
}
