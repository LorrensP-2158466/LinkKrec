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
	// zijn de optionals echt nodig hier?
	q := fmt.Sprintf(`
		PREFIX lr: <http://linkrec.example.org/schema#>
		PREFIX schema: <http://schema.org/>
		PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>

		SELECT ?name ?id ?email ?isEmployer ?location ?lookingForOpportunities
			(GROUP_CONCAT(DISTINCT ?skill; separator=", ") AS ?skills)
			(GROUP_CONCAT(DISTINCT ?connectionName; separator=", ") AS ?connections)
			(GROUP_CONCAT(DISTINCT ?educationEntry; separator=", ") AS ?educations)
			(GROUP_CONCAT(DISTINCT ?experienceEntry; separator=", ") AS ?experiences)
		WHERE {
		?user a lr:User ;
				lr:Id ?id ;
				lr:hasName ?name ;
				lr:hasEmail ?email ;
				lr:isEmployer ?isEmployer ;
				lr:hasLocation ?location ;
				lr:isLookingForOpportunities ?isLookingForOpportunities ;
				lr:hasSkill ?skill .
		BIND(?isLookingForOpportunities AS ?lookingForOpportunities)

		OPTIONAL {
			?user lr:hasConnection ?connection .
			?connection lr:Id ?connectionName .
		}
		OPTIONAL {
			?user lr:hasEducation ?education .
			?education lr:Id ?educationEntry .
		}
		OPTIONAL {
			?user lr:hasExperience ?experience .
			?experience lr:Id ?experienceEntry .
		}

		FILTER(%s)
		}
		GROUP BY ?name ?id ?email ?isEmployer ?location ?lookingForOpportunities
	`, filter)
	res, err := u.Repo.Query(q)
	if err != nil {
		return nil, []error{err}
	}

	fmt.Println("res: ", res)

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

// getVacancies implements a batch function that can retrieve many vacancies by ID,
// for use in a dataloader
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

		SELECT ?id ?title ?description ?location ?postedById ?startDate ?endDate ?status ?education (GROUP_CONCAT(DISTINCT ?experienceType; separator=", ") AS ?experienceTypes) (GROUP_CONCAT(DISTINCT ?experienceDuration; separator=", ") AS ?experienceDurations)
		WHERE {
			?vacancy a lr:Vacancy ;
			lr:Id ?id ;
			lr:vacancyTitle ?title ;
			lr:vacancyDescription ?description ;
			lr:vacancyLocation ?location ;
			lr:postedBy ?postedBy ;
			lr:vacancyStartDate ?startDate ;
			lr:vacancyEndDate ?endDate ;
			lr:vacancyStatus ?status ;
			lr:requiredEducation ?education ;
			lr:requiredExperienceType ?experienceType ;
			lr:requiredExperienceDuration ?experienceDuration .
			?postedBy lr:Id ?postedById .

		FILTER(%s)
		}
		GROUP BY ?id ?title ?description ?location ?postedById ?startDate ?endDate ?status ?education
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
	// fill return array with empty objects so the lengths match
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

func (u *DataBase) getEmployers(ctx context.Context, employerIDs []string) ([]*model.Employer, []error) {
	var ids []string
	for _, id := range employerIDs {
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
		?employer a lr:Employer ;
			lr:Id ?id ;
			lr:employerName ?name ;
			lr:employerEmail ?email ;
			lr:employerLocation ?location ;
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

	employers := make([]*model.Employer, len(employerIDs))
	errs := make([]error, len(employerIDs))

	var foundEmployers = make(map[string]*model.Employer)
	for _, m := range res.Solutions() {
		employer, err := util.MapRdfEmployerToGQL(m)
		if err != nil {
			return nil, []error{err}
		}
		foundEmployers[employer.ID] = employer
	}
	// fill return array with empty objects so the lengths match
	for i, id := range employerIDs {
		if employer, found := foundEmployers[id]; found {
			employers[i] = employer
			errs[i] = nil
		} else {
			employers[i] = &model.Employer{ID: id}
			errs[i] = fmt.Errorf("employer not found for ID: %s", id)
		}
	}
	return employers, errs
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

		SELECT ?id ?institution ?info ?degree ?field
		WHERE {
			?education a lr:EducationEntry ;
			lr:Id ?id ;
			lr:institutionName ?institution ;
			lr:institutionInfo ?info ;
			lr:degreeType ?degree ;
			lr:degreeField ?field .

		FILTER(%s)
		}
		GROUP BY ?id ?institution ?info ?degree ?field
	`, filter)
	res, err := u.Repo.Query(q)
	if err != nil {
		return nil, []error{err}
	}
	fmt.Println("res: ", res)

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

func (u *DataBase) getExperienceEntries(ctx context.Context, experienceEntryIDs []string) ([]*model.ExperienceEntry, []error) {
	var ids []string
	for _, id := range experienceEntryIDs {
		s := fmt.Sprintf("?id = \"%s\"", id)
		ids = append(ids, s)
	}
	filter := strings.Join(ids, " || ")
	q := fmt.Sprintf(`
		PREFIX lr: <http://linkrec.example.org/schema#>
		PREFIX schema: <http://schema.org/>
		PREFIX rdfs: <http://www.w3.org/2000/01/rdf-schema#>

		SELECT ?id ?title ?description ?startDate ?endDate ?experienceType
		WHERE {
			?experience a lr:ExperienceEntry ;
			lr:Id ?id ;
			lr:experienceTitle ?title ;
			lr:experienceDescription ?description ;
			lr:experienceType ?experienceType ;
			lr:experienceStartDate ?startDate ;
			lr:experienceEndDate ?endDate .

		FILTER(%s)
		}
	`, filter)
	res, err := u.Repo.Query(q)
	if err != nil {
		return nil, []error{err}
	}

	experienceEntries := make([]*model.ExperienceEntry, len(experienceEntryIDs))
	errs := make([]error, len(experienceEntryIDs))

	var foundExperienceEntries = make(map[string]*model.ExperienceEntry)
	for _, m := range res.Solutions() {
		experienceEntry, err := util.MapRdfExperienceEntryToGQL(m)
		if err != nil {
			return nil, []error{err}
		}
		foundExperienceEntries[experienceEntry.ID] = experienceEntry
	}
	// fill return array with empty objects so the lengths match
	for i, id := range experienceEntryIDs {
		if experienceEntry, found := foundExperienceEntries[id]; found {
			experienceEntries[i] = experienceEntry
			errs[i] = nil
		} else {
			experienceEntries[i] = &model.ExperienceEntry{ID: id}
			errs[i] = fmt.Errorf("experienceEntry not found for ID: %s", id)
		}
	}
	return experienceEntries, errs
}

func (u *DataBase) getConnectionRequests(ctx context.Context, connectionRequestIDs []string) ([]*model.ConnectionRequest, []error) {
	var ids []string
	for _, id := range connectionRequestIDs {
		s := fmt.Sprintf("?id = \"%s\"", id)
		ids = append(ids, s)
	}
	filter := strings.Join(ids, " || ")
	fmt.Println("filter: ", filter)
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
			lr:status ?status .
			?fromUser lr:Id ?fromUserId .
			?connectedToUser lr:Id ?connectedToUserId .

		FILTER(%s)
		}
	`, filter)
	fmt.Println("q: ", q)
	res, err := u.Repo.Query(q)
	fmt.Println("res: ", res)
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
