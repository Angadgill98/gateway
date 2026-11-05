package routes

import (
	handlers "gateway/Server/Handlers"
	"net/http"
)



func RegisterEndpointRoutes(mux *http.ServeMux, handler *handlers.Handlers) {
    
	mux.HandleFunc("/register-endpoint", handler.EndpointHandler.RegisterEndpoint)
    
	mux.HandleFunc("/update-endpoint", handler.EndpointHandler.UpdateEndpoint)
	
	mux.HandleFunc("/get-all-endpoints", handler.EndpointHandler.GetAllEndpoints)
}
