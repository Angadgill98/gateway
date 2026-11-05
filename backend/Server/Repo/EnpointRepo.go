package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)


type EndpointRepo struct{

}


func CreateEndpointRepo()*EndpointRepo{
	return &EndpointRepo{
	}
}


func (*EndpointRepo) CreateEndpoint(db *pgxpool.Pool, route string, url string, rate_limiting int, app_id uuid.UUID) (uuid.UUID, error) {
    query, values := BuildInsertQuery("endpoint", map[string]any{
        "route":      route,
        "url":        url,
        "rate_limit": rate_limiting,
        "app_id":     app_id,
    }, []string{"id"})

    var id uuid.UUID

    err := db.QueryRow(context.Background(), query, values...).Scan(&id)

    return id, err
}

func (*EndpointRepo) GetEndpoint_AppId(db *pgxpool.Pool, route string, url string) (uuid.UUID, uuid.UUID, error) {
    query, values := BuildSelectQuery("endpoint", []string{"id", "app_id"}, map[string][]any{
        "route": {"=", route},
        "url":   {"=", url},
    })

    var id uuid.UUID
    var app_id uuid.UUID

    err := db.QueryRow(context.Background(), query, values...).Scan(&id, &app_id)

    return id, app_id, err
}