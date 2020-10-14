package helpers

import "strings"

type Validator struct {
	Message string                 `json:"message"`
	Errors  map[string]interface{} `json:"errors"`
}

func ValidatorRemoveNamespace(s string) string {
	splittedErrorNamespace := strings.Split(s, ".")
	slicedErrorField := splittedErrorNamespace[1:]
	errorField := strings.Join(slicedErrorField, ".")

	splittedErrorNamespaceTwo := strings.Split(errorField, "]")
	errorField = strings.Join(splittedErrorNamespaceTwo, ".")

	splittedErrorNamespaceThree := strings.Split(errorField, "[")
	errorField = strings.Join(splittedErrorNamespaceThree, ".")

	splittedErrorNamespaceFour := strings.Split(errorField, "..")
	errorField = strings.Join(splittedErrorNamespaceFour, ".")

	return errorField
}

func ValidatorMessage(fieldName string, validationType string, validationParam string) string {

	splittedFieldName := strings.Split(fieldName, "_")
	fieldName = strings.Join(splittedFieldName, " ")
	message := "The "

	if validationType == "required" {
		message += fieldName + " field is required"
	}

	if validationType == "min" {
		message += fieldName + " minimal length is " + validationParam
	}

	return message
}
