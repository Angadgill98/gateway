package handlers

import (
	"encoding/json"
	repo "gateway/Server/Repo"
	services "gateway/Server/Services"
	"log"
	"net/http"

	"github.com/google/uuid"
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
        log.Printf("Invalid create destination request body: %v", err)
        WriteResponse(w, http.StatusBadRequest, false, "Invalid request body")
        return
    }

    destinationID, err := h.Repo.DestinationRepo.CreateDestination(h.Repo.Db, body.EndpointID, body.URL, body.Weight)
    if err != nil {
        log.Printf("Failed to create destination: %v", err)
        WriteResponse(w, http.StatusInternalServerError, false, "Failed to create destination")
        return
    }

    if err := h.Repo.EndpointRepo.RegisterDestinationID(h.Repo.Db, body.EndpointID, destinationID); err != nil {
        log.Printf("Failed to register destination ID in endpoint: %v", err)
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
        log.Printf("Invalid update destination request body: %v", err)
        WriteResponse(w, http.StatusBadRequest, false, "Invalid request body")
        return
    }

    err := h.Repo.DestinationRepo.UpdateDestination(h.Repo.Db, body.ID, body.EndpointID, body.URL, body.Weight, body.IsActive)
    if err != nil {
        log.Printf("Failed to update destination in PostgreSQL: %v", err)
        WriteResponse(w, http.StatusInternalServerError, false, "Failed to update destination")
        return
    }

    destinationKey := h.Services.ProxyService.CreateDestinationHashKey(body.AppName, body.Route)

    err = h.Repo.DestinationRepo.UpdateDestinationActive_RDB(
        h.Repo.Rdb,
        destinationKey,
        body.ID,
        body.IsActive,
    )
    if err != nil {
        log.Printf("Failed to update destination in Redis: %v", err)
        WriteResponse(w, http.StatusInternalServerError, false, "Failed to update destination in Redis")
        return
    }

    WriteResponse(w, http.StatusOK, true, "Updated destination", map[string]any{
        "id": body.ID,
    })
}

func (h *DestinationHandler) GetAllDestinations(w http.ResponseWriter, req *http.Request) {
    endpointIDStr := req.URL.Query().Get("endpoint_id")

    endpointID, err := uuid.Parse(endpointIDStr)
    if err != nil {
        log.Printf("Invalid endpoint ID: %v", err)
        WriteResponse(w, http.StatusBadRequest, false, "Invalid endpoint ID")
        return
    }

    destinations, err := h.Repo.DestinationRepo.GetAllDestinationsByEndpointID(h.Repo.Db, endpointID)
    if err != nil {
        log.Printf("Failed to retrieve destinations: %v", err)
        WriteResponse(w, http.StatusInternalServerError, false, "Failed to retrieve destinations")
        return
    }

    WriteResponse(w, http.StatusOK, true, "Retrieved destinations successfully", destinations)
}