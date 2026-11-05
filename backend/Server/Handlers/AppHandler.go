package handlers

import (
	"encoding/json"
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



func (h *AppHandler) RegisterApp(w http.ResponseWriter, req *http.Request) {
    var body RegisterAppRequest

    if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
        WriteResponse(w, http.StatusBadRequest, false, "Invalid request body")
        return
    }

    id, err := h.Repo.AppRepo.CreateApp(h.Repo.Db, body.Name)

    if err != nil {
        WriteResponse(w, http.StatusInternalServerError, false, "Failed to create app")
        return
    }

    WriteResponse(w, http.StatusCreated, true, "Created the app", map[string]any{
        "id": id,
    })
}


func (h *AppHandler) GetAllApps(w http.ResponseWriter, req *http.Request) {
    apps, err := h.Repo.AppRepo.GetAllApps(h.Repo.Db)

    if err != nil {
        WriteResponse(w, http.StatusInternalServerError, false, "Failed to retrieve apps")
        return
    }

    WriteResponse(w, http.StatusOK, true, "Apps retrieved successfully", apps)
}