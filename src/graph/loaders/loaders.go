package loaders

import (
	"LinkKrec/graph/model"
	"LinkKrec/graph/util"
	"context"
	"fmt"
	"strings"
)

// getUsers implements a batch function that can retrieve many users by ID,
// for use in a dataloader
func (u *DataBase) getUsers(ctx context.Context, userIDs []string) ([]*model.User, []error) {
	fmt.Println("qmslghmlqsdhgmqsnhjg")
	var ids []string
	for _, id := range userIDs {
		s := fmt.Sprintf("?id = \"%s\"", id)
		ids = append(ids, s)
	}
	filter := strings.Join(ids, " || ")

	q := fmt.Sprintf(`
		PREFIX lr: <http://linkrec.example.org/schema#>
		PREFIX schema: <http://schema.org/>
		PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
		PREFIX skos: <http://www.w3.org/2004/02/skos/core#>
		PREFIX foaf: <http://xmlns.com/foaf/0.1/>
		PREFIX esco: <http://data.europa.eu/esco/model#>
		PREFIX esco_skill: <http://data.europa.eu/esco/skill/>

		SELECT ?id ?name ?email ?locationId ?lookingForOpportunities ?isProfileComplete
			(GROUP_CONCAT(DISTINCT ?connectionName; separator=", ") AS ?connections)
			(GROUP_CONCAT(DISTINCT ?educationEntry; separator=", ") AS ?educations)
			(GROUP_CONCAT(DISTINCT ?companyId; separator=", ") AS ?companies)
			(GROUP_CONCAT(CONCAT(STRAFTER(STR(?escoSkill), STR(esco_skill:)), "|", ?skill); separator=", ") as ?skillIdsAndLabels)
		WHERE {
			?user a lr:User ;
				lr:Id ?id ;
				foaf:name ?name ;
				foaf:mbox ?email ;
				lr:isLookingForOpportunities ?isLookingForOpportunities ;
				lr:isProfileComplete ?isProfileComplete ;
			BIND(?isLookingForOpportunities AS ?lookingForOpportunities)

		OPTIONAL {
			?user lr:hasSkill ?escoSkill .
			?escoSkill skos:prefLabel ?skill .
			FILTER(LANG(?skill) = "en")
		}
		OPTIONAL {
			?user lr:hasConnection ?connection .
			?connection lr:Id ?connectionName .
		}
		OPTIONAL {
			?user lr:hasEducation ?education .
			?education lr:Id ?educationEntry .
		}
		OPTIONAL {
			?user foaf:based_near ?location .
			?location lr:Id ?locationId .
		}
		OPTIONAL {
			?user lr:hasCompany ?company.
			?company lr:Id ?companyId .
		}

		FILTER(%s)
		}
		GROUP BY ?id ?name ?email ?locationId ?lookingForOpportunities ?isProfileComplete
	`, filter)
	res, err := u.Repo.Query(q)
	fmt.Println(res, err)
	if err != nil {
		return nil, []error{err}
	}
	users := make([]*model.User, len(userIDs))
	errs := make([]error, len(userIDs))
	fmt.Println(res.Bindings())

	var foundUsers = make(map[string]*model.User)
	for _, m := range res.Solutions() {
		user, err := util.MapRdfUserToGQL(m)
		if err != nil {
			return nil, []error{err}
		}
		foundUsers[user.ID] = user
	}
	fmt.Println(foundUsers)
	// fill return array with empty objects so the lengths match
	for i, id := range userIDs {
		if user, found := foundUsers[id]; found {
			users[i] = user
			errs[i] = nil
		} else {
			users[i] = &model.User{ID: id}
			errs[i] = fmt.Errorf("user not found for ID: %s", id)
		}
	}
	return users, errs
}

