package handlers

import (
	repo "gateway/Server/Repo"
	services "gateway/Server/Services"
	"net/http"
)





type ProxyHandler struct {
	Services *services.Services
	Repo *repo.Repo
}

func CreateProxyHandler(services *services.Services, repo *repo.Repo) *ProxyHandler{
	return &ProxyHandler{
		Services: services,
		Repo: repo,
	}
}



func (h *ProxyHandler) RegisterEndpoint(w http.ResponseWriter, req *http.Request){
	h.Repo.EndpointRepo.CreateEndpoint()
	h.Repo.AppRepo.RegisterEndpoint()
}

