package validation

var v *Validator

type Validator interface {
	Var(field interface{}, options interface{}) error
	Struct(s interface{}) error
}

// Singleton will return the global validator variable and create it if necessary
func Singleton() *Validator {
	if v == nil {
		validator := NewPlaygroundValidator()
		v = &validator
	}

	return v
}
