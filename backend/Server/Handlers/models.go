package handlers

import "github.com/google/uuid"

type RegisterAppRequest struct {
    Name       string  `json:"name"`
}

type RegisterEndpointRequest struct {
    AppID      uuid.UUID `json:"app_id"`
    Route      string    `json:"route"`
    URL        string    `json:"url"`
    RateLimit  int       `json:"rate_limit"`
    RefillRate float64   `json:"refill_rate"`
}

type UpdateEndpointRequest struct {
    ID         uuid.UUID `json:"id"`
    AppID      uuid.UUID `json:"app_id"`
    Route      string    `json:"route"`
    URL        string    `json:"url"`
    RateLimit  int       `json:"rate_limit"`
    RefillRate float64   `json:"refill_rate"`
}