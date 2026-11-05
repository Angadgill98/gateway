package routes

import (
	handlers "gateway/Server/Handlers"
	"net/http"
)



func RegisterEndpointRoutes(mux *http.ServeMux,handler *handlers.Handlers){
	app_mux:=http.NewServeMux()

	app_mux.HandleFunc("/register-endpoint",handler.EndpointHandler.RegisterEndpoint)

	mux.Handle("/",http.StripPrefix("/",app_mux))
}

