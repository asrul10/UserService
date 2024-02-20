package handler

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestValidateContainsUppercase(t *testing.T) {
	customValidator := &CustomValidator{
		validator: validator.New(),
	}
	customValidator.validator.RegisterValidation("contains-uppercase", ValidateContainsUppercase)

	type TestStruct struct {
		Password string `validate:"contains-uppercase=1"`
	}

	tests := []struct {
		caseName string
		input    TestStruct
		expected string
	}{
		{
			caseName: "Single uppercase",
			input:    TestStruct{Password: "aA"},
			expected: "",
		},
		{
			caseName: "No uppercase",
			input:    TestStruct{Password: "aa"},
			expected: "password contains-uppercase 1",
		},
		{
			caseName: "Multiple uppercase",
			input:    TestStruct{Password: "lkajwdljiaowjIjsalkdj"},
			expected: "",
		},
		{
			caseName: "Special characters and multiple uppercase",
			input:    TestStruct{Password: "aaaasj982(*&(*))AaljsdlakjUUU"},
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			err := customValidator.Validate(test.input)
			if err != nil && err.Error() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, err.Error())
			}
			if err == nil && test.expected != "" {
				t.Errorf("Expected %s, got nil", test.expected)
			}
		})
	}
}

func TestValidateContainsLowercase(t *testing.T) {
	customValidator := &CustomValidator{
		validator: validator.New(),
	}
	customValidator.validator.RegisterValidation("contains-lowercase", ValidateContainsLowercase)

	type TestStruct struct {
		Password string `validate:"contains-lowercase=1"`
	}

	tests := []struct {
		caseName string
		input    TestStruct
		expected string
	}{
		{
			caseName: "Single lowercase",
			input:    TestStruct{Password: "aA"},
			expected: "",
		},
		{
			caseName: "No lowercase",
			input:    TestStruct{Password: "AA"},
			expected: "password contains-lowercase 1",
		},
		{
			caseName: "Multiple lowercase",
			input:    TestStruct{Password: "lkajwdljiaowjIjsalkdj"},
			expected: "",
		},
		{
			caseName: "Special characters and multiple lowercase",
			input:    TestStruct{Password: "aaaasj982(*&(*))AaljsdlakjUUU"},
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			err := customValidator.Validate(test.input)
			if err != nil && err.Error() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, err.Error())
			}
			if err == nil && test.expected != "" {
				t.Errorf("Expected %s, got nil", test.expected)
			}
		})
	}
}

func TestValidateContainsNumber(t *testing.T) {
	customValidator := &CustomValidator{
		validator: validator.New(),
	}
	customValidator.validator.RegisterValidation("contains-number", ValidateContainsNumber)

	type TestStruct struct {
		Password string `validate:"contains-number=1"`
	}

	tests := []struct {
		caseName string
		input    TestStruct
		expected string
	}{
		{
			caseName: "No number",
			input:    TestStruct{Password: "aA"},
			expected: "password contains-number 1",
		},
		{
			caseName: "No number",
			input:    TestStruct{Password: "AA"},
			expected: "password contains-number 1",
		},
		{
			caseName: "No number",
			input:    TestStruct{Password: "lkajwdljiaowjIjsalkdj"},
			expected: "password contains-number 1",
		},
		{
			caseName: "Special characters and multiple numbers",
			input:    TestStruct{Password: "aaaasj982(*&(*))AaljsdlakjUUU"},
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			err := customValidator.Validate(test.input)
			if err != nil && err.Error() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, err.Error())
			}
			if err == nil && test.expected != "" {
				t.Errorf("Expected %s, got nil", test.expected)
			}
		})
	}
}

func TestValidateContainsSpecialCharacter(t *testing.T) {
	customValidator := &CustomValidator{
		validator: validator.New(),
	}
	customValidator.validator.RegisterValidation("contains-special-char", ValidateContainsSpecialCharacter)

	type TestStruct struct {
		Password string `validate:"contains-special-char=1"`
	}

	tests := []struct {
		caseName string
		input    TestStruct
		expected string
	}{
		{
			caseName: "No special character",
			input:    TestStruct{Password: "aA"},
			expected: "password contains-special-char 1",
		},
		{
			caseName: "No special character",
			input:    TestStruct{Password: "AA"},
			expected: "password contains-special-char 1",
		},
		{
			caseName: "No special character",
			input:    TestStruct{Password: "lkajwdljiaowjIjsalkdj"},
			expected: "password contains-special-char 1",
		},
		{
			caseName: "Special characters and multiple numbers",
			input:    TestStruct{Password: "aaaasj982(*&(*))AaljsdlakjUUU"},
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.caseName, func(t *testing.T) {
			err := customValidator.Validate(test.input)
			if err != nil && err.Error() != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, err.Error())
			}
			if err == nil && test.expected != "" {
				t.Errorf("Expected %s, got nil", test.expected)
			}
		})
	}
}
