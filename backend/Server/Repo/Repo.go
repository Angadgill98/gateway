package repo

import (
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var ProjectName="gateway"


type Repo struct{
	AppRepo *AppRepo
	EndpointRepo *EndpointRepo
    DestinationRepo *DestinationRepo
    RateLimitingRepo *RateLimitingRepo
	Rdb *redis.Client
	Db *pgxpool.Pool
}

func CreateRepo(db *pgxpool.Pool,rdb *redis.Client)*Repo{
	return &Repo{
		AppRepo: CreateApprepo(),
		EndpointRepo: CreateEndpointRepo(),
        DestinationRepo: CreateDestinationRepo(),
        RateLimitingRepo:CreateRateLimitingRepo(),
		Db: db,
		Rdb: rdb,
	}
}

func BuildInsertQuery(table string, data map[string]any, returning []string) (string, []any) {
    columns := make([]string, 0, len(data))
    placeholders := make([]string, 0, len(data))
    values := make([]any, 0, len(data))

    i := 1

    for column, value := range data {
        columns = append(columns, column)
        placeholders = append(placeholders, fmt.Sprintf("$%d", i))
        values = append(values, value)
        i++
    }

    query := fmt.Sprintf(
        "INSERT INTO %s (%s) VALUES (%s) RETURNING %s",
        table,
        strings.Join(columns, ", "),
        strings.Join(placeholders, ", "),
        strings.Join(returning, ", "),
    )

    return query, values
}


func BuildUpdateQuery(table string, data map[string]any, conditions map[string][]any) (string, []any) {
    sets := make([]string, 0, len(data))
    where := make([]string, 0, len(conditions))
    values := make([]any, 0, len(data)+len(conditions))

    i := 1

    for column, value := range data {
        sets = append(sets, fmt.Sprintf("%s = $%d", column, i))
        values = append(values, value)
        i++
    }

    for column, condition := range conditions {
        where = append(where, fmt.Sprintf("%s %v $%d", column, condition[0], i))
        values = append(values, condition[1])
        i++
    }

    query := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
        table,
        strings.Join(sets, ", "),
        strings.Join(where, " AND "),
    )

    return query, values
}

func BuildSelectQuery(table string, columns []string, conditions map[string][]any) (string, []any) {
    query := "SELECT " + strings.Join(columns, ", ") + " FROM " + table

    values := make([]any, 0, len(conditions))
    where := make([]string, 0, len(conditions))

    i := 1

    for column, condition := range conditions {
        where = append(where, fmt.Sprintf("%s %v $%d", column, condition[0], i))
        values = append(values, condition[1])
        i++
    }

    if len(where) > 0 {
        query += " WHERE " + strings.Join(where, " AND ")
    }

    return query, values
}