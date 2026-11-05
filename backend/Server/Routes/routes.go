package routes

import (
	handlers "gateway/Server/Handlers"
	"net/http"
)







func RegisterRoutes(mux *http.ServeMux){
	handlers:=handlers.CreateHandlers()
	RegisterAppRoutes(mux,handlers)
}