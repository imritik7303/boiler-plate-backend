package errs

import (
	"net/http"
)

func NewUnAuthorizedError(message string, override bool) *HTTPError {
	return &HTTPError{
		Code: MakeUpperCaseWithUnderscores(http.StatusText(http.StatusUnauthorized)),
		Message: message,
		Status: http.StatusUnauthorized,
		Override: override,
	}
}

func NewForbiddenError(message string ,override bool) *HTTPError {
	return &HTTPError{
		Code: MakeUpperCaseWithUnderscores(http.StatusText(http.StatusForbidden)),
		Message: message,
		Status: http.StatusForbidden,
		Override: override,
	}
}

func NewBadRequestError(message string , override bool , code *string , errors []FieldError , action *Action) *HTTPError {
      formattedcode := MakeUpperCaseWithUnderscores(http.StatusText(http.StatusBadRequest))

	  if code != nil {
		formattedcode = *code
	  }

	  return  &HTTPError{
		Code: formattedcode,
		Message: message,
		Status: http.StatusBadRequest,
		Override: override,
		Errors: errors,
		Action: action,
	  }
}

func NewNotFoundError(message string , override bool ,code *string) *HTTPError {
	formattedcode := MakeUpperCaseWithUnderscores(http.StatusText(http.StatusNotFound))

	if code != nil {
		formattedcode = *code
	}

	return &HTTPError{
		Code: formattedcode,
		Message: message,
		Status: http.StatusNotFound,
		Override: override,
	}
}

func NewInternalServerError() *HTTPError {
	return &HTTPError{
		Code: MakeUpperCaseWithUnderscores(http.StatusText(http.StatusInternalServerError)),
		Message: http.StatusText(http.StatusInternalServerError),
		Status: http.StatusInternalServerError,
		Override: false,
	}
}

func ValidationError(err error) *HTTPError {
	return NewBadRequestError("validation failed :"  + err.Error() , false , nil , nil ,nil)
}