package schema

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

var commomValidators = map[string]func(in string) (bool, error){
	"email": func(in string) (bool, error) {
		if !rxEmail.MatchString(in) {
			return false, errors.New(fmt.Sprintf("%s is not passed", "email"))
		}
		return true, nil
	},
	"url": func(in string) (bool, error) {
		if !rxURL.MatchString(in) {
			return false, errors.New(fmt.Sprintf("%s is not passed", "url"))
		}
		return true, nil
	},
}

type SchemaValidator struct {
	formatFuncMap    map[string]func(in string) (out interface{}, err error)
	validatorFuncMap map[string]func(in string) (bool, error)
}

type Config struct {
	FormatFuncMap    map[string]func(in string) (out interface{}, err error)
	ValidatorFuncMap map[string]func(in string) (bool, error)
}

func NewSchemaValidator(Config Config) (SchemaValidator, error) {
	schemaValidator := SchemaValidator{}
	schemaValidator.formatFuncMap = Config.FormatFuncMap
	if Config.ValidatorFuncMap != nil {
		for k, v := range commomValidators {
			if _, ok := Config.ValidatorFuncMap[k]; !ok {
				Config.ValidatorFuncMap[k] = v
			}
		}
		schemaValidator.validatorFuncMap = Config.ValidatorFuncMap
	} else {
		schemaValidator.validatorFuncMap = commomValidators
	}

	return schemaValidator, nil
}

func (schemaValidator SchemaValidator) Encode(in interface{}, req *http.Request) error {

	value := reflect.Indirect(reflect.ValueOf(in))
	if value.Kind() != reflect.Struct {
		return errors.New(fmt.Sprintf("%s should be a struct kind", in))
	}
	if !value.CanSet() {
		return errors.New(fmt.Sprintf("%+v should be addressable", in))

	}

	inType := reflect.TypeOf(value.Interface())
	for i := 0; i < inType.NumField(); i++ {
		field := inType.Field(i)
		fieldValue := value.Field(i)
		name := strings.ToLower(field.Name)
		// name := field.Name
		tag := field.Tag
		var required bool
		if fieldName := tag.Get("field"); fieldName != "" {

			name = fieldName
		}
		if _, ok := tag.Lookup("required"); ok {
			required = true
		}

		formStr := req.FormValue(name)
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
			if v, err := schemaValidator.Format(formatName, formStr); err != nil {
				return err
			} else {
				formValue = v
			}
		}

		switch fieldValue.Kind() {
		case reflect.Int:
			i, err := strconv.Atoi(formValue.(string))
			if err == nil {
				fieldValue.SetInt(int64(i))
			} else {
				fieldValue.SetInt(int64(0))
			}
		case reflect.String:
			fieldValue.SetString(formValue.(string))
		}

	}
	return nil
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