func (u *DataBase) getVacancies(ctx context.Context, vacancyIDs []string) ([]*model.Vacancy, []error) {
	var ids []string
	for _, id := range vacancyIDs {
		s := fmt.Sprintf("?id = \"%s\"", id)
		ids = append(ids, s)
	}
	filter := strings.Join(ids, " || ")
	q := fmt.Sprintf(`
		PREFIX lr: <http://linkrec.example.org/schema#>
		PREFIX schema: <http://schema.org/>
		PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>
		PREFIX esco_skill: <http://data.europa.eu/esco/skill/>
		PREFIX foaf: <http://xmlns.com/foaf/0.1/>
		PREFIX skos: <http://www.w3.org/2004/02/skos/core#>

		SELECT ?id ?title ?description ?locationId ?postedById ?startDate ?endDate ?status ?experienceDuration 
			(STRAFTER(STR(?degField), "#") AS ?degreeField) 
			(STRAFTER(STR(?degType), "#") AS ?degreeType)
			(GROUP_CONCAT(CONCAT(STRAFTER(STR(?escoSkill), STR(esco_skill:)), "|", ?skill); separator=", ") as ?skillIdsAndLabels)
		WHERE {
			?vacancy a lr:Vacancy ;
				lr:Id ?id ;
				lr:vacancyTitle ?title ;
				lr:vacancyDescription ?description ;
				foaf:based_near ?location ;
				lr:postedBy ?postedBy ;
				lr:vacancyStartDate ?startDate ;
				lr:vacancyEndDate ?endDate ;
				lr:vacancyStatus ?status .
				?postedBy lr:Id ?postedById .
				?location lr:Id ?locationId

			OPTIONAL { 
				?vacancy lr:requiredSkill ?escoSkill .
				?escoSkill skos:prefLabel ?skill .
				FILTER(LANG(?skill) = "en")
			}
			OPTIONAL { 
				?vacancy lr:requiredDegreeType ?degType .
			}
			OPTIONAL { 
				?vacancy lr:requiredDegreeField ?degField .
			}
			OPTIONAL { 
				?vacancy lr:requiredExperienceDuration ?experienceDuration .
			}
			FILTER(%s)
		}
		GROUP BY ?id ?title ?description ?locationId ?postedById ?startDate ?endDate ?status ?degField ?degType ?experienceDuration
	`, filter)
	res, err := u.Repo.Query(q)

	if err != nil {
		return nil, []error{err}
	}

	vacancies := make([]*model.Vacancy, len(vacancyIDs))
	errs := make([]error, len(vacancyIDs))

	var foundVacancies = make(map[string]*model.Vacancy)
	for _, m := range res.Solutions() {
		vacancy, err := util.MapRdfVacancyToGQL(m)
		if err != nil {
			return nil, []error{err}
		}
		foundVacancies[vacancy.ID] = vacancy
	}

	for i, id := range vacancyIDs {
		if vacancy, found := foundVacancies[id]; found {
			vacancies[i] = vacancy
			errs[i] = nil
		} else {
			vacancies[i] = &model.Vacancy{ID: id}
			errs[i] = fmt.Errorf("vacancy not found for ID: %s", id)
		}
	}
	return vacancies, errs
}

func (u *DataBase) getCompanies(ctx context.Context, companyIDs []string) ([]*model.Company, []error) {
	var ids []string
	for _, id := range companyIDs {
		s := fmt.Sprintf("?id = \"%s\"", id)
		ids = append(ids, s)
	}
	filter := strings.Join(ids, " || ")
	q := fmt.Sprintf(`
		PREFIX lr: <http://linkrec.example.org/schema#>
		PREFIX rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#>
		PREFIX list: <http://jena.hpl.hp.com/ARQ/list#>

		SELECT ?id ?name ?email ?locationId 
			(GROUP_CONCAT(DISTINCT ?vacancyId; separator=", ") AS ?vacancies) 
			(GROUP_CONCAT(DISTINCT ?employeeId; separator=", ") AS ?employees)   
		WHERE {
			?company a lr:Company ;
				lr:Id ?id ;
				lr:companyName ?name ;
				lr:companyEmail ?email ;
				lr:companyLocation ?location .

			?location lr:Id ?locationId

			OPTIONAL{
				?company lr:hasEmployee ?employee .
				?employee lr:Id ?employeeId .
			}
			OPTIONAL{
				?company lr:hasVacancy ?vacancy .
				?vacancy lr:Id ?vacancyId .
			}

			FILTER(%s)
		}
		GROUP BY ?id ?name ?email ?location ?locationId
	`, filter)
	res, err := u.Repo.Query(q)
	if err != nil {
		return nil, []error{err}
	}

	companies := make([]*model.Company, len(companyIDs))
	errs := make([]error, len(companyIDs))

	var foundEmployers = make(map[string]*model.Company)
	for _, m := range res.Solutions() {
		company, err := util.MapRdfCompanyToGQL(m)
		if err != nil {
			return nil, []error{err}
		}
		foundEmployers[company.ID] = company
	}
	// fill return array with empty objects so the lengths match
	for i, id := range companyIDs {
		if company, found := foundEmployers[id]; found {
			companies[i] = company
			errs[i] = nil
		} else {
			companies[i] = &model.Company{ID: id}
			errs[i] = fmt.Errorf("company not found for ID: %s", id)
		}
	}
	return companies, errs
}

