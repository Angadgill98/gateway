package handlers

import (
	repo "gateway/Server/Repo"
	services "gateway/Server/Services"
	"net/http"
)





type EndpointHandler struct {
	Services *services.Services
	Repo *repo.Repo
}

func CreateEndpointHandler(services *services.Services, repo *repo.Repo) *EndpointHandler{
	return &EndpointHandler{
		Services: services,
		Repo: repo,
	}
}



func (h *EndpointHandler) RegisterEndpoint(w http.ResponseWriter, req *http.Request){
	h.Repo.EndpointRepo.CreateEndpoint()
	h.Repo.AppRepo.RegisterEndpoint()
}