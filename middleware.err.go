package main

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator"
)

// Validation Middleware
// https://github.com/gin-gonic/gin/issues/430

var (
	errorInternalError = errors.New("Woops! Something went wrong :(")
)

func validationErrorToText(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", e.Field(), e.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s", e.Field(), e.Param())
	case "email":
		return fmt.Sprintf("Invalid email format")
	case "len":
		return fmt.Sprintf("%s must be %s characters long", e.Field(), e.Param())
	case "alphanum":
		return fmt.Sprintf("%s must contain letters and numbers %s", e.Field(), e.Param())
	}
	return fmt.Sprintf("%s is not valid", e.Field())
}

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &defaultValidator{}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}

	return nil
}

func (v *defaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")

		// add any custom validations etc. here
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func middlewareErrors() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {

				switch e.Type {
				case gin.ErrorTypePublic:

					if !c.Writer.Written() {
						c.Set("httpStatus", c.Writer.Status())
						render(c, gin.H{"error": e.Error()}, "error.html")
					}
				case gin.ErrorTypeBind:
					errs := e.Err.(validator.ValidationErrors)
					list := make(map[string]string)
					for _, err := range errs {
						list[err.Field()] = validationErrorToText(err)
					}

					c.Set("httpStatus", c.Writer.Status())

					render(c, gin.H{"error": list}, "error.html")
				}

				if !c.Writer.Written() {
					c.Set("httpStatus", 500)
					render(c, gin.H{"error": errorInternalError.Error()}, "error.html")
				}
			}
		}
	}
}
