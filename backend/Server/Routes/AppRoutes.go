package routes

import (
	handlers "gateway/Server/Handlers"
	"net/http"
)



func RegisterAppRoutes(mux *http.ServeMux,handler *handlers.Handlers){
	app_mux:=http.NewServeMux()

	app_mux.HandleFunc("/get-all-apps",handler.AppHandlers.GetAllApps)

	app_mux.HandleFunc("/register-app",handler.AppHandlers.RegisterApp)

	mux.Handle("/api/", http.StripPrefix("/api", app_mux))
}