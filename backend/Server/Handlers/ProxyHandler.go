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


func (h *ProxyHandler) ProxyReqProcessing(w http.ResponseWriter, req *http.Request) {
    reqInfo, err := h.Services.ProxyService.ExtractReqInfo(req)
    if err != nil {
        WriteResponse(w, http.StatusBadRequest, false, "Failed to extract request information")
        return
    }

    appName, endpoint, clientIP, err := h.Services.ProxyService.GetReqAndClientMetaD(reqInfo)
    if err != nil {
        WriteResponse(w, http.StatusBadRequest, false, err.Error())
        return
    }

    appKey := h.Services.ProxyService.CreateAppBucketKey(appName)
    endpointKey := h.Services.ProxyService.CreateEndpointBucketKey(appName, endpoint)
    clientKey := h.Services.ProxyService.CreateClientBucketKey(appName, clientIP)

    result, err := h.Repo.RateLimitingRepo.CheckRateLimit(appKey, endpointKey, clientKey, h.Repo.Rdb)
    if err != nil {
        WriteResponse(w, http.StatusInternalServerError, false, "Rate limit check failed")
        return
    }

    allowed, message := h.Services.ProxyService.HandleRateLimitResult(result, appName, endpoint, h.Repo)

	if !allowed {
		WriteResponse(w, http.StatusTooManyRequests, false, message)
		return
	}

    destinationKey := h.Services.ProxyService.CreateDestinationHashKey(appName,endpoint)
	
	destination, exists, err := h.Repo.DestinationRepo.GetProxyReqDestination_RDM(h.Repo.Rdb, destinationKey)
	if err != nil {
		WriteResponse(w, http.StatusInternalServerError, false, err.Error())
		return
	}

	if !exists {
		app, err := h.Repo.AppRepo.GetAppByName(h.Repo.Db, appName)
		if err != nil {
			WriteResponse(w, http.StatusInternalServerError, false, "Failed to retrieve application")
			return
		}

		endpointData, err := h.Repo.EndpointRepo.GetEndpoint_AppId(h.Repo.Db, endpoint, app.ID)
		if err != nil {
			WriteResponse(w, http.StatusInternalServerError, false, "Failed to retrieve endpoint")
			return
		}

		destinations, err := h.Repo.DestinationRepo.GetDestinationsByEndpointID(h.Repo.Db, endpointData.ID)
		if err != nil {
			WriteResponse(w, http.StatusInternalServerError, false, "Failed to retrieve destinations")
			return
		}

		if len(destinations) == 0 {
			WriteResponse(w, http.StatusServiceUnavailable, false, "No destinations configured")
			return
		}

		destination, err = h.Services.ProxyService.GetDestination(destinations, destinationKey)
		if err != nil {
			WriteResponse(w, http.StatusServiceUnavailable, false, err.Error())
			return
		}
		
		err = h.Repo.DestinationRepo.InitializeDestination_RDB(h.Repo.Rdb, destinationKey, destinations)
		if err != nil {
			WriteResponse(w, http.StatusInternalServerError, false, "Failed to initialize destinations")
			return
		}
	}

	_ = destination // Continue processing later.

    // Continue forwarding the request.
    h.Services.ProxyService.ForwardRequest(w, req, destination.URL)
}