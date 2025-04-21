package handlers

import (
	"github.com/go-playground/validator/v10"
	"pvz/pkg/errors"
	"regexp"
)

var (
	uuidV4Pattern = regexp.MustCompile(
		`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`,
	)
	emailPattern = regexp.MustCompile(
		`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`,
	)
)

var apiValidator = NewApiValidator()

type ApiValidator struct {
	validate *validator.Validate
}

func NewApiValidator() *ApiValidator {
	v := validator.New()
	apiVal := &ApiValidator{validate: v}

	if err := v.RegisterValidation("uuid", validateUuid); err != nil {
		panic(err)
	}
	if err := v.RegisterValidation("email", validateEmail); err != nil {
		panic(err)
	}

	return apiVal
}

func (av *ApiValidator) Validate(i interface{}) error {
	return av.ValidateRequest(i)
}

func (av *ApiValidator) ValidateRequest(req interface{}) error {
	reqErrors := av.validate.Struct(req)
	if reqErrors != nil {
		if validationErrors, ok := reqErrors.(validator.ValidationErrors); ok {
			reqError := validationErrors[0]
			reqName := reqError.Field()
			reqTag := reqError.Tag()
			reqValue := reqError.Value()

			switch reqTag {
			case "required":
				return errors.NewPropertyMissing(reqName)
			case "min":
				return errors.NewPropertyTooSmall(reqName)
			case "max":
				return errors.NewPropertyTooBig(reqName)
			case "oneof":
				return errors.NewWrongPropertyValue(reqName, reqValue.(string))
			case "uuid", "email":
				return errors.NewBadPropertyValue(reqName, reqValue.(string))
			default:
				return errors.NewInternalError()
			}
		}
		return reqErrors
	}
	return nil
}

func (av *ApiValidator) ValidateParam(param interface{}, tag string) error {
	parErrors := av.validate.Var(param, tag)
	if parErrors != nil {
		if validationErrors, ok := parErrors.(validator.ValidationErrors); ok {
			parError := validationErrors[0]
			parTag := parError.Tag()
			parValue := parError.Value()

			switch parTag {
			case "required":
				return errors.NewParamMissing()
			case "uuid":
				return errors.NewBadParamValue(parValue.(string))
			default:
				return errors.NewInternalError()
			}
		}
		return parErrors
	}
	return nil
}

func validateUuid(fl validator.FieldLevel) bool {
	tag, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return uuidV4Pattern.MatchString(tag)
}

func validateEmail(fl validator.FieldLevel) bool {
	tag, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return emailPattern.MatchString(tag)
}
