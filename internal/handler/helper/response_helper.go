package helper

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/fajaramaulana/simple_bank_project/internal/handler/response"
	"github.com/gin-gonic/gin"
)

func GlobalCheckingErrorBindJson(errMessage string) (message string, returnError map[string]string) {
	if errMessage == "EOF" {
		message := "Request body is empty"
		return message, map[string]string{
			"error": errMessage,
		}
	}
	returnDataErrorCheck, isExistError := ExtractFieldNameFromError(errMessage)
	if isExistError {
		message := "Validation error"
		return message, returnDataErrorCheck
	} else {
		mapReturn := map[string]string{
			"error": errMessage,
		}
		return errMessage, mapReturn
	}
}

func ExtractFieldNameFromError(errorMessage string) (fieldErrorsReturn map[string]string, boolReturn bool) {
	fieldErrors := make(map[string]string)
	// Define a regular expression pattern to match the field name in the error message
	regexPattern := `Key: '([^']+)' Error:Field validation for '([^']+)' failed on the '([^']+)' tag`
	regex := regexp.MustCompile(regexPattern)

	boolReturn = false

	// Find all matches in the error message
	matches := regex.FindAllStringSubmatch(errorMessage, -1)
	if len(matches) > 0 {
		for _, match := range matches {
			fieldName := match[2]
			errorMessage := match[3]

			// Combine the key and field name to form a unique identifier
			identifier := fieldName

			// Store the error message in the map using the identifier as the key
			switch errorMessage {
			case "required":
				fieldErrors[identifier] = fmt.Sprintf("%s is required", identifier)
			case "datetime":
				fieldErrors[identifier] = fmt.Sprintf("%s is not valid datetime", identifier)
			case "gt": // greater than
				fieldErrors[identifier] = fmt.Sprintf("%s must be greater than %s", identifier, errorMessage)
			case "gte": // greater than or equal
				fieldErrors[identifier] = fmt.Sprintf("%s must be greater than or equal %s", identifier, errorMessage)
			case "lt": // less than
				fieldErrors[identifier] = fmt.Sprintf("%s must be less than %s", identifier, errorMessage)
			case "lte": // less than or equal
				fieldErrors[identifier] = fmt.Sprintf("%s must be less than or equal %s", identifier, errorMessage)
			case "max": // max length
				fieldErrors[identifier] = fmt.Sprintf("%s must be less than %s characters", identifier, errorMessage)
			case "min": // min length
				fieldErrors[identifier] = fmt.Sprintf("%s must be greater than %s characters", identifier, errorMessage)
			case "email":
				fieldErrors[identifier] = fmt.Sprintf("%s must be a valid email", identifier)
			case "eqfield":
				fieldErrors[identifier] = fmt.Sprintf("%s must be equal %s", identifier, errorMessage)
			case "nefield":
				fieldErrors[identifier] = fmt.Sprintf("%s must not be equal %s", identifier, errorMessage)
			case "eqcsfield":
				fieldErrors[identifier] = fmt.Sprintf("%s must be equal %s", identifier, errorMessage)
			case "necsfield":
				fieldErrors[identifier] = fmt.Sprintf("%s must not be equal %s", identifier, errorMessage)
			case "unique":
				fieldErrors[identifier] = fmt.Sprintf("%s is already exists", identifier)
			case "uuid4":
				fieldErrors[identifier] = fmt.Sprintf("%s is not valid uuid", identifier)
			case "uuid":
				fieldErrors[identifier] = fmt.Sprintf("%s is not valid uuid", identifier)
			case "numeric":
				fieldErrors[identifier] = fmt.Sprintf("%s must be numeric", identifier)
			case "alpha":
				fieldErrors[identifier] = fmt.Sprintf("%s must be alpha", identifier)
			case "alphanum":
				fieldErrors[identifier] = fmt.Sprintf("%s must be alphanumeric", identifier)
			case "alphanumunicode":
				fieldErrors[identifier] = fmt.Sprintf("%s must be alphanumeric unicode", identifier)
			case "alphaunicode":
				fieldErrors[identifier] = fmt.Sprintf("%s must be alpha unicode", identifier)
			case "ascii":
				fieldErrors[identifier] = fmt.Sprintf("%s must be ascii", identifier)
			case "contains":
				fieldErrors[identifier] = fmt.Sprintf("%s must contain %s", identifier, errorMessage)
			case "containsany":
				fieldErrors[identifier] = fmt.Sprintf("%s must contain any %s", identifier, errorMessage)
			case "containsrune":
				fieldErrors[identifier] = fmt.Sprintf("%s must contain %s", identifier, errorMessage)
			case "excludes":
				fieldErrors[identifier] = fmt.Sprintf("%s must exclude %s", identifier, errorMessage)
			case "excludesall":
				fieldErrors[identifier] = fmt.Sprintf("%s must exclude all %s", identifier, errorMessage)
			case "excludesrune":
				fieldErrors[identifier] = fmt.Sprintf("%s must exclude %s", identifier, errorMessage)
			case "startswith":
				fieldErrors[identifier] = fmt.Sprintf("%s must start with %s", identifier, errorMessage)
			case "endswith":
				fieldErrors[identifier] = fmt.Sprintf("%s must end with %s", identifier, errorMessage)
			case "customDate":
				fieldErrors[identifier] = fmt.Sprintf("%s must be in format dd/mm/yyyy", identifier)
			case "currency":
				fieldErrors[identifier] = fmt.Sprintf("%s must be valid currency (%s)", identifier, Implode(ValidCurrencies, ","))
			default:
				fieldErrors[identifier] = fmt.Sprintf("%s is %s", identifier, errorMessage)
			}
		}

		return fieldErrors, true
	} else {
		// Define a regular expression pattern to match the field name in the error message
		patternErrJsonUnMarshal := `cannot unmarshal (\S|\s)+ into Go struct field (\S+) of type (\S+)`
		// get value from regexPatternJsonUnmarshallErr
		re := regexp.MustCompile(patternErrJsonUnMarshal)

		// Find matches in the error message
		matches := re.FindStringSubmatch(errorMessage)
		if len(matches) > 0 {
			fieldName := matches[2]

			// Split the field name by dots and get the last part
			parts := strings.Split(fieldName, ".")
			fieldName = parts[len(parts)-1]

			fieldType := matches[3]

			// Combine the key and field name to form a unique identifier
			fieldErrors[fieldName] = fmt.Sprintf("%s is not valid, must %s type", fieldName, fieldType)

			return fieldErrors, true
		}
	}

	return fieldErrors, boolReturn
}

func ReturnJSON(ctx *gin.Context, code int, message string, data interface{}) {
	meta := response.Meta{
		Code:    code,
		Status:  http.StatusText(code),
		Message: message,
	}

	response := response.Response{
		Meta: meta,
		Data: data,
	}

	ctx.JSON(code, response)
}

func ReturnJSONError(ctx *gin.Context, code int, message string, data interface{}, err interface{}) {
	meta := response.Meta{
		Code:    code,
		Status:  http.StatusText(code),
		Message: message,
	}

	response := response.ResponseError{
		Meta:  meta,
		Data:  data,
		Error: err,
	}

	ctx.JSON(code, response)
}

func ReturnJSONWithMetaPage(ctx *gin.Context, code int, message string, data interface{}, totalRecords int, totalFiltered int, page int, perPage int) {
	meta := response.MetaPagination{
		Code:          code,
		Status:        http.StatusText(code),
		Message:       message,
		TotalRecords:  totalRecords,
		TotalFiltered: totalFiltered,
		Page:          page,
		PerPage:       perPage,
	}

	response := response.ResponsePagination{
		Meta: meta,
		Data: data,
	}

	ctx.JSON(code, response)
}
