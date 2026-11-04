package routes

import "net/http"



func RegisterAppRoutes(mux *http.ServeMux){
	app_mux:=http.NewServeMux()

	app_mux.HandleFunc("/register-app",)




	mux.Handle("/",http.StripPrefix("/",app_mux))
}