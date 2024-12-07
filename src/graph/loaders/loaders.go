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

		SELECT ?id ?name ?email ?location ?lookingForOpportunities
			(GROUP_CONCAT(DISTINCT ?skill; separator=", ") AS ?skills)
			(GROUP_CONCAT(DISTINCT ?connectionName; separator=", ") AS ?connections)
			(GROUP_CONCAT(DISTINCT ?educationEntry; separator=", ") AS ?educations)
			(GROUP_CONCAT(DISTINCT ?companyId; separator=", ") AS ?companies)
		WHERE {
		?user a lr:User ;
				lr:Id ?id ;
				lr:hasName ?name ;
				lr:hasEmail ?email ;
				lr:isLookingForOpportunities ?isLookingForOpportunities ;
		BIND(?isLookingForOpportunities AS ?lookingForOpportunities)

		OPTIONAL {
			?user lr:hasSkill ?escoSkill .
			?escoSkill skos:prefLabel ?skill .
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
			?user lr:hasLocation ?location .
			?location lr:Id ?locationEntry .
		}
		OPTIONAL {
			?user a lr:User ;
			lr:hasCompany ?company .
			?company lr:Id ?companyId .
		}

		FILTER(%s)
		FILTER(LANG(?skill) = "en")
		}
		GROUP BY ?id ?name ?email ?location ?lookingForOpportunities
	`, filter)
	res, err := u.Repo.Query(q)
	if err != nil {
		return nil, []error{err}
	}

	users := make([]*model.User, len(userIDs))
	errs := make([]error, len(userIDs))

	var foundUsers = make(map[string]*model.User)
	for _, m := range res.Solutions() {
		user, err := util.MapRdfUserToGQL(m)
		if err != nil {
			return nil, []error{err}
		}
		foundUsers[user.ID] = user
	}
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
		PREFIX esco_skill: <http://data.europa.eu/esco/Skill>
		PREFIX skos: <http://www.w3.org/2004/02/skos/core#>

		SELECT ?id ?title ?description ?location ?postedById ?startDate ?endDate ?status ?degreeType ?degreeField ?experienceDuration (GROUP_CONCAT(DISTINCT ?skill; separator=", ") AS ?skills)
		WHERE {
			?vacancy a lr:Vacancy ;
				lr:Id ?id ;
				lr:vacancyTitle ?title ;
				lr:vacancyDescription ?description ;
				lr:vacancyLocation ?location ;
				lr:postedBy ?postedBy ;
				lr:vacancyStartDate ?startDate ;
				lr:vacancyEndDate ?endDate ;
				lr:vacancyStatus ?status .
				?postedBy lr:Id ?postedById .

			OPTIONAL { 
				?vacancy lr:requiredSkill ?skill .
				?skill skos:prefLabel ?skillLabel .
				FILTER(LANG(?skillLabel) = "en")
			}
			OPTIONAL { 
				?vacancy lr:requiredDegreeType ?degreeType .
			}
			OPTIONAL { 
				?vacancy lr:requiredDegreeField ?degreeField .
			}
			OPTIONAL { 
				?vacancy lr:requiredExperienceDuration ?experienceDuration .
			}
			FILTER(%s)
		}
		GROUP BY ?id ?title ?description ?location ?postedById ?startDate ?endDate ?status ?degreeType ?degreeField ?experienceDuration

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

		SELECT ?id ?name ?email ?location (GROUP_CONCAT(DISTINCT ?vacancyId; separator=", ") AS ?vacancies) (GROUP_CONCAT(DISTINCT ?employeeId; separator=", ") AS ?employees)   
		WHERE {
		?company a lr:Company ;
			lr:Id ?id ;
			lr:companyName ?name ;
			lr:companyEmail ?email ;
			lr:companyLocation ?location ;
			lr:hasVacancy ?vacancy ;
			lr:hasEmployee ?employee .
			?vacancy lr:Id ?vacancyId .
			?employee lr:Id ?employeeId .

		FILTER(%s)
		}
		GROUP BY ?id ?name ?email ?location
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

		SELECT ?id ?institution ?info 
			(STRAFTER(STR(?degField), "#") AS ?field) 
			(STRAFTER(STR(?degType), "#") AS ?degree)
		WHERE {
			?education a lr:EducationEntry ;
			lr:Id ?id ;
			lr:institutionName ?institution ;
			lr:institutionInfo ?info ;
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
