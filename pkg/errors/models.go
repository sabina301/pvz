package errors

import "fmt"

type commonError struct {
	Message string `json:"message"`
}

func (e commonError) Error() string {
	return e.Message
}

type MalformedBody struct {
	commonError
}

func NewMalformedBody() MalformedBody {
	msg := "malformed body"
	return MalformedBody{
		commonError: commonError{
			Message: msg,
		},
	}
}

type InternalError struct {
	commonError
}

func NewInternalError() InternalError {
	msg := "internal error"
	return InternalError{
		commonError: commonError{
			Message: msg,
		},
	}
}

type ObjectNotFound struct {
	commonError
}

func NewObjectNotFound(objectType string) ObjectNotFound {
	msg := fmt.Sprintf("%s not found", objectType)
	return ObjectNotFound{
		commonError: commonError{
			Message: msg,
		},
	}
}

type InvalidCredentials struct {
	commonError
}

func NewInvalidCredentials() InvalidCredentials {
	msg := "invalid credentials"
	return InvalidCredentials{
		commonError: commonError{
			Message: msg,
		},
	}
}

type StartDateAfterEndDate struct {
	commonError
}

func NewStartDateAfterEndDate() StartDateAfterEndDate {
	msg := "start date after end date"
	return StartDateAfterEndDate{
		commonError: commonError{
			Message: msg,
		},
	}
}

type ReceptionIsNotClosed struct {
	commonError
}

func NewReceptionIsNotClosed(pvzId string) ReceptionIsNotClosed {
	msg := fmt.Sprintf("reception in pvz with id %s is not closed", pvzId)
	return ReceptionIsNotClosed{
		commonError: commonError{
			Message: msg,
		},
	}
}

type ReceptionIsNotInProgress struct {
	commonError
}

func NewReceptionIsNotInProgress(recId string) ReceptionIsNotInProgress {
	msg := fmt.Sprintf("reception with id %s is not closed", recId)
	return ReceptionIsNotInProgress{
		commonError: commonError{
			Message: msg,
		},
	}
}

type NoInProgressReception struct {
	commonError
}

func NewNoInProgressReception() NoInProgressReception {
	msg := "no in-progress reception"
	return NoInProgressReception{
		commonError: commonError{
			Message: msg,
		},
	}
}

type ObjectHasNotSubObjects struct {
	commonError
}

func NewObjectHasNotSubObjects(object, subObject string) ObjectHasNotSubObjects {
	msg := fmt.Sprintf("%s has not %s", object, subObject)
	return ObjectHasNotSubObjects{
		commonError: commonError{
			Message: msg,
		},
	}
}

type ObjectAlreadyExists struct {
	commonError
}

func NewObjectAlreadyExists(object, property, value string) ObjectAlreadyExists {
	msg := fmt.Sprintf("%s with %s %s already exists", object, property, value)
	return ObjectAlreadyExists{
		commonError: commonError{
			Message: msg,
		},
	}
}

type PropertyMissing struct {
	commonError
}

func NewPropertyMissing(property string) PropertyMissing {
	msg := fmt.Sprintf("property %s is missing", property)
	return PropertyMissing{
		commonError: commonError{
			Message: msg,
		},
	}
}

type ParamMissing struct {
	commonError
}

func NewParamMissing() ParamMissing {
	msg := "param is missing"
	return ParamMissing{
		commonError: commonError{
			Message: msg,
		},
	}
}

type PropertyTooSmall struct {
	commonError
	Property string `json:"property"`
}

func NewPropertyTooSmall(property string) PropertyTooSmall {
	msg := fmt.Sprintf("property %s is too small", property)
	return PropertyTooSmall{
		commonError: commonError{
			Message: msg,
		},
		Property: property,
	}
}

type PropertyTooBig struct {
	commonError
	Property string `json:"property"`
}

func NewPropertyTooBig(property string) PropertyTooBig {
	msg := fmt.Sprintf("property %s is too big", property)
	return PropertyTooBig{
		commonError: commonError{
			Message: msg,
		},
		Property: property,
	}
}

type WrongPropertyValue struct {
	commonError
}

func NewWrongPropertyValue(property, value string) WrongPropertyValue {
	msg := fmt.Sprintf("property %s has wrong value %s", property, value)
	return WrongPropertyValue{
		commonError: commonError{
			Message: msg,
		},
	}
}

type BadPropertyValue struct {
	commonError
}

func NewBadPropertyValue(property, value string) BadPropertyValue {
	msg := fmt.Sprintf("property %s has bad value format %s", property, value)
	return BadPropertyValue{
		commonError: commonError{
			Message: msg,
		},
	}
}

type BadParamValue struct {
	commonError
}

func NewBadParamValue(value string) BadParamValue {
	msg := fmt.Sprintf("bad param value %s", value)
	return BadParamValue{
		commonError: commonError{
			Message: msg,
		},
	}
}

type AccessForbidden struct {
	commonError
}

func NewAccessForbidden() AccessForbidden {
	msg := "access forbidden"
	return AccessForbidden{
		commonError: commonError{
			Message: msg,
		},
	}
}
