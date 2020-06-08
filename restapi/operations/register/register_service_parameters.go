// Code generated by go-swagger; DO NOT EDIT.

package register

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"gobully/models"
)

// NewRegisterServiceParams creates a new RegisterServiceParams object
// no default values defined in spec.
func NewRegisterServiceParams() RegisterServiceParams {

	return RegisterServiceParams{}
}

// RegisterServiceParams contains all the bound params for the register service operation
// typically these are obtained from a http.Request
//
// swagger:parameters registerService
type RegisterServiceParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*send register information to get in the network
	  Required: true
	  In: body
	*/
	Service *models.RegisterInfoDTO
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewRegisterServiceParams() beforehand.
func (o *RegisterServiceParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.RegisterInfoDTO
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("service", "body"))
			} else {
				res = append(res, errors.NewParseError("service", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.Service = &body
			}
		}
	} else {
		res = append(res, errors.Required("service", "body"))
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
