package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)



type AppRepo struct {
}


func CreateApprepo()*AppRepo{
	return &AppRepo{
	}
}


func (*AppRepo) CreateApp(db *pgxpool.Pool, name string) (uuid.UUID, error) {
    data := map[string]any{
        "name": name,
    }

    query, values := BuildInsertQuery("app", data, []string{"id"})

    var id uuid.UUID

    err := db.QueryRow(context.Background(), query, values...).Scan(&id)

    return id, err
}

func (*AppRepo) RegisterEndpoint(db *pgxpool.Pool, appID uuid.UUID, endpointID uuid.UUID) error {
    query, values := BuildUpdateQuery(
        "app",
        map[string]any{
            "endpoint_id": endpointID,
        },
        map[string][]any{
            "id": {"=", appID},
        },
    )

    _, err := db.Exec(context.Background(), query, values...)

    return err
}


func (*AppRepo) GetAppById(db *pgxpool.Pool, appID uuid.UUID) (App, error) {
    query, values := BuildSelectQuery("app", []string{"id", "name", "endpoint_ids", "rate_limit"}, map[string][]any{
        "id": {"=", appID},
    })

    var app App

    err := db.QueryRow(context.Background(), query, values...).Scan(&app.ID, &app.Name, &app.EndpointIDs, &app.RateLimit)

    return app, err
}

func (*AppRepo) GetAppByName(db *pgxpool.Pool, appName string) (App, error) {
    query, values := BuildSelectQuery("app", []string{"id", "name", "endpoint_ids", "rate_limit"}, map[string][]any{
        "name": {"=", appName},
    })

    var app App

    err := db.QueryRow(context.Background(), query, values...).Scan(&app.ID, &app.Name, &app.EndpointIDs, &app.RateLimit)

    return app, err
}

func (*AppRepo) GetAllApps(db *pgxpool.Pool) ([]App, error) {
    query, values := BuildSelectQuery("app", []string{"id", "name", "endpoint_ids", "rate_limit", "refill_rate"}, nil)

    rows, err := db.Query(context.Background(), query, values...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    apps := make([]App, 0)

    for rows.Next() {
        var app App

        err := rows.Scan(&app.ID, &app.Name, &app.EndpointIDs, &app.RateLimit, &app.RefillRate)
        if err != nil {
            return nil, err
        }

        apps = append(apps, app)
    }

    return apps, rows.Err()
}