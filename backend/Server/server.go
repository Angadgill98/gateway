package server

import (
	routes "gateway/Server/Routes"
	"net/http"

	"github.com/rs/cors"
)




type Server struct {
    httpServer *http.Server
    mux        *http.ServeMux
}

func Initialize(addr string) *Server {
    mux := http.NewServeMux()

    return &Server{
        mux: mux,
        httpServer: &http.Server{
            Addr:    addr,
            Handler: mux,
        },
    }
}


func (server *Server) StartServer() {
    routes.RegisterRoutes(server.mux)

    corsHandler := cors.AllowAll()

    server.httpServer.Handler = corsHandler.Handler(server.mux)

    server.httpServer.ListenAndServe()
}


