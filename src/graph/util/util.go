package util

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/knakk/rdf"
)

func MapPrimitiveBindingsToStruct[T any](bindings map[string]rdf.Term) (T, error) {
	// Check if T is a pointer type
	var zeroValue T
	isPtr := reflect.TypeOf(zeroValue).Kind() == reflect.Ptr

	// Get the underlying type (dereference if it's a pointer)
	var structType reflect.Type
	if isPtr {
		structType = reflect.TypeOf(zeroValue).Elem()
	} else {
		structType = reflect.TypeOf(zeroValue)
	}

	// Create a new value of the struct type
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
			} else if field.Type.Kind() == reflect.Bool {
				v, _ := strconv.ParseBool(bindingVal.String())
				structValue.Field(i).Set(reflect.ValueOf(v))
			} else if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.String {
				var sliceValues []string
				for _, v := range strings.Split(bindingVal.String(), ", ") {
					sliceValues = append(sliceValues, v)
				}
				structValue.Field(i).Set(reflect.ValueOf(sliceValues))
			}
		}
	}

	// Convert the result to the requested type
	if isPtr {
		result := structValue.Addr().Interface().(T)
		return result, nil
	}
	result := structValue.Interface().(T)
	return result, nil
}

func MapPrimitiveBindingsToStructArray[T any](solutions []map[string]rdf.Term) ([]T, error) {
	var arr = make([]T, 0)
	for _, skill := range solutions {
		mapped, err := MapPrimitiveBindingsToStruct[T](skill)
		if err != nil {
			return nil, err
		}
		arr = append(arr, mapped)
	}

	return arr, nil
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

func createSkillsQuery(skills []string) string {
	if len(skills) == 0 {
		return "."
	}

	var skillParts []string
	for _, skill := range skills {
		skillParts = append(skillParts, fmt.Sprintf("lr:requiredSkill lr:%s", skill))
	}
	return strings.Join(skillParts, " ; ") + " ."
}
