package handler

import (
	"errors"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	// Register custom validation
	cv.validator.RegisterValidation(
		"contains-uppercase",
		ValidateContainsUppercase,
	)
	cv.validator.RegisterValidation(
		"contains-lowercase",
		ValidateContainsLowercase,
	)
	cv.validator.RegisterValidation(
		"contains-number",
		ValidateContainsNumber,
	)
	cv.validator.RegisterValidation(
		"contains-special-char",
		ValidateContainsSpecialCharacter,
	)

	// Format the error message
	if err := cv.validator.Struct(i); err != nil {
		var errStr []string
		for _, e := range err.(validator.ValidationErrors) {
			errStr = append(
				errStr,
				strings.Trim(
					strings.ToLower(e.Field()[0:1])+e.Field()[1:]+" "+e.Tag()+" "+e.Param(),
					" ",
				),
			)
		}
		return errors.New(strings.Join(errStr, ", "))
	}
	return nil
}

func ValidateCharacterContainsLength(fl validator.FieldLevel, chars string) bool {
	minStr := fl.Param()
	minChar, err := strconv.Atoi(minStr)
	if err != nil {
		return false
	}
	count := 0
	for _, c := range fl.Field().String() {
		if strings.ContainsAny(string(c), chars) {
			count++
		}
	}
	return count >= minChar
}

func ValidateContainsUppercase(fl validator.FieldLevel) bool {
	return ValidateCharacterContainsLength(fl, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

func ValidateContainsLowercase(fl validator.FieldLevel) bool {
	return ValidateCharacterContainsLength(fl, "abcdefghijklmnopqrstuvwxyz")
}

func ValidateContainsNumber(fl validator.FieldLevel) bool {
	return ValidateCharacterContainsLength(fl, "0123456789")
}

func ValidateContainsSpecialCharacter(fl validator.FieldLevel) bool {
	specialCharacters := "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
	return ValidateCharacterContainsLength(fl, specialCharacters)
}
