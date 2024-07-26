package helper

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/fajaramaulana/simple_bank_project/internal/httpapi/handler/response"
	"github.com/gin-gonic/gin"
)

func GlobalCheckingErrorBindJson(errMessage string, request interface{}) (message string, returnError map[string]string) {
	if errMessage == "EOF" {
		message := "Request body is empty"
		return message, map[string]string{
			"error": errMessage,
		}
	}
	returnDataErrorCheck, isExistError := ExtractFieldNameFromError(errMessage, request)
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

func ExtractFieldNameFromError(errorMessage string, request interface{}) (fieldErrorsReturn map[string]string, boolReturn bool) {
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
			fieldNameKey := getJSONTagName(reflect.TypeOf(request), identifier)

			valueValidation := getValueFromTag(reflect.TypeOf(request), identifier, errorMessage)

			// Store the error message in the map using the identifier as the key
			switch errorMessage {
			case "required":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s is required", fieldNameKey)
			case "datetime":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s is not valid datetime", fieldNameKey)
			case "gt": // greater than
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must be greater than %s", fieldNameKey, errorMessage)
			case "gte": // greater than or equal
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must be greater than or equal %s", fieldNameKey, errorMessage)
			case "lt": // less than
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must be less than %s", fieldNameKey, errorMessage)
			case "lte": // less than or equal
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must be less than or equal %s", fieldNameKey, errorMessage)
			case "max": // max length
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must be less than %s characters", fieldNameKey, valueValidation)
			case "min": // min length
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must be greater than %s characters", fieldNameKey, valueValidation)
			case "email":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must be a valid email", fieldNameKey)
			case "eqfield":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must be equal %s", fieldNameKey, errorMessage)
			case "nefield":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must not be equal %s", fieldNameKey, errorMessage)
			case "eqcsfield":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must be equal %s", fieldNameKey, errorMessage)
			case "necsfield":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must not be equal %s", fieldNameKey, errorMessage)
			case "unique":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s is already exists", fieldNameKey)
			case "uuid4":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s is not valid uuid", fieldNameKey)
			case "uuid":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s is not valid uuid", fieldNameKey)
			case "numeric":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must be numeric", fieldNameKey)
			case "alpha":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must be alpha", fieldNameKey)
			case "alphanum":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must be alphanumeric", fieldNameKey)
			case "alphanumunicode":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must be alphanumeric unicode", fieldNameKey)
			case "alphaunicode":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must be alpha unicode", fieldNameKey)
			case "ascii":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must be ascii", fieldNameKey)
			case "contains":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must contain %s", fieldNameKey, errorMessage)
			case "containsany":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must contain any %s", fieldNameKey, errorMessage)
			case "containsrune":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must contain %s", fieldNameKey, errorMessage)
			case "excludes":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must exclude %s", fieldNameKey, errorMessage)
			case "excludesall":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must exclude all %s", fieldNameKey, errorMessage)
			case "excludesrune":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must exclude %s", fieldNameKey, errorMessage)
			case "startswith":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must start with %s", fieldNameKey, errorMessage)
			case "endswith":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must end with %s", fieldNameKey, errorMessage)
			case "customDate":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must be in format dd/mm/yyyy", fieldNameKey)
			case "currency":
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s must be valid currency (%s)", fieldNameKey, Implode(ValidCurrencies, ","))
			default:
				fieldErrors[fieldNameKey] = fmt.Sprintf("%s is %s", fieldNameKey, errorMessage)
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

func getJSONTagName(t reflect.Type, fieldName string) string {
	field, _ := t.FieldByName(fieldName)
	tag := field.Tag.Get("json")
	if tag == "" {
		if tagParts := strings.Split(field.Tag.Get("form"), ","); len(tagParts) > 0 {
			return tagParts[0]
		}
		return fieldName
	}
	tagParts := strings.Split(tag, ",")
	return tagParts[0]
}

func getValueFromTag(t reflect.Type, fieldName string, key string) string {
	field, _ := t.FieldByName(fieldName)
	tag := field.Tag.Get("binding")

	if tag != "" {
		valueValidation := getValueValidation(tag, key)
		return key + " " + valueValidation
	}

	return ""
}

func getValueValidation(data string, key string) string {
	parts := strings.Split(data, ",")
	prefix := key + "="
	for _, part := range parts {
		if strings.HasPrefix(part, prefix) {
			return strings.TrimPrefix(part, prefix)
		}
	}
	return ""
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

func ReturnJSONAbort(ctx *gin.Context, code int, message string, data interface{}) {
	meta := response.Meta{
		Code:    code,
		Status:  http.StatusText(code),
		Message: message,
	}

	response := response.Response{
		Meta: meta,
		Data: data,
	}

	ctx.AbortWithStatusJSON(code, response)
}
