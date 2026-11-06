package routes

import (
	handlers "gateway/Server/Handlers"
	"net/http"
)







func RegisterRoutes(mux *http.ServeMux) {
    handlers := handlers.CreateHandlers()

    api_mux := http.NewServeMux()

    RegisterAppRoutes(api_mux, handlers)
    RegisterEndpointRoutes(api_mux, handlers)
	RegisterDestinationRoutes(api_mux,handlers)

    mux.Handle("/api/", http.StripPrefix("/api", api_mux))
}