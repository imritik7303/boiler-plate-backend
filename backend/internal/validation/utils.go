package utils

import (
	"strings"

	"github.com/labstack/echo/v4"
	"guthub.com/imritik7303/boiler-plate-backend/internal/errs"
)

type Validatable interface {
	Validate() error
}

type CustomValidationError struct {
	Field   string
	Message string
}

type CustomValidationErrors []CustomValidationError

func (c CustomValidationErrors) Error() string {
	return "validation failed"
}

func BindAndValidate(c echo.Context , payload Validatable) error {
	if err := c.Bind(payload) ; err != nil {
		message := strings.Split(strings.Split(err.Error(),",")[1] , "message")[1]
		return errs.NewBadRequestError(message , false , nil , nil , nil)
	}
	
	if msg, fieldErrors := validateStruct(payload) ; fieldErrors != nil {
		return errs.NewBadRequestError(msg ,true , nil , fieldErrors , nil)
	}

	return nil
}


func BindAndValidateQuery(c echo.Context , payload Validatable) error {
   if err := c.Bind(payload) ; err != nil {
	return errs.NewBadRequestError("Invalid query paramter" , false , nil , nil ,nil )
   }

  
	if msg, fieldErrors := validateStruct(payload) ; fieldErrors != nil {
		return errs.NewBadRequestError(msg ,true , nil , fieldErrors , nil)
	}

	return nil
} 


func validateStruct(v Validatable) (string , []errs.FieldError) {
	if err := v.Validate() ; err != nil {
		return  extractValidationErrors(err)
	}

	return "" , nil
}