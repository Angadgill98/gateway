package handlers

import (
	"encoding/json"
	db "gateway/Server/Db"
	repo "gateway/Server/Repo"
	services "gateway/Server/Services"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)



type Handlers struct{
	AppHandlers *AppHandler
    EndpointHandler *EndpointHandler

    Services *services.Services
    Repo *repo.Repo
    Db *pgxpool.Pool
}

func CreateHandlers()*Handlers{
	Service:=services.CreateService()
    Rdb:=db.InitRedis()
    Db:=db.InitDB()
    Repo:=repo.CreateRepo(Db,Rdb)
	return &Handlers{
		AppHandlers: CreateAppHandler(Service,Repo),
        EndpointHandler: CreateEndpointHandler(Service,Repo),
        Services: Service,
        Repo: Repo,
        Db: Db,
	}
}

type Response struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Data    any    `json:"data,omitempty"`
}

func WriteResponse(w http.ResponseWriter, status int, success bool, message string, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

    json.NewEncoder(w).Encode(Response{
        Success: success,
        Message: message,
        Data:    data,
    })
}