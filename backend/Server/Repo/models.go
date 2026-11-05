package repo

import "github.com/google/uuid"

type App struct {
    ID          uuid.UUID   `json:"id"`
    Name        string      `json:"name"`
    EndpointIDs []uuid.UUID `json:"endpoint_ids"`
    RateLimit   int         `json:"rate_limit"`
}