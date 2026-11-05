package routes

import (
	handlers "gateway/Server/Handlers"
	"net/http"
)



func RegisterAppRoutes(mux *http.ServeMux, handler *handlers.Handlers) {

    mux.HandleFunc("/get-all-apps", handler.AppHandlers.GetAllApps)

    mux.HandleFunc("/register-app", handler.AppHandlers.RegisterApp)

}