func (u *DataBase) getEducationEntries(ctx context.Context, educationEntryIDs []string) ([]*model.EducationEntry, []error) {
	var ids []string
	for _, id := range educationEntryIDs {
		s := fmt.Sprintf("?id = \"%s\"", id)
		ids = append(ids, s)
	}
	filter := strings.Join(ids, " || ")
	q := fmt.Sprintf(`
		PREFIX lr: <http://linkrec.example.org/schema#>
		PREFIX schema: <http://schema.org/>
		PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>

		SELECT ?id ?institution ?info ?startDate ?endDate
			(STRAFTER(STR(?degField), "#") AS ?field) 
			(STRAFTER(STR(?degType), "#") AS ?degree)
		WHERE {
			?education a lr:EducationEntry ;
				lr:Id ?id ;
				lr:institutionName ?institution ;
    			lr:startDate ?startDate ;
    			lr:endDate ?endDate ;
                lr:degreeType ?degType ;
				lr:degreeField ?degField .
			FILTER(%s)
		}
	`, filter)
	res, err := u.Repo.Query(q)
	if err != nil {
		return nil, []error{err}
	}

	educationEntries := make([]*model.EducationEntry, len(educationEntryIDs))
	errs := make([]error, len(educationEntryIDs))

	var foundEducationEntries = make(map[string]*model.EducationEntry)
	for _, m := range res.Solutions() {
		educationEntry, err := util.MapRdfEducationEntryToGQL(m)
		if err != nil {
			return nil, []error{err}
		}
		foundEducationEntries[educationEntry.ID] = educationEntry
	}
	// fill return array with empty objects so the lengths match
	for i, id := range educationEntryIDs {
		if educationEntry, found := foundEducationEntries[id]; found {
			educationEntries[i] = educationEntry
			errs[i] = nil
		} else {
			educationEntries[i] = &model.EducationEntry{ID: id}
			errs[i] = fmt.Errorf("educationEntry not found for ID: %s", id)
		}
	}
	return educationEntries, errs
}

func (u *DataBase) getConnectionRequests(ctx context.Context, connectionRequestIDs []string) ([]*model.ConnectionRequest, []error) {
	var ids []string
	for _, id := range connectionRequestIDs {
		s := fmt.Sprintf("?id = \"%s\"", id)
		ids = append(ids, s)
	}
	filter := strings.Join(ids, " || ")
	q := fmt.Sprintf(`
		PREFIX lr: <http://linkrec.example.org/schema#>
		PREFIX schema: <http://schema.org/>
		PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>

		SELECT ?id ?fromUserId ?connectedToUserId ?status
		WHERE {
			?connectionRequest a lr:ConnectionRequest ;
			lr:Id ?id ;
			lr:fromUser ?fromUser ;
			lr:connectedToUser ?connectedToUser ;
			lr:requestStatus ?status .
			?fromUser lr:Id ?fromUserId .
			?connectedToUser lr:Id ?connectedToUserId .

		FILTER(%s)
		}
	`, filter)
	res, err := u.Repo.Query(q)
	if err != nil {
		return nil, []error{err}
	}

	connectionRequests := make([]*model.ConnectionRequest, len(connectionRequestIDs))
	errs := make([]error, len(connectionRequestIDs))

	var foundConnectionRequests = make(map[string]*model.ConnectionRequest)
	for _, m := range res.Solutions() {
		connectionRequest, err := util.MapRdfConnectionRequestToGQL(m)
		if err != nil {
			return nil, []error{err}
		}
		foundConnectionRequests[connectionRequest.ID] = connectionRequest
	}
	// fill return array with empty objects so the lengths match
	for i, id := range connectionRequestIDs {
		if connectionRequest, found := foundConnectionRequests[id]; found {
			connectionRequests[i] = connectionRequest
			errs[i] = nil
		} else {
			connectionRequests[i] = &model.ConnectionRequest{ID: id}
			errs[i] = fmt.Errorf("connectionRequest not found for ID: %s", id)
		}
	}
	return connectionRequests, errs
}

