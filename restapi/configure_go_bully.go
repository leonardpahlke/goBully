// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"gobully/restapi/operations"
	"gobully/restapi/operations/election"
	"gobully/restapi/operations/register"
)

//go:generate swagger generate server --target ../../goBully --name GoBully --spec ../api/swagger.yml

func configureFlags(api *operations.GoBullyAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.GoBullyAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()
	api.XMLConsumer = runtime.XMLConsumer()

	api.JSONProducer = runtime.JSONProducer()
	api.XMLProducer = runtime.XMLProducer()

	if api.ElectionElectionMessageHandler == nil {
		api.ElectionElectionMessageHandler = election.ElectionMessageHandlerFunc(func(params election.ElectionMessageParams) middleware.Responder {
			return middleware.NotImplemented("operation election.ElectionMessage has not yet been implemented")
		})
	}
	if api.RegisterRegisterServiceHandler == nil {
		api.RegisterRegisterServiceHandler = register.RegisterServiceHandlerFunc(func(params register.RegisterServiceParams) middleware.Responder {
			return middleware.NotImplemented("operation register.RegisterService has not yet been implemented")
		})
	}
	if api.RegisterSendUnregisterToServicesHandler == nil {
		api.RegisterSendUnregisterToServicesHandler = register.SendUnregisterToServicesHandlerFunc(func(params register.SendUnregisterToServicesParams) middleware.Responder {
			return middleware.NotImplemented("operation register.SendUnregisterToServices has not yet been implemented")
		})
	}
	if api.ElectionStartElectionMessageHandler == nil {
		api.ElectionStartElectionMessageHandler = election.StartElectionMessageHandlerFunc(func(params election.StartElectionMessageParams) middleware.Responder {
			return middleware.NotImplemented("operation election.StartElectionMessage has not yet been implemented")
		})
	}
	if api.RegisterTriggerRegisterToServiceHandler == nil {
		api.RegisterTriggerRegisterToServiceHandler = register.TriggerRegisterToServiceHandlerFunc(func(params register.TriggerRegisterToServiceParams) middleware.Responder {
			return middleware.NotImplemented("operation register.TriggerRegisterToService has not yet been implemented")
		})
	}
	if api.RegisterUnregisterFromServiceHandler == nil {
		api.RegisterUnregisterFromServiceHandler = register.UnregisterFromServiceHandlerFunc(func(params register.UnregisterFromServiceParams) middleware.Responder {
			return middleware.NotImplemented("operation register.UnregisterFromService has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
