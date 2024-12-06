package util

import (
	"LinkKrec/graph/model"
	"fmt"
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
	if user["connections"] != nil {
		for _, con := range strings.Split(user["connections"].String(), ", ") {
			connections = append(connections, &model.User{ID: con})
		}
	}
	userObj.Connections = connections

	var educations = make([]*model.EducationEntry, 0)
	if user["educations"] != nil {
		for _, edu := range strings.Split(user["educations"].String(), ", ") {
			educations = append(educations, &model.EducationEntry{ID: edu})
		}
	}
	userObj.Education = educations

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

func MapRdfNotificationToGQL(notification map[string]rdf.Term) (*model.Notification, error) {
	notificationObj, err := MapPrimitiveBindingsToStruct[model.Notification](notification)
	if err != nil {
		return nil, err
	}

	notificationObj.ForUser = &model.User{ID: notification["forUserId"].String()}
	startDate := notification["createdAt"].String()
	notificationObj.CreatedAt = &startDate

	return &notificationObj, nil
}

func MapRdfConnectionRequestToGQL(connectionRequest map[string]rdf.Term) (*model.ConnectionRequest, error) {
	connectionRequestObj, err := MapPrimitiveBindingsToStruct[model.ConnectionRequest](connectionRequest)
	if err != nil {
		return nil, err
	}

	connectionRequestObj.FromUser = &model.User{ID: connectionRequest["fromUserId"].String()}
	connectionRequestObj.ConnectedToUser = &model.User{ID: connectionRequest["connectedToUserId"].String()}

	return &connectionRequestObj, nil
}

func MapRdfEducationEntryToGQL(educationEntry map[string]rdf.Term) (*model.EducationEntry, error) {
	educationEntryObj, err := MapPrimitiveBindingsToStruct[model.EducationEntry](educationEntry)
	if err != nil {
		return nil, err
	}
	fmt.Println(educationEntry)
	educationType := educationEntry["degree"].String()
	var degree model.DegreeType
	for _, d := range model.AllDegreeType {
		if d.String() == educationType {
			degree = d
			break
		}
	}
	educationEntryObj.Degree = degree

	educationField := educationEntry["field"].String()
	var field model.DegreeField
	for _, f := range model.AllDegreeField {
		if f.String() == educationField {
			field = f
			break
		}
	}
	educationEntryObj.Field = field

	return &educationEntryObj, nil
}

// func MapRdfExperienceEntryToGQL(experienceEntry map[string]rdf.Term) (*model.ExperienceEntry, error) {
// 	fmt.Println("experienceEntry: ", experienceEntry)
// 	experienceEntryObj, err := MapPrimitiveBindingsToStruct[model.ExperienceEntry](experienceEntry)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println("primitive experienceEntryObj: ", experienceEntryObj)

// 	experienceType := (experienceEntry["experienceType"].String())
// 	var experience model.ExperienceType
// 	for _, e := range model.AllExperienceType {
// 		if e.String() == experienceType {
// 			experience = e
// 			break
// 		}
// 	}
// 	experienceEntryObj.ExperienceType = experience

// 	startDate := experienceEntry["startDate"].String()
// 	experienceEntryObj.StartDate = &startDate

// 	endDate := experienceEntry["endDate"].String()
// 	experienceEntryObj.EndDate = &endDate

// 	return &experienceEntryObj, nil
// }
