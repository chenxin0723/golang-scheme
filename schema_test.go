package schema

import (
	"net/http"
	"net/url"
	"testing"
	"unicode/utf8"
)

type People struct {
	Name string `field:"name" required:"true"`
	Url  string `field:"url" validator:"email" required:"true"`
}

func TestAddressable(t *testing.T) {
	req := &http.Request{Form: url.Values{"name": []string{"222"}, "url": []string{"baidu.com"}}}
	schemaValidator, _ := NewSchemaValidator(Config{
		ValidatorFuncMap: map[string]func(in string) (bool, error){
			"max_len_50": func(in string) (bool, error) {
				if utf8.RuneCountInString(in) > 50 {
					return false, nil
				}
				return true, nil
			},
			"max_len_300": func(in string) (bool, error) {
				if utf8.RuneCountInString(in) > 300 {
					return false, nil
				}
				return true, nil
			},
			"max_len_30": func(in string) (bool, error) {
				if utf8.RuneCountInString(in) > 30 {
					return false, nil
				}
				return true, nil
			},
		},
	})
	var people People
	if err := schemaValidator.Encode(&people, req); err != nil {
		t.Errorf("err is -------------- %s\n", err)
	}
}
