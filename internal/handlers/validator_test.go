package handlers_test

import (
	"github.com/stretchr/testify/assert"
	"pvz/internal/handlers"
	"pvz/pkg/errors"
	"strings"
	"testing"
)

type TestStruct struct {
	Email    string `json:"email" validate:"required,email"`
	UUID     string `json:"uuid" validate:"required,uuid"`
	Username string `json:"username" validate:"required,min=3,max=50"`
}

func TestApiValidator_ValidateRequest(t *testing.T) {
	tests := []struct {
		name          string
		input         interface{}
		expectedError error
	}{
		{
			name: "valid email and uuid",
			input: &TestStruct{
				Email:    "test@example.com",
				UUID:     "f284df64-34de-4c29-b04c-075b1e660850",
				Username: "validuser",
			},
			expectedError: nil,
		},
		{
			name: "invalid email format",
			input: &TestStruct{
				Email:    "invalid-email",
				UUID:     "f284df64-34de-4c29-b04c-075b1e660850",
				Username: "validuser",
			},
			expectedError: errors.NewBadPropertyValue("Email", "invalid-email"),
		},
		{
			name: "invalid UUID format",
			input: &TestStruct{
				Email:    "test@example.com",
				UUID:     "invalid-uuid",
				Username: "validuser",
			},
			expectedError: errors.NewBadPropertyValue("UUID", "invalid-uuid"),
		},
		{
			name: "missing required email",
			input: &TestStruct{
				Email:    "",
				UUID:     "f284df64-34de-4c29-b04c-075b1e660850",
				Username: "validuser",
			},
			expectedError: errors.NewPropertyMissing("Email"),
		},
		{
			name: "username too short",
			input: &TestStruct{
				Email:    "test@example.com",
				UUID:     "f284df64-34de-4c29-b04c-075b1e660850",
				Username: "us",
			},
			expectedError: errors.NewPropertyTooSmall("Username"),
		},
		{
			name: "username too long",
			input: &TestStruct{
				Email:    "test@example.com",
				UUID:     "f284df64-34de-4c29-b04c-075b1e660850",
				Username: strings.Repeat("a", 200),
			},
			expectedError: errors.NewPropertyTooBig("Username"),
		},
	}

	apiValidator := handlers.NewApiValidator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apiValidator.ValidateRequest(tt.input)
			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestApiValidator_ValidateParam(t *testing.T) {
	tests := []struct {
		name          string
		param         interface{}
		tag           string
		expectedError error
	}{
		{
			name:          "valid uuid",
			param:         "f284df64-34de-4c29-b04c-075b1e660850",
			tag:           "uuid",
			expectedError: nil,
		},
		{
			name:          "invalid uuid",
			param:         "invalid-uuid",
			tag:           "uuid",
			expectedError: errors.NewBadParamValue("invalid-uuid"),
		},
		{
			name:          "valid email",
			param:         "test@example.com",
			tag:           "email",
			expectedError: nil,
		},
		{
			name:          "missing required parameter",
			param:         "",
			tag:           "required",
			expectedError: errors.NewParamMissing(),
		},
	}

	apiValidator := handlers.NewApiValidator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apiValidator.ValidateParam(tt.param, tt.tag)
			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
