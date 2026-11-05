package handlers

import (
	"encoding/json"
	repo "gateway/Server/Repo"
	services "gateway/Server/Services"
	"net/http"
)




type DestinationHandler struct {
    Services *services.Services
    Repo     *repo.Repo
}

func CreateDestinationHandler(services *services.Services, repo *repo.Repo) *DestinationHandler {
    return &DestinationHandler{
        Services: services,
        Repo:     repo,
    }
}

func (h *DestinationHandler) CreateDestination(w http.ResponseWriter, req *http.Request) {
    var body CreateDestinationRequest

    if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
        WriteResponse(w, http.StatusBadRequest, false, "Invalid request body")
        return
    }

    destinationID, err := h.Repo.DestinationRepo.CreateDestination(h.Repo.Db, body.EndpointID, body.URL, body.Weight)
    if err != nil {
        WriteResponse(w, http.StatusInternalServerError, false, "Failed to create destination")
        return
    }

    if err := h.Repo.EndpointRepo.RegisterDestinationID(h.Repo.Db, body.EndpointID, destinationID); err != nil {
        WriteResponse(w, http.StatusInternalServerError, false, "Failed to register destination")
        return
    }

    WriteResponse(w, http.StatusCreated, true, "Created destination", map[string]any{
        "id": destinationID,
    })
}

func (h *DestinationHandler) UpdateDestination(w http.ResponseWriter, req *http.Request) {
    var body UpdateDestinationRequest

    if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
        WriteResponse(w, http.StatusBadRequest, false, "Invalid request body")
        return
    }

    err := h.Repo.DestinationRepo.UpdateDestination(h.Repo.Db, body.ID, body.EndpointID, body.URL, body.Weight, body.IsActive)
    if err != nil {
        WriteResponse(w, http.StatusInternalServerError, false, "Failed to update destination")
        return
    }

    WriteResponse(w, http.StatusOK, true, "Updated destination", map[string]any{
        "id": body.ID,
    })
}