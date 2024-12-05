package util

import (
	"encoding/gob"
	"reflect"
	"strconv"
	"strings"

	"github.com/knakk/rdf"
)

const (
	SessionInfoKey = "sessInfo"
	GinContextKey  = "GinContextKey"
	QueryRepoKey   = "queryrepo"
	UpdateRepoKey  = "updateRepo"
)

type UserSessionInfo struct {
	// convenience to quickly determine of this user has a completed account
	IsComplete bool `json:"ProfileCompleted"`
	IsUser     bool
	Email      string `json:"email"`
	Id         string `json:"id"`
	Cookie     string
	// TODO: More?
}

func init() {
	gob.Register(UserSessionInfo{}) // Register the type for serialization
}

func MapPrimitiveBindingsToStruct[T any](bindings map[string]rdf.Term) (T, error) {
	structType := reflect.TypeOf(*new(T))
	structValue := reflect.New(structType).Elem()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		bindingKey := strings.Split(field.Tag.Get("json"), ",")[0]
		bindingVal := bindings[bindingKey]

		if bindingVal != nil {
			if isEnumType(field.Type) {
				continue
			}
			if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Ptr && field.Type.Elem().Elem().Kind() == reflect.String {
				var sliceValues []*string
				for _, v := range strings.Split(bindingVal.String(), ", ") {
					strPtr := new(string)
					*strPtr = v
					sliceValues = append(sliceValues, strPtr)
				}
				structValue.Field(i).Set(reflect.ValueOf(sliceValues))
			} else if field.Type.Kind() == reflect.String {
				structValue.Field(i).Set(reflect.ValueOf(bindingVal.String()))
			} else if field.Type.Kind() == reflect.Ptr {
				if field.Type.Elem().Kind() == reflect.Bool {
					v, _ := strconv.ParseBool(bindingVal.String())
					structValue.Field(i).Set(reflect.ValueOf(&v))
				} else if field.Type.Elem().Kind() == reflect.String {
					strPtr := new(string)
					*strPtr = bindingVal.String()
					structValue.Field(i).Set(reflect.ValueOf(strPtr))
				}
			}
		}
	}
	result := structValue.Interface().(T)
	return result, nil
}

func isEnumType(fieldType reflect.Type) bool {
	return fieldType.PkgPath() != "" && (fieldType.Kind() == reflect.Int || fieldType.Kind() == reflect.String)
}

// map array to new array using a function
func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}
