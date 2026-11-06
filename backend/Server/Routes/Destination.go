package routes



import (
	handlers "gateway/Server/Handlers"
	"net/http"
)



func RegisterDestinationRoutes(mux *http.ServeMux,handler *handlers.Handlers){

	mux.HandleFunc("/create-destination",handler.DestinationHandler.CreateDestination)

	mux.HandleFunc("/update-destination", handler.DestinationHandler.UpdateDestination)

	mux.HandleFunc("/get-all-destinations", handler.DestinationHandler.GetAllDestinations)
}

