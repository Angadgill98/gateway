package services

import (
	"errors"
	"fmt"
	repo "gateway/Server/Repo"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/google/uuid"
)



type ProxyService struct {

}

func CreateProxyService() *ProxyService{
	return &ProxyService{}
}

type ProxyRequestInfo struct {
    Method     string
    Host       string
    Path       string
    Query      string
    Headers    http.Header
    Cookies    []*http.Cookie
    Body       []byte
    RemoteAddr string
}

func (*ProxyService)ExtractReqInfo(r *http.Request) (*ProxyRequestInfo, error) {
    var body []byte

    if r.Body != nil {
        var err error
        body, err = io.ReadAll(r.Body)
        if err != nil {
            return nil, err
        }
    }

    return &ProxyRequestInfo{
        Method:     r.Method,
        Host:       r.Host,
        Path:       r.URL.Path,
        Query:      r.URL.RawQuery,
        Headers:    r.Header,
        Cookies:    r.Cookies(),
        Body:       body,
        RemoteAddr: r.RemoteAddr,
    }, nil
}


type ProxyResponseInfo struct {
    StatusCode int
    Headers    http.Header
    Body       []byte
}

func (s *ProxyService) ForwardRequest(w http.ResponseWriter, req *http.Request, destination string) {
    target, err := url.Parse(destination)
    if err != nil {
        http.Error(w, "Invalid destination URL", http.StatusBadRequest)
        return
    }

    proxy := &httputil.ReverseProxy{
        Rewrite: func(pr *httputil.ProxyRequest) {
            pr.SetURL(target)
        },

        ModifyResponse: func(res *http.Response) error {
            // TODO: Execute response callback here.
            // You can inspect or modify the response before sending it.

            return nil
        },

        ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
            http.Error(w, "Proxy request failed", http.StatusBadGateway)
        },
    }

    proxy.ServeHTTP(w, req)
}

func (*ProxyService)ExtractProxyRes(res *http.Response) (*ProxyResponseInfo, error) {
    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, err
    }

    return &ProxyResponseInfo{
        StatusCode: res.StatusCode,
        Headers:    res.Header,
        Body:       body,
    }, nil
}

func (s *ProxyService) GetReqAndClientMetaD(info *ProxyRequestInfo) (string, string, string, error) {
    parts := strings.SplitN(strings.Trim(info.Path, "/"), "/", 2)

    if len(parts) != 2 {
        return "", "", "", fmt.Errorf("invalid request path")
    }

    appName := parts[0]
    endpoint := "/" + parts[1]

    clientIP, _, err := net.SplitHostPort(info.RemoteAddr)
    if err != nil {
        clientIP = info.RemoteAddr
    }

    return appName, endpoint, clientIP, nil
}

func (s *ProxyService) CreateAppBucketKey(app_name string) string {
    return repo.ProjectName + ":rate_limit:app:" + app_name
}

func (s *ProxyService) CreateEndpointBucketKey(app_name string, endpoint string) string {
    return repo.ProjectName + ":rate_limit:endpoint:" + app_name + ":" + endpoint
}

func (s *ProxyService) CreateClientBucketKey(app_name string, client_ip string) string {
    return repo.ProjectName + ":rate_limit:client:" + app_name + ":" + client_ip
}

func (s *ProxyService) CreateDestinationHashKey(appName string, endpoint string) string {
    return repo.ProjectName + ":app:" + appName + ":endpoint:" + endpoint + ":destinations"
}


func (s *ProxyService) HandleRateLimitResult(result int64, app_name string, route string, repo *repo.Repo) (bool, string) {
    switch result {
    case 0:
        return true, "Request allowed"

    case 1:
        return false, "Application rate limit exceeded"

    case 2:
        return false, "Endpoint rate limit exceeded"
    }

    // Handle missing buckets.
    app, err := repo.AppRepo.GetAppByName(repo.Db, app_name)
    if err != nil {
        return false, "Failed to retrieve application"
    }

    switch result {
    case 3:
        err = repo.RateLimitingRepo.InitializeApplicationBucket(app_name,app.RateLimit-1, app.RateLimit, app.RefillRate, repo.Rdb)
        if err != nil {
            return false, "Failed to initialize application bucket"
        }

        return true, "Application bucket initialized"

    case 4:
        endpoint, err := repo.EndpointRepo.GetEndpoint_AppId(repo.Db, route, app.ID)
        if err != nil {
            return false, "Failed to retrieve endpoint"
        }

        err = repo.RateLimitingRepo.InitializeEndpointBucket(app_name, route, app.RateLimit-1, endpoint.RateLimit, endpoint.RefillRate, repo.Rdb)
        if err != nil {
            return false, "Failed to initialize endpoint bucket"
        }

        return true , "Endpoint bucket initialized"

    default:
        return false, "Unknown rate limit result"
    }
}


func (h *ProxyService) GetDestination(destinations []repo.Destination, destinationKey string) (repo.Destination, error) {
    var destination repo.Destination

    for i := range destinations {
        if !destinations[i].IsActive || destinations[i].Weight <= 0 {
            continue
        }

        if destination.ID == uuid.Nil || destinations[i].Weight > destination.Weight {
            destination = destinations[i]
        }
    }

    if destination.ID == uuid.Nil {
        return destination, errors.New("no active destinations available")
    }

    return destination, nil
}