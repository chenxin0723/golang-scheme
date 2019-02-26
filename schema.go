package schema

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

type SchemaValidator struct {
	req             *http.Request
	formatFuncMap   map[string]func(in interface{}) (out interface{}, err error)
	validateFuncMap map[string]func(in interface{}) bool
}

type Config struct {
	Request         *http.Request
	FormatFuncMap   map[string]func(in interface{}) (out interface{}, err error)
	ValidateFuncMap map[string]func(in interface{}) bool
}

func NewSchemaValidator(Config Config) (SchemaValidator, error) {
	schemaValidator := SchemaValidator{}
	if Config.Request != nil {
		schemaValidator.req = Config.Request
	} else {
		return nil, errors.New("request is nil")
	}

	schemaValidator.formatFuncMap = Config.FormatFuncMap
	schemaValidator.validateFuncMap = Config.ValidateFuncMap
	return schemaValidator, nil
}

func (schemaValidator SchemaValidator) Encode(in interface{}) error {
	value := reflect.ValueOf(in)
	if value.Kind() != reflect.Struct {
		return errors.New(fmt.Sprintf("%s should be a struct kind", in))
	}

}
