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
	fmt.Printf("%# v\n", isExistError)
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
			fieldErrors[identifier] = fmt.Sprintf("%s is %s", identifier, errorMessage)
		}

		return fieldErrors, true
	} else {
		fmt.Printf("%# v\n", errorMessage)
		// Define a regular expression pattern to match the field name in the error message
		patternErrJsonUnMarshal := `cannot unmarshal (\S|\s)+ into Go struct field (\S+) of type (\S+)`
		// get value from regexPatternJsonUnmarshallErr
		re := regexp.MustCompile(patternErrJsonUnMarshal)

		// Find matches in the error message
		matches := re.FindStringSubmatch(errorMessage)
		fmt.Printf("%# v\n", matches)
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
