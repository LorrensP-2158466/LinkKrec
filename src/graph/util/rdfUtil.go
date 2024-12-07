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

	var companies = make([]*model.Company, 0)
	if user["companies"] != nil {
		for _, comp := range strings.Split(user["companies"].String(), ", ") {
			companies = append(companies, &model.Company{ID: comp})
		}
	}
	userObj.Companies = companies
	return &userObj, nil
}

func MapRdfCompanyToGQL(company map[string]rdf.Term) (*model.Company, error) {
	companyObj, err := MapPrimitiveBindingsToStruct[model.Company](company)
	if err != nil {
		return nil, err
	}

	if company["location"] != nil {
		var location = company["location"].String()
		companyObj.Location = &location
	} else {
		companyObj.Location = nil
	}

	if company["vacancies"] != nil {
		var vacancies = make([]*model.Vacancy, 0)
		for _, vac := range strings.Split(company["vacancies"].String(), ", ") {
			vacancies = append(vacancies, &model.Vacancy{ID: vac})
		}
		companyObj.Vacancies = vacancies
	} else {
		companyObj.Vacancies = nil
	}

	if company["employees"] != nil {
		var employees = make([]*model.User, 0)
		for _, emp := range strings.Split(company["employees"].String(), ", ") {
			employees = append(employees, &model.User{ID: emp})
		}
		companyObj.Employees = employees
	} else {
		companyObj.Employees = nil
	}

	return &companyObj, nil
}

func MapRdfVacancyToGQL(vacancy map[string]rdf.Term) (*model.Vacancy, error) {
	vacancyObj, err := MapPrimitiveBindingsToStruct[model.Vacancy](vacancy)
	if err != nil {
		return nil, err
	}

	startDate := vacancy["startDate"].String()
	vacancyObj.StartDate = startDate

	endDate := vacancy["endDate"].String()
	vacancyObj.EndDate = endDate

	vacancyObj.PostedBy = &model.Company{ID: vacancy["postedById"].String()}

	degreeType := vacancy["degreeType"].String()
	var degreeTypeObj model.DegreeType
	for _, d := range model.AllDegreeType {
		if d.String() == degreeType {
			degreeTypeObj = d
			break
		}
	}
	vacancyObj.RequiredDegreeType = &degreeTypeObj

	degreeField := vacancy["degreeField"].String()
	var degreeFieldObj model.DegreeField
	for _, d := range model.AllDegreeField {
		if d.String() == degreeField {
			degreeFieldObj = d
			break
		}
	}
	vacancyObj.RequiredDegreeField = &degreeFieldObj

	experienceDuration, err := strconv.Atoi(vacancy["experienceDuration"].String())
	if err != nil {
		return nil, err
	}
	vacancyObj.RequiredExperienceDuration = &experienceDuration

	return &vacancyObj, nil
}

func MapRdfNotificationToGQL(notification map[string]rdf.Term) (*model.Notification, error) {
	notificationObj, err := MapPrimitiveBindingsToStruct[model.Notification](notification)
	if err != nil {
		return nil, err
	}

	notificationObj.ForUser = &model.User{ID: notification["forUserId"].String()}
	if notification["createdAt"] != nil {
		startDate := notification["createdAt"].String()
		notificationObj.CreatedAt = &startDate
	} else {
		notificationObj.CreatedAt = nil
	}

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
