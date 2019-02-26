package schema

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

// type People struct {
// 	Name string `field: "name" format: "" validator "" required`
// }

type SchemaValidator struct {
	req              *http.Request
	formatFuncMap    map[string]func(in string) (out interface{}, err error)
	validatorFuncMap map[string]func(in string) (bool, error)
}

type Config struct {
	Request          *http.Request
	FormatFuncMap    map[string]func(in string) (out interface{}, err error)
	ValidatorFuncMap map[string]func(in string) (bool, error)
}

func NewSchemaValidator(Config Config) (SchemaValidator, error) {
	schemaValidator := SchemaValidator{}
	if Config.Request != nil {
		schemaValidator.req = Config.Request
	} else {
		return nil, errors.New("request is nil")
	}

	schemaValidator.formatFuncMap = Config.FormatFuncMap
	schemaValidator.validatorFuncMap = Config.ValidatorFuncMap
	return schemaValidator, nil
}

func (schemaValidator SchemaValidator) Encode(in interface{}) error {
	value := reflect.ValueOf(in)
	if value.Kind() != reflect.Struct {
		return errors.New(fmt.Sprintf("%s should be a struct kind", in))
	}
	inType := reflect.TypeOf(in)

	for i := 0; i < inType.NumField(); i++ {
		field := inType.Field(i)
		name := field.Name
		tag := field.Tag
		var required bool
		if fieldName := tag.Get("field"); fieldName != "" {
			name = fieldName
		}
		if _, ok := tag.Lookup("required"); ok {
			required = true
		}

		formStr := schemaValidator.req.FormValue(name)
		var formValue interface{}
		formValue = formStr

		if formStr == "" && required {
			return errors.New(fmt.Sprintf("%s is required", name))
		}

		if validatorName := tag.Get("validator"); validatorName != "" {
			if passed, err := schemaValidator.Validate(validatorName, formStr); passed == false {
				return err
			}
		}

		if formatName := tag.Get("format"); formatName != "" {
			if v, err := schemaValidator.Format(validatorName, formStr); err != nil {
				return err
			} else {
				formValue = v
			}
		}

	}
}

func (schemaValidator SchemaValidator) Validate(validatorName, formStr string) (bool, error) {
	if validator, ok := schemaValidator.validatorFuncMap[validatorName]; ok {
		return validator(formStr)
	}
	return false, errors.New(fmt.Sprintf("validator %s is nil", validatorName))
}

func (schemaValidator SchemaValidator) Format(formatName, formStr string) (out interface{}, err error) {
	if formatName, ok := schemaValidator.formatFuncMap[formatName]; ok {
		return formatName(formStr)
	}
	return nil, errors.New(fmt.Sprintf("format %s is nil", formatName))
}
