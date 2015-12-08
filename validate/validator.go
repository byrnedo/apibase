package validate

import (
	validator "gopkg.in/bluesuncorp/validator.v8"
)

var V *validator.Validate

func init() {
	V = validator.New(&validator.Config{TagName: "validate"})
}

func ValidateStruct(any interface{}) (valErrors validator.ValidationErrors) {

	if errs := V.Struct(any); errs != nil {
		valErrors = errs.(validator.ValidationErrors)
	}
	return
}
