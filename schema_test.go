package schema

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

type People struct {
	Name  string `format:"uppercase" required:"true"`
	Age   int    `validator:"max_age_150"`
	Email string `field:"my_email" validator:"email" required:"true"`
}

func TestAddressable(t *testing.T) {
	req := &http.Request{Form: url.Values{"name": []string{"xin"}, "age": []string{"22"}, "my_email": []string{"xin@theplant.jp"}}}
	schemaValidator, _ := NewSchemaValidator(Config{
		ValidatorFuncMap: map[string]func(in string) (bool, error){
			"max_age_150": func(in string) (bool, error) {
				i, err := strconv.Atoi(in)
				if err != nil {
					return false, errors.New("age is not a integer")
				}
				if i > 150 {
					return false, errors.New("Age should be within 150")
				}
				return true, nil
			},
		},
		FormatFuncMap: map[string]func(in string) (out interface{}, err error){
			"uppercase": func(in string) (out interface{}, err error) {
				return strings.ToUpper(in), nil
			},
		},
	})
	var people People
	if err := schemaValidator.Encode(&people, req); err != nil {
		t.Errorf("err is -------------- %s\n", err)
	}
}
