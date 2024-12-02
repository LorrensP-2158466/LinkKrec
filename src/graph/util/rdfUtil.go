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

func MapRdfVacancyToGQL(vacancy map[string]rdf.Term) (*model.Vacancy, error) {
	vacancyObj, err := MapPrimitiveBindingsToStruct[model.Vacancy](vacancy)
	if err != nil {
		return nil, err
	}

	startDate := vacancy["startDate"].String()
	vacancyObj.StartDate = &startDate

	endDate := vacancy["endDate"].String()
	vacancyObj.EndDate = &endDate

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
