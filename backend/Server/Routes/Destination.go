package routes



import (
	handlers "gateway/Server/Handlers"
	"net/http"
)



func RegisterDestinationRoutes(mux *http.ServeMux,handler *handlers.Handlers){
	app_mux:=http.NewServeMux()

	app_mux.HandleFunc("/create-destination",handler.DestinationHandler.CreateDestination)

	app_mux.HandleFunc("/update-destination", handler.DestinationHandler.UpdateDestination)

	mux.Handle("/api/", http.StripPrefix("/api", app_mux))
}

