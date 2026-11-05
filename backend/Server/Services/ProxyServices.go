package services

import (
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
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

