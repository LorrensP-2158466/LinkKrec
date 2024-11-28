package graph

import (
	"reflect"
	"strings"

	"github.com/knakk/rdf"
)

func MapStringBindingsToStruct[T any](bindings map[string][]rdf.Term) (T, error) {
	structType := reflect.TypeOf(*new(T))
	structValue := reflect.New(structType).Elem()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		bindingKey := strings.Split(field.Tag.Get("json"), ",")[0]
		bindingVal := bindings[bindingKey]

		if bindingVal != nil {
			if field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Ptr && field.Type.Elem().Elem().Kind() == reflect.String {
				var sliceValues []*string
				for _, v := range bindingVal {
					strPtr := new(string)
					*strPtr = v.String()
					sliceValues = append(sliceValues, strPtr)
				}
				structValue.Field(i).Set(reflect.ValueOf(sliceValues))
			} else if field.Type.Kind() == reflect.Slice {
				var sliceValues []string
				for _, v := range bindingVal {
					sliceValues = append(sliceValues, v.String())
				}
				structValue.Field(i).Set(reflect.ValueOf(sliceValues))
			} else {
				structValue.Field(i).Set(reflect.ValueOf(bindingVal[0].String()))
			}
		}
	}
	result := structValue.Interface().(T)
	return result, nil
}