func (u *DataBase) getNotifications(ctx context.Context, notificationIDs []string) ([]*model.Notification, []error) {
	var ids []string
	for _, id := range notificationIDs {
		s := fmt.Sprintf("?id = \"%s\"", id)
		ids = append(ids, s)
	}
	filter := strings.Join(ids, " || ")
	q := fmt.Sprintf(`
		PREFIX lr: <http://linkrec.example.org/schema#>
		PREFIX schema: <http://schema.org/>
		PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>

		SELECT ?id ?title ?message ?userId ?createdAt
		WHERE {
			?notification a lr:Notification ;
			lr:Id ?id ;
			lr:notifcationTitle ?title ;
			lr:notificationMessage ?message ;
			lr:forUser ?user ;
			notificationCreatedAt ?createdAt ;
			?user lr:Id ?userId .

		FILTER(%s)
		}
	`, filter)
	res, err := u.Repo.Query(q)
	if err != nil {
		return nil, []error{err}
	}

	notifications := make([]*model.Notification, len(notificationIDs))
	errs := make([]error, len(notificationIDs))

	var foundNotifications = make(map[string]*model.Notification)
	for _, m := range res.Solutions() {
		notification, err := util.MapRdfNotificationToGQL(m)
		if err != nil {
			return nil, []error{err}
		}
		foundNotifications[notification.ID] = notification
	}
	// fill return array with empty objects so the lengths match
	for i, id := range notificationIDs {
		if notification, found := foundNotifications[id]; found {
			notifications[i] = notification
			errs[i] = nil
		} else {
			notifications[i] = &model.Notification{ID: id}
			errs[i] = fmt.Errorf("notification not found for ID: %s", id)
		}
	}
	return notifications, errs
}

func (u *DataBase) getLocations(ctx context.Context, locationIDs []string) ([]*model.Location, []error) {
	print("qsfqsgqimjoji")
	var ids []string
	for _, id := range locationIDs {
		s := fmt.Sprintf("?id = \"%s\"", id)
		ids = append(ids, s)
	}
	filter := strings.Join(ids, " || ")
	q := fmt.Sprintf(`
		PREFIX lr: <http://linkrec.example.org/schema#>
		PREFIX schema: <http://schema.org/>
		PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>

		SELECT ?id ?country ?city ?street ?houseNumber ?longitude ?latitude
			WHERE {
				?location a lr:Location ;
				lr:Id ?id ;
				lr:inCountry ?country ;
				lr:inCity ?city ;
				lr:inStreet ?street ;
				lr:houseNumber ?houseNumber ;

			FILTER(%s)
		}
	`, filter)
	res, err := u.Repo.Query(q)
	if err != nil {
		return nil, []error{err}
	}

	locations := make([]*model.Location, len(locationIDs))
	errs := make([]error, len(locationIDs))

	var foundLocations = make(map[string]*model.Location)
	for _, m := range res.Solutions() {
		location, err := util.MapRdfLocationToGQL(m)
		if err != nil {
			return nil, []error{err}
		}
		foundLocations[location.ID] = location
	}
	// fill return array with empty objects so the lengths match
	for i, id := range locationIDs {
		if location, found := foundLocations[id]; found {
			locations[i] = location
			errs[i] = nil
		} else {
			locations[i] = &model.Location{ID: id}
			errs[i] = fmt.Errorf("location not found for ID: %s", id)
		}
	}
	return locations, errs
}
