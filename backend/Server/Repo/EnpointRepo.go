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

func (*EndpointRepo) CreateEndpoint(db *pgxpool.Pool, route string, rate_limiting int, refill_rate float64, app_id uuid.UUID) (uuid.UUID, error) {
    query, values := BuildInsertQuery("endpoint", map[string]any{
        "route":       route,
        "rate_limit":  rate_limiting,
        "refill_rate": refill_rate,
        "app_id":      app_id,
    }, []string{"id"})

    var id uuid.UUID
    err := db.QueryRow(context.Background(), query, values...).Scan(&id)

    return id, err
}

func (*EndpointRepo) GetEndpoint_AppId(db *pgxpool.Pool, route string, app_id uuid.UUID) (uuid.UUID, uuid.UUID, error) {
    query, values := BuildSelectQuery("endpoint", []string{"id", "app_id"}, map[string][]any{
        "route":  {"=", route},
        "app_id": {"=", app_id},
    })

    var id uuid.UUID
    var appID uuid.UUID

    err := db.QueryRow(context.Background(), query, values...).Scan(&id, &appID)

    return id, appID, err
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