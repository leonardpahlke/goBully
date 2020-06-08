// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/security"
	"github.com/go-openapi/spec"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"gobully/restapi/operations/election"
	"gobully/restapi/operations/register"
)

// NewGoBullyAPI creates a new GoBully instance
func NewGoBullyAPI(spec *loads.Document) *GoBullyAPI {
	return &GoBullyAPI{
		handlers:            make(map[string]map[string]http.Handler),
		formats:             strfmt.Default,
		defaultConsumes:     "application/json",
		defaultProduces:     "application/json",
		customConsumers:     make(map[string]runtime.Consumer),
		customProducers:     make(map[string]runtime.Producer),
		PreServerShutdown:   func() {},
		ServerShutdown:      func() {},
		spec:                spec,
		ServeError:          errors.ServeError,
		BasicAuthenticator:  security.BasicAuth,
		APIKeyAuthenticator: security.APIKeyAuth,
		BearerAuthenticator: security.BearerAuth,

		JSONConsumer: runtime.JSONConsumer(),
		XMLConsumer:  runtime.XMLConsumer(),

		JSONProducer: runtime.JSONProducer(),
		XMLProducer:  runtime.XMLProducer(),

		ElectionElectionMessageHandler: election.ElectionMessageHandlerFunc(func(params election.ElectionMessageParams) middleware.Responder {
			return middleware.NotImplemented("operation election.ElectionMessage has not yet been implemented")
		}),
		RegisterRegisterServiceHandler: register.RegisterServiceHandlerFunc(func(params register.RegisterServiceParams) middleware.Responder {
			return middleware.NotImplemented("operation register.RegisterService has not yet been implemented")
		}),
		RegisterSendUnregisterToServicesHandler: register.SendUnregisterToServicesHandlerFunc(func(params register.SendUnregisterToServicesParams) middleware.Responder {
			return middleware.NotImplemented("operation register.SendUnregisterToServices has not yet been implemented")
		}),
		ElectionStartElectionMessageHandler: election.StartElectionMessageHandlerFunc(func(params election.StartElectionMessageParams) middleware.Responder {
			return middleware.NotImplemented("operation election.StartElectionMessage has not yet been implemented")
		}),
		RegisterTriggerRegisterToServiceHandler: register.TriggerRegisterToServiceHandlerFunc(func(params register.TriggerRegisterToServiceParams) middleware.Responder {
			return middleware.NotImplemented("operation register.TriggerRegisterToService has not yet been implemented")
		}),
		RegisterUnregisterFromServiceHandler: register.UnregisterFromServiceHandlerFunc(func(params register.UnregisterFromServiceParams) middleware.Responder {
			return middleware.NotImplemented("operation register.UnregisterFromService has not yet been implemented")
		}),
	}
}

/*GoBullyAPI This project implements the bully algorithm with docker containers.
Several containers are served, each of which is accessible with a rest API.
For more information, see the code comments */
type GoBullyAPI struct {
	spec            *loads.Document
	context         *middleware.Context
	handlers        map[string]map[string]http.Handler
	formats         strfmt.Registry
	customConsumers map[string]runtime.Consumer
	customProducers map[string]runtime.Producer
	defaultConsumes string
	defaultProduces string
	Middleware      func(middleware.Builder) http.Handler

	// BasicAuthenticator generates a runtime.Authenticator from the supplied basic auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BasicAuthenticator func(security.UserPassAuthentication) runtime.Authenticator
	// APIKeyAuthenticator generates a runtime.Authenticator from the supplied token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	APIKeyAuthenticator func(string, string, security.TokenAuthentication) runtime.Authenticator
	// BearerAuthenticator generates a runtime.Authenticator from the supplied bearer token auth function.
	// It has a default implementation in the security package, however you can replace it for your particular usage.
	BearerAuthenticator func(string, security.ScopedTokenAuthentication) runtime.Authenticator

	// JSONConsumer registers a consumer for the following mime types:
	//   - application/json
	JSONConsumer runtime.Consumer
	// XMLConsumer registers a consumer for the following mime types:
	//   - application/xml
	XMLConsumer runtime.Consumer

	// JSONProducer registers a producer for the following mime types:
	//   - application/json
	JSONProducer runtime.Producer
	// XMLProducer registers a producer for the following mime types:
	//   - application/xml
	XMLProducer runtime.Producer

	// ElectionElectionMessageHandler sets the operation handler for the election message operation
	ElectionElectionMessageHandler election.ElectionMessageHandler
	// RegisterRegisterServiceHandler sets the operation handler for the register service operation
	RegisterRegisterServiceHandler register.RegisterServiceHandler
	// RegisterSendUnregisterToServicesHandler sets the operation handler for the send unregister to services operation
	RegisterSendUnregisterToServicesHandler register.SendUnregisterToServicesHandler
	// ElectionStartElectionMessageHandler sets the operation handler for the start election message operation
	ElectionStartElectionMessageHandler election.StartElectionMessageHandler
	// RegisterTriggerRegisterToServiceHandler sets the operation handler for the trigger register to service operation
	RegisterTriggerRegisterToServiceHandler register.TriggerRegisterToServiceHandler
	// RegisterUnregisterFromServiceHandler sets the operation handler for the unregister from service operation
	RegisterUnregisterFromServiceHandler register.UnregisterFromServiceHandler
	// ServeError is called when an error is received, there is a default handler
	// but you can set your own with this
	ServeError func(http.ResponseWriter, *http.Request, error)

	// PreServerShutdown is called before the HTTP(S) server is shutdown
	// This allows for custom functions to get executed before the HTTP(S) server stops accepting traffic
	PreServerShutdown func()

	// ServerShutdown is called when the HTTP(S) server is shut down and done
	// handling all active connections and does not accept connections any more
	ServerShutdown func()

	// Custom command line argument groups with their descriptions
	CommandLineOptionsGroups []swag.CommandLineOptionsGroup

	// User defined logger function.
	Logger func(string, ...interface{})
}

// SetDefaultProduces sets the default produces media type
func (o *GoBullyAPI) SetDefaultProduces(mediaType string) {
	o.defaultProduces = mediaType
}

// SetDefaultConsumes returns the default consumes media type
func (o *GoBullyAPI) SetDefaultConsumes(mediaType string) {
	o.defaultConsumes = mediaType
}

// SetSpec sets a spec that will be served for the clients.
func (o *GoBullyAPI) SetSpec(spec *loads.Document) {
	o.spec = spec
}

// DefaultProduces returns the default produces media type
func (o *GoBullyAPI) DefaultProduces() string {
	return o.defaultProduces
}

// DefaultConsumes returns the default consumes media type
func (o *GoBullyAPI) DefaultConsumes() string {
	return o.defaultConsumes
}

// Formats returns the registered string formats
func (o *GoBullyAPI) Formats() strfmt.Registry {
	return o.formats
}

// RegisterFormat registers a custom format validator
func (o *GoBullyAPI) RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	o.formats.Add(name, format, validator)
}

// Validate validates the registrations in the GoBullyAPI
func (o *GoBullyAPI) Validate() error {
	var unregistered []string

	if o.JSONConsumer == nil {
		unregistered = append(unregistered, "JSONConsumer")
	}
	if o.XMLConsumer == nil {
		unregistered = append(unregistered, "XMLConsumer")
	}

	if o.JSONProducer == nil {
		unregistered = append(unregistered, "JSONProducer")
	}
	if o.XMLProducer == nil {
		unregistered = append(unregistered, "XMLProducer")
	}

	if o.ElectionElectionMessageHandler == nil {
		unregistered = append(unregistered, "election.ElectionMessageHandler")
	}
	if o.RegisterRegisterServiceHandler == nil {
		unregistered = append(unregistered, "register.RegisterServiceHandler")
	}
	if o.RegisterSendUnregisterToServicesHandler == nil {
		unregistered = append(unregistered, "register.SendUnregisterToServicesHandler")
	}
	if o.ElectionStartElectionMessageHandler == nil {
		unregistered = append(unregistered, "election.StartElectionMessageHandler")
	}
	if o.RegisterTriggerRegisterToServiceHandler == nil {
		unregistered = append(unregistered, "register.TriggerRegisterToServiceHandler")
	}
	if o.RegisterUnregisterFromServiceHandler == nil {
		unregistered = append(unregistered, "register.UnregisterFromServiceHandler")
	}

	if len(unregistered) > 0 {
		return fmt.Errorf("missing registration: %s", strings.Join(unregistered, ", "))
	}

	return nil
}

// ServeErrorFor gets a error handler for a given operation id
func (o *GoBullyAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return o.ServeError
}

// AuthenticatorsFor gets the authenticators for the specified security schemes
func (o *GoBullyAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]runtime.Authenticator {
	return nil
}

// Authorizer returns the registered authorizer
func (o *GoBullyAPI) Authorizer() runtime.Authorizer {
	return nil
}

// ConsumersFor gets the consumers for the specified media types.
// MIME type parameters are ignored here.
func (o *GoBullyAPI) ConsumersFor(mediaTypes []string) map[string]runtime.Consumer {
	result := make(map[string]runtime.Consumer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONConsumer
		case "application/xml":
			result["application/xml"] = o.XMLConsumer
		}

		if c, ok := o.customConsumers[mt]; ok {
			result[mt] = c
		}
	}
	return result
}

// ProducersFor gets the producers for the specified media types.
// MIME type parameters are ignored here.
func (o *GoBullyAPI) ProducersFor(mediaTypes []string) map[string]runtime.Producer {
	result := make(map[string]runtime.Producer, len(mediaTypes))
	for _, mt := range mediaTypes {
		switch mt {
		case "application/json":
			result["application/json"] = o.JSONProducer
		case "application/xml":
			result["application/xml"] = o.XMLProducer
		}

		if p, ok := o.customProducers[mt]; ok {
			result[mt] = p
		}
	}
	return result
}

// HandlerFor gets a http.Handler for the provided operation method and path
func (o *GoBullyAPI) HandlerFor(method, path string) (http.Handler, bool) {
	if o.handlers == nil {
		return nil, false
	}
	um := strings.ToUpper(method)
	if _, ok := o.handlers[um]; !ok {
		return nil, false
	}
	if path == "/" {
		path = ""
	}
	h, ok := o.handlers[um][path]
	return h, ok
}

// Context returns the middleware context for the go bully API
func (o *GoBullyAPI) Context() *middleware.Context {
	if o.context == nil {
		o.context = middleware.NewRoutableContext(o.spec, o, nil)
	}

	return o.context
}

func (o *GoBullyAPI) initHandlerCache() {
	o.Context() // don't care about the result, just that the initialization happened
	if o.handlers == nil {
		o.handlers = make(map[string]map[string]http.Handler)
	}

	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/election"] = election.NewElectionMessage(o.context, o.ElectionElectionMessageHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/register"] = register.NewRegisterService(o.context, o.RegisterRegisterServiceHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/sendunregister"] = register.NewSendUnregisterToServices(o.context, o.RegisterSendUnregisterToServicesHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/startelection"] = election.NewStartElectionMessage(o.context, o.ElectionStartElectionMessageHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/sendregister"] = register.NewTriggerRegisterToService(o.context, o.RegisterTriggerRegisterToServiceHandler)
	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/unregister"] = register.NewUnregisterFromService(o.context, o.RegisterUnregisterFromServiceHandler)
}

// Serve creates a http handler to serve the API over HTTP
// can be used directly in http.ListenAndServe(":8000", api.Serve(nil))
func (o *GoBullyAPI) Serve(builder middleware.Builder) http.Handler {
	o.Init()

	if o.Middleware != nil {
		return o.Middleware(builder)
	}
	return o.context.APIHandler(builder)
}

// Init allows you to just initialize the handler cache, you can then recompose the middleware as you see fit
func (o *GoBullyAPI) Init() {
	if len(o.handlers) == 0 {
		o.initHandlerCache()
	}
}

// RegisterConsumer allows you to add (or override) a consumer for a media type.
func (o *GoBullyAPI) RegisterConsumer(mediaType string, consumer runtime.Consumer) {
	o.customConsumers[mediaType] = consumer
}

// RegisterProducer allows you to add (or override) a producer for a media type.
func (o *GoBullyAPI) RegisterProducer(mediaType string, producer runtime.Producer) {
	o.customProducers[mediaType] = producer
}

// AddMiddlewareFor adds a http middleware to existing handler
func (o *GoBullyAPI) AddMiddlewareFor(method, path string, builder middleware.Builder) {
	um := strings.ToUpper(method)
	if path == "/" {
		path = ""
	}
	o.Init()
	if h, ok := o.handlers[um][path]; ok {
		o.handlers[method][path] = builder(h)
	}
}
