package util

import (
	"LinkKrec/graph/model"
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
