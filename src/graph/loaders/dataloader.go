package loaders

import (
	"LinkKrec/graph/model"
	"LinkKrec/graph/util"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/knakk/sparql"
	"github.com/vikstrous/dataloadgen"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

type userReader struct {
	Repo *sparql.Repo
}

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

// Loaders wrap your data loaders to inject via middleware
type Loaders struct {
	UserLoader *dataloadgen.Loader[string, *model.User]
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders(conn *sparql.Repo) *Loaders {
	// define the data loader
	ur := &userReader{Repo: conn}
	return &Loaders{
		UserLoader: dataloadgen.NewLoader(ur.getUsers, dataloadgen.WithWait(time.Millisecond)),
	}
}

// Middleware injects data loaders into the context
func Middleware(conn *sparql.Repo, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loader := NewLoaders(conn)
		r = r.WithContext(context.WithValue(r.Context(), loadersKey, loader))
		next.ServeHTTP(w, r)
	})
}

// For returns the dataloader for a given context
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

// GetUser returns single user by id efficiently
func GetUser(ctx context.Context, userID string) (*model.User, error) {
	loaders := For(ctx)
	return loaders.UserLoader.Load(ctx, userID)
}

// GetUsers returns many users by ids efficiently
func GetUsers(ctx context.Context, userIDs []string) ([]*model.User, error) {
	loaders := For(ctx)
	return loaders.UserLoader.LoadAll(ctx, userIDs)
}
