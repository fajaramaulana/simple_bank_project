package helper

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type CustomValidation struct {
	Validator *validator.Validate
}

func DoValidation(i interface{}) map[string]string {
	message := map[string]string{}

	val := CustomValidation{validator.New()}

	// add custom validation
	val.Validator.RegisterValidation("customDate", validateCustomDate)

	if err := val.Validator.Struct(i); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			switch e.Tag() {
			case "required":
				message[e.Field()] = fmt.Sprintf("%s is required", e.Field())
			case "datetime":
				message[e.Field()] = fmt.Sprintf("%s is not valid datetime", e.Field())
			case "gt": // greater than
				message[e.Field()] = fmt.Sprintf("%s must be greater than %s", e.Field(), e.Param())
			case "gte": // greater than or equal
				message[e.Field()] = fmt.Sprintf("%s must be greater than or equal %s", e.Field(), e.Param())
			case "lt": // less than
				message[e.Field()] = fmt.Sprintf("%s must be less than %s", e.Field(), e.Param())
			case "lte": // less than or equal
				message[e.Field()] = fmt.Sprintf("%s must be less than or equal %s", e.Field(), e.Param())
			case "max": // max length
				message[e.Field()] = fmt.Sprintf("%s must be less than %s characters", e.Field(), e.Param())
			case "min": // min length
				message[e.Field()] = fmt.Sprintf("%s must be greater than %s characters", e.Field(), e.Param())
			case "email":
				message[e.Field()] = fmt.Sprintf("%s must be a valid email", e.Field())
			case "eqfield":
				message[e.Field()] = fmt.Sprintf("%s must be equal %s", e.Field(), e.Param())
			case "nefield":
				message[e.Field()] = fmt.Sprintf("%s must not be equal %s", e.Field(), e.Param())
			case "eqcsfield":
				message[e.Field()] = fmt.Sprintf("%s must be equal %s", e.Field(), e.Param())
			case "necsfield":
				message[e.Field()] = fmt.Sprintf("%s must not be equal %s", e.Field(), e.Param())
			case "unique":
				message[e.Field()] = fmt.Sprintf("%s is already exists", e.Field())
			case "uuid4":
				message[e.Field()] = fmt.Sprintf("%s is not valid uuid", e.Field())
			case "uuid":
				message[e.Field()] = fmt.Sprintf("%s is not valid uuid", e.Field())
			case "numeric":
				message[e.Field()] = fmt.Sprintf("%s must be numeric", e.Field())
			case "alpha":
				message[e.Field()] = fmt.Sprintf("%s must be alpha", e.Field())
			case "alphanum":
				message[e.Field()] = fmt.Sprintf("%s must be alphanumeric", e.Field())
			case "alphanumunicode":
				message[e.Field()] = fmt.Sprintf("%s must be alphanumeric unicode", e.Field())
			case "alphaunicode":
				message[e.Field()] = fmt.Sprintf("%s must be alpha unicode", e.Field())
			case "ascii":
				message[e.Field()] = fmt.Sprintf("%s must be ascii", e.Field())
			case "contains":
				message[e.Field()] = fmt.Sprintf("%s must contain %s", e.Field(), e.Param())
			case "containsany":
				message[e.Field()] = fmt.Sprintf("%s must contain any %s", e.Field(), e.Param())
			case "containsrune":
				message[e.Field()] = fmt.Sprintf("%s must contain %s", e.Field(), e.Param())
			case "excludes":
				message[e.Field()] = fmt.Sprintf("%s must exclude %s", e.Field(), e.Param())
			case "excludesall":
				message[e.Field()] = fmt.Sprintf("%s must exclude all %s", e.Field(), e.Param())
			case "excludesrune":
				message[e.Field()] = fmt.Sprintf("%s must exclude %s", e.Field(), e.Param())
			case "startswith":
				message[e.Field()] = fmt.Sprintf("%s must start with %s", e.Field(), e.Param())
			case "endswith":
				message[e.Field()] = fmt.Sprintf("%s must end with %s", e.Field(), e.Param())
			case "customDate":
				message[e.Field()] = fmt.Sprintf("%s must be in format dd/mm/yyyy", e.Field())
			case "currency":
				message[e.Field()] = fmt.Sprintf("%s must be valid currency", e.Field())
			}

		}

		return message
	}

	return nil
}

func validateCustomDate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	_, err := time.Parse("02/01/2006", dateStr)
	return err == nil
}

var ValidCurrencies = map[string]bool{
	"USD": true,
	"EUR": true,
	"IDR": true,
}

func CurrencyValidator(fl validator.FieldLevel) bool {
	if currency, ok := fl.Field().Interface().(string); ok {
		return ValidCurrencies[currency]
	}
	return false
}
