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
func (u *userReader) getUsers(ctx context.Context, userIDs []string) ([]*model.User, []error) {
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

		SELECT ?name ?id ?email ?isEmployer ?location ?isLookingForOpportunities
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

		OPTIONAL {
			?user lr:hasConnection ?connection .
			?connection lr:Id ?connectionName .
		}
		OPTIONAL {
			?user lr:hasEducation ?education .
			?education rdfs:label ?educationEntry .
		}
		OPTIONAL {
			?user lr:hasExperience ?experience .
			?experience rdfs:label ?experienceEntry .
		}

		FILTER(%s)
		}
		GROUP BY ?name ?id ?email ?isEmployer ?location ?isLookingForOpportunities
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

// getVacancies implements a batch function that can retrieve many vacancies by ID,
// for use in a dataloader
func (u *userReader) getVacancies(ctx context.Context, vacancyIDs []string) ([]*model.Vacancy, []error) {
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

func (u *userReader) getEmployers(ctx context.Context, employerIDs []string) ([]*model.Employer, []error) {
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
