package repo

import (
    "context"

    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgxpool"
)

type DestinationRepo struct {
}

func CreateDestinationRepo() *DestinationRepo {
    return &DestinationRepo{}
}

func (*DestinationRepo) CreateDestination(db *pgxpool.Pool, endpointID uuid.UUID, url string, weight int) (uuid.UUID, error) {
    query, values := BuildInsertQuery("endpoint_destination", map[string]any{
        "endpoint_id": endpointID,
        "url":         url,
        "weight":      weight,
    }, []string{"id"})

    var id uuid.UUID
    err := db.QueryRow(context.Background(), query, values...).Scan(&id)

    return id, err
}

func (*DestinationRepo) UpdateDestination(db *pgxpool.Pool, id uuid.UUID, endpointID uuid.UUID, url string, weight int, isActive bool) error {
    query, values := BuildUpdateQuery("endpoint_destination", map[string]any{
        "endpoint_id": endpointID,
        "url":         url,
        "weight":      weight,
        "is_active":   isActive,
    }, map[string][]any{
        "id": {"=", id},
    })

    _, err := db.Exec(context.Background(), query, values...)

    return err
}