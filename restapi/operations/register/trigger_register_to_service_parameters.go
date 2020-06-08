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
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"

	"gobully/models"
)

// NewTriggerRegisterToServiceParams creates a new TriggerRegisterToServiceParams object
// no default values defined in spec.
func NewTriggerRegisterToServiceParams() TriggerRegisterToServiceParams {

	return TriggerRegisterToServiceParams{}
}

// TriggerRegisterToServiceParams contains all the bound params for the trigger register to service operation
// typically these are obtained from a http.Request
//
// swagger:parameters triggerRegisterToService
type TriggerRegisterToServiceParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*start election algorithm - to get a coordinator
	  Required: true
	  In: body
	*/
	ElectionInformation *models.InformationElectionDTO
	/*trigger registration, service sends registration message to other
	  Required: true
	  In: query
	*/
	IP float64
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewTriggerRegisterToServiceParams() beforehand.
func (o *TriggerRegisterToServiceParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	if runtime.HasBody(r) {
		defer r.Body.Close()
		var body models.InformationElectionDTO
		if err := route.Consumer.Consume(r.Body, &body); err != nil {
			if err == io.EOF {
				res = append(res, errors.Required("electionInformation", "body"))
			} else {
				res = append(res, errors.NewParseError("electionInformation", "body", "", err))
			}
		} else {
			// validate body object
			if err := body.Validate(route.Formats); err != nil {
				res = append(res, err)
			}

			if len(res) == 0 {
				o.ElectionInformation = &body
			}
		}
	} else {
		res = append(res, errors.Required("electionInformation", "body"))
	}
	qIP, qhkIP, _ := qs.GetOK("ip")
	if err := o.bindIP(qIP, qhkIP, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindIP binds and validates parameter IP from query.
func (o *TriggerRegisterToServiceParams) bindIP(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("ip", "query")
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false
	if err := validate.RequiredString("ip", "query", raw); err != nil {
		return err
	}

	value, err := swag.ConvertFloat64(raw)
	if err != nil {
		return errors.InvalidType("ip", "query", "float64", raw)
	}
	o.IP = value

	return nil
}
