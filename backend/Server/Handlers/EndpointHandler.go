package handlers

import (
	"encoding/json"
	repo "gateway/Server/Repo"
	services "gateway/Server/Services"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type EndpointHandler struct {
	Services *services.Services
	Repo     *repo.Repo
}

func CreateEndpointHandler(services *services.Services, repo *repo.Repo) *EndpointHandler {
	return &EndpointHandler{
		Services: services,
		Repo:     repo,
	}
}

func (h *EndpointHandler) RegisterEndpoint(w http.ResponseWriter, req *http.Request) {
	var body RegisterEndpointRequest

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		log.Printf("Invalid register endpoint request body: %v", err)
		WriteResponse(w, http.StatusBadRequest, false, "Invalid request body")
		return
	}

	endpointID, err := h.Repo.EndpointRepo.CreateEndpoint(h.Repo.Db, body.Route, body.AppID)
	if err != nil {
		log.Printf("Failed to create endpoint: %v", err)
		WriteResponse(w, http.StatusInternalServerError, false, "Failed to create endpoint")
		return
	}

	if err := h.Repo.AppRepo.RegisterEndpoint(h.Repo.Db, body.AppID, endpointID); err != nil {
		log.Printf("Failed to register endpoint in app: %v", err)
		WriteResponse(w, http.StatusInternalServerError, false, "Failed to register endpoint")
		return
	}

	WriteResponse(w, http.StatusCreated, true, "Created the endpoint", map[string]any{
		"id": endpointID,
	})
}

func (h *EndpointHandler) UpdateEndpoint(w http.ResponseWriter, req *http.Request) {
	var body UpdateEndpointRequest

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		log.Printf("Invalid update endpoint request body: %v", err)
		WriteResponse(w, http.StatusBadRequest, false, "Invalid request body")
		return
	}

	err := h.Repo.EndpointRepo.UpdateEndpoint(h.Repo.Db, body.ID, body.Route, body.RateLimit, body.RefillRate, body.AppID)
	if err != nil {
		log.Printf("Failed to update endpoint: %v", err)
		WriteResponse(w, http.StatusInternalServerError, false, "Failed to update endpoint")
		return
	}

	WriteResponse(w, http.StatusOK, true, "Updated the endpoint", map[string]any{
		"id": body.ID,
	})
}

func (h *EndpointHandler) GetAllEndpoints(w http.ResponseWriter, req *http.Request) {
    appIDStr := req.URL.Query().Get("app_id")

    appID, err := uuid.Parse(appIDStr)
    if err != nil {
        log.Printf("Invalid app ID: %v", err)
        WriteResponse(w, http.StatusBadRequest, false, "Invalid app ID")
        return
    }

    endpoints, err := h.Repo.EndpointRepo.GetAllEndpointsByAppID(h.Repo.Db, appID)
    if err != nil {
        log.Printf("Failed to retrieve endpoints: %v", err)
        WriteResponse(w, http.StatusInternalServerError, false, "Failed to retrieve endpoints")
        return
    }

    WriteResponse(w, http.StatusOK, true, "Retrieved endpoints successfully", endpoints)
}