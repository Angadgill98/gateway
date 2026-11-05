package server

import "net/http"




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


func (server *Server) StartServer(){
	server.httpServer.ListenAndServe()
}


