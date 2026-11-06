package repo

import (
    "context"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
)

type EndpointRepo struct{}

func CreateEndpointRepo() *EndpointRepo {
    return &EndpointRepo{}
}

func (*EndpointRepo) CreateEndpoint(db *pgxpool.Pool, route string, app_id uuid.UUID) (uuid.UUID, error) {
    query, values := BuildInsertQuery("endpoint", map[string]any{
        "route":       route,
        "app_id":      app_id,
    }, []string{"id"})

    var id uuid.UUID
    err := db.QueryRow(context.Background(), query, values...).Scan(&id)

    return id, err
}

func (*EndpointRepo) GetAllEndpointsByAppID(db *pgxpool.Pool, appID uuid.UUID) ([]Endpoint, error) {
    query, values := BuildSelectQuery("endpoint", []string{
        "id", "app_id", "route", "destination_ids", "rate_limit", "refill_rate",
    }, map[string][]any{
        "app_id": {"=", appID},
    })

    rows, err := db.Query(context.Background(), query, values...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    endpoints := []Endpoint{}

    for rows.Next() {
        var endpoint Endpoint

        err := rows.Scan(
            &endpoint.ID,
            &endpoint.AppID,
            &endpoint.Route,
            &endpoint.DestinationIDs,
            &endpoint.RateLimit,
            &endpoint.RefillRate,
        )
        if err != nil {
            return nil, err
        }

        endpoints = append(endpoints, endpoint)
    }

    return endpoints, rows.Err()
}

func (*EndpointRepo) GetEndpoint_AppId(db *pgxpool.Pool, route string, app_id uuid.UUID) (Endpoint, error) {
    query, values := BuildSelectQuery("endpoint", []string{
        "id",
        "app_id",
        "route",
        "destination_ids",
        "rate_limit",
        "refill_rate",
    }, map[string][]any{
        "route":  {"=", route},
        "app_id": {"=", app_id},
    })

    var endpoint Endpoint

    err := db.QueryRow(context.Background(), query, values...).Scan(
        &endpoint.ID,
        &endpoint.AppID,
        &endpoint.Route,
        &endpoint.DestinationIDs,
        &endpoint.RateLimit,
        &endpoint.RefillRate,
    )

    return endpoint, err
}

func (*EndpointRepo) UpdateEndpoint(db *pgxpool.Pool, id uuid.UUID, route string, rate_limiting int, refill_rate float64, app_id uuid.UUID) error {
    query, values := BuildUpdateQuery("endpoint", map[string]any{
        "route":       route,
        "rate_limit":  rate_limiting,
        "refill_rate": refill_rate,
        "app_id":      app_id,
    }, map[string][]any{
        "id": {"=", id},
    })

    _, err := db.Exec(context.Background(), query, values...)

    return err
}


func (*EndpointRepo) RegisterDestinationID(db *pgxpool.Pool, endpointID uuid.UUID, destinationID uuid.UUID) error {
    query := `
        UPDATE endpoint
        SET destination_ids = array_append(destination_ids, $1)
        WHERE id = $2
    `

    _, err := db.Exec(context.Background(), query, destinationID, endpointID)
    return err
}