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
		user, err := util.MapPrimitiveBindingsToStruct[model.User](m)
		if err != nil {
			return nil, []error{err}
		}
		var connections = make([]*model.User, 0)
		for _, con := range strings.Split(m["connections"].String(), ", ") {
			connections = append(connections, &model.User{ID: con})
		}
		user.Connections = connections
		foundUsers[user.ID] = &user
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
