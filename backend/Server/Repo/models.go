package repo

import "github.com/google/uuid"

type App struct {
    ID          uuid.UUID   `json:"id"`
    Name        string      `json:"name"`
    EndpointIDs []uuid.UUID `json:"endpoint_ids"`
    RateLimit   int         `json:"rate_limit"`
    RefillRate  float64     `json:"refill_rate"`
}

type Endpoint struct {
    ID             uuid.UUID   `json:"id"`
    AppID          uuid.UUID   `json:"app_id"`
    Route          string      `json:"route"`
    DestinationIDs []uuid.UUID `json:"destination_ids"`
    RateLimit      int         `json:"rate_limit"`
    RefillRate     float64     `json:"refill_rate"`
}

type Destination struct {
    ID         uuid.UUID `json:"id"`
    EndpointID uuid.UUID `json:"endpoint_id"`
    URL        string    `json:"url"`
    Weight     int       `json:"weight"`
    IsActive   bool      `json:"is_active"`
}