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
	if user["locationId"] != nil {
		userObj.Location = &model.Location{ID: user["locationId"].String()}
	} else {
		userObj.Location = nil
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

	var skills = make([]*model.Skill, 0)
	if user["skillIdsAndLabels"] != nil {
		for _, skillIdAndLabel := range strings.Split(user["skillIdsAndLabels"].String(), ", ") {
			splitted := strings.SplitN(skillIdAndLabel, "|", 2)
			skills = append(skills, &model.Skill{ID: splitted[0], Label: splitted[1]})
		}
	}
	userObj.Skills = skills

	var expr = make([]*model.Experience, 0)
	if user["occupIdsLabelsDurs"] != nil {
		for _, idsLabelsDurs := range strings.Split(user["occupIdsLabelsDurs"].String(), ", ") {
			splitted := strings.Split(idsLabelsDurs, "|")
			dur, _ := strconv.Atoi(splitted[2])
			expr = append(expr, &model.Experience{ID: splitted[0], Label: splitted[1], DurationInMonths: dur})
		}
	}
	userObj.Experiences = expr

	return &userObj, nil
}

func MapRdfCompanyToGQL(company map[string]rdf.Term) (*model.Company, error) {
	companyObj, err := MapPrimitiveBindingsToStruct[model.Company](company)
	if err != nil {
		return nil, err
	}

	if company["locationId"] != nil {
		companyObj.Location = &model.Location{ID: company["locationId"].String()}
	} else {
		companyObj.Location = nil
	}

	var vacancies = make([]*model.Vacancy, 0)
	if company["vacancies"] != nil {
		for _, vac := range strings.Split(company["vacancies"].String(), ", ") {
			vacancies = append(vacancies, &model.Vacancy{ID: vac})
		}
	}
	companyObj.Vacancies = vacancies

	var employees = make([]*model.User, 0)
	if company["employees"] != nil {
		for _, emp := range strings.Split(company["employees"].String(), ", ") {
			employees = append(employees, &model.User{ID: emp})
		}
	}
	companyObj.Employees = employees

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

	LocObj := model.Location{
		ID: vacancy["locationId"].String(),
	}
	vacancyObj.Location = &LocObj

	var skills = make([]*model.Skill, 0)
	if vacancy["skillIdsAndLabels"] != nil {
		for _, skillIdAndLabel := range strings.Split(vacancy["skillIdsAndLabels"].String(), ", ") {
			splitted := strings.SplitN(skillIdAndLabel, "|", 2)
			skills = append(skills, &model.Skill{ID: splitted[0], Label: splitted[1]})
		}
	}
	vacancyObj.RequiredSkills = skills

	var expr = make([]*model.Experience, 0)
	if vacancy["occupIdsLabelsDurs"] != nil {
		for _, idsLabelsDurs := range strings.Split(vacancy["occupIdsLabelsDurs"].String(), ", ") {
			splitted := strings.SplitN(idsLabelsDurs, "|", 3)
			dur, _ := strconv.Atoi(splitted[2])
			expr = append(expr, &model.Experience{ID: splitted[0], Label: splitted[1], DurationInMonths: dur})
		}
	}
	vacancyObj.RequiredExperience = expr

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

	if connectionRequest["fromUserId"] != nil {
		connectionRequestObj.FromUser = &model.User{ID: connectionRequest["fromUserId"].String()}
	} else {
		connectionRequestObj.FromUser = nil
	}

	if connectionRequest["connectedToUserId"] != nil {
		connectionRequestObj.ConnectedToUser = &model.User{ID: connectionRequest["connectedToUserId"].String()}
	} else {
		connectionRequestObj.ConnectedToUser = nil
	}

	return &connectionRequestObj, nil
}

func MapRdfEducationEntryToGQL(educationEntry map[string]rdf.Term) (*model.EducationEntry, error) {
	educationEntryObj, err := MapPrimitiveBindingsToStruct[model.EducationEntry](educationEntry)
	if err != nil {
		return nil, err
	}

	if educationEntry["degree"] != nil {
		educationType := (educationEntry["degree"].String())
		var degree model.DegreeType
		for _, d := range model.AllDegreeType {
			if d.String() == educationType {
				degree = d
				break
			}
		}
		educationEntryObj.Degree = degree
	}

	if educationEntry["field"] != nil {
		educationField := (educationEntry["field"].String())
		var field model.DegreeField
		for _, f := range model.AllDegreeField {
			if f.String() == educationField {
				field = f
				break
			}
		}
		educationEntryObj.Field = field
	}

	return &educationEntryObj, nil
}

func MapRdfLocationToGQL(location map[string]rdf.Term) (*model.Location, error) {
	locationObj, err := MapPrimitiveBindingsToStruct[model.Location](location)
	if err != nil {
		return nil, err
	}
	return &locationObj, nil
}
