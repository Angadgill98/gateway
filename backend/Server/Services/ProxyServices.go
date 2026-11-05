package services

import (
	"io"
	"net/http"
)



type ProxyService struct {

}

func CreateProxyService() *ProxyService{
	return &ProxyService{}
}

type ProxyRequestInfo struct {
    Method  string
    Host    string
    Path    string
    Query   string
    Headers http.Header
    Body    []byte
    RemoteAddr string
}

func ExtractReqInfo(r *http.Request) (*ProxyRequestInfo, error) {
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
        Body:       body,
        RemoteAddr: r.RemoteAddr,
    }, nil
}


type ProxyResponseInfo struct {
    StatusCode int
    Headers    http.Header
    Body       []byte
}

func ExtractProxyRes(res *http.Response) (*ProxyResponseInfo, error) {
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

