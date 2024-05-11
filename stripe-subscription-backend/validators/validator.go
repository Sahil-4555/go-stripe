package validators

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator"
)

func GetError(field, tag string) string {
	if tag == "required" {
		return field + " is " + tag
	} else if tag == "email" {
		return field + " not valid"
	}
	return ""
}

func ValidateStruct(req interface{}, key string) (string, bool) {
	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		valErrs := err.(validator.ValidationErrors)
		for _, v := range valErrs {
			fieldName := strings.Replace(strings.Replace(v.Namespace(), key+".", "", 1), ".", " ", 3)
			reg, _ := regexp.Compile("[^A-Z`[]]+")
			fieldName = strings.Replace(reg.ReplaceAllString(fieldName, ""), "[", "", 2)
			errorString := GetError(fieldName, v.Tag())
			if errorString == "" {
				errorString = "Some of required field are missing or invalid"
			}
			return errorString, false
		}
	}
	return "", true
}
