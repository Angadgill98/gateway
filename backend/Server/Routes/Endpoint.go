package routes

import (
	handlers "gateway/Server/Handlers"
	"net/http"
)



func RegisterEndpointRoutes(mux *http.ServeMux,handler *handlers.Handlers){
	app_mux:=http.NewServeMux()

	app_mux.HandleFunc("/register-endpoint",handler.EndpointHandler.RegisterEndpoint)

	app_mux.HandleFunc("/update-endpoint", handler.EndpointHandler.UpdateEndpoint)

	mux.Handle("/api/", http.StripPrefix("/api", app_mux))
}

