

schema-validator
======

Package schema converts structs to and from form values and also provide some custom verification and formatting.

## Example

Here's a quick example: we parse POST form values and then decode them into a struct:

```go
// Set a Decoder instance as a package global, because it caches
// meta-data about structs, and an instance can be shared safely.

type People struct {
	Name  string `format:"uppercase" required:"true"`
	Age   int    `validator:"max_age_150"`
	Email string `field:"my_email" validator:"email" required:"true"`
}


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

func MyHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    if err != nil {
        // Handle error
    }

    var person Person

    // r.PostForm is a map of our POST form values
    err = schemaValidator.Decode(&person, r)
    if err != nil {
        // Handle error
    }

    // Do something with person.Name or person.Phone
}
```


## License

BSD licensed. See the LICENSE file for details.
