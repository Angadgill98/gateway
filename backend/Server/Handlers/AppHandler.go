package handlers

import (
	repo "gateway/Server/Repo"
	services "gateway/Server/Services"
	"net/http"
)





type AppHandler struct {
	Services *services.Services
	Repo *repo.Repo
}

func CreateAppHandler(services *services.Services, repo *repo.Repo) *AppHandler{
	return &AppHandler{
		Services: services,
		Repo: repo,
	}
}



func (h *AppHandler) RegisterApp(w http.ResponseWriter, req *http.Request){
	h.Repo.AppRepo.CreateApp(h.Repo.Db,)
}