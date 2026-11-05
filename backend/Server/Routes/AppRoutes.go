package routes

import (
	handlers "gateway/Server/Handlers"
	"net/http"
)



func RegisterAppRoutes(mux *http.ServeMux,handler *handlers.Handlers){
	app_mux:=http.NewServeMux()

	app_mux.HandleFunc("/register-app",handler.AppHandlers.RegisterApp)

	mux.Handle("/",http.StripPrefix("/",app_mux))
}