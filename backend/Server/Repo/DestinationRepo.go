package repo

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type DestinationRepo struct {
}

func CreateDestinationRepo() *DestinationRepo {
    return &DestinationRepo{}
}

func (*DestinationRepo) CreateDestination(db *pgxpool.Pool, endpointID uuid.UUID, url string, weight int) (uuid.UUID, error) {
    query, values := BuildInsertQuery("destination", map[string]any{
        "endpoint_id": endpointID,
        "url":         url,
        "weight":      weight,
        "is_active":   false,
    }, []string{"id"})

    var id uuid.UUID
    err := db.QueryRow(context.Background(), query, values...).Scan(&id)

    return id, err
}

func (*DestinationRepo) UpdateDestination(db *pgxpool.Pool, id uuid.UUID, endpointID uuid.UUID, url string, weight int, isActive bool) error {
    query, values := BuildUpdateQuery("destination", map[string]any{
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

func (*DestinationRepo) GetAllDestinationsByEndpointID(db *pgxpool.Pool, endpointID uuid.UUID) ([]Destination, error) {
    query, values := BuildSelectQuery("destination", []string{
        "id",
        "endpoint_id",
        "url",
        "weight",
        "is_active",
    }, map[string][]any{
        "endpoint_id": {"=", endpointID},
    })

    rows, err := db.Query(context.Background(), query, values...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    destinations := []Destination{}

    for rows.Next() {
        var destination Destination

        err := rows.Scan(
            &destination.ID,
            &destination.EndpointID,
            &destination.URL,
            &destination.Weight,
            &destination.IsActive,
        )
        if err != nil {
            return nil, err
        }

        destinations = append(destinations, destination)
    }

    return destinations, rows.Err()
}

func (*DestinationRepo) InitializeDestination_RDB(rdb *redis.Client, destination_key string, destinations []Destination) error {
    values := make(map[string]any)

    for _, destination := range destinations {
        data, err := json.Marshal(map[string]any{
            "id":             destination.ID,
            "endpoint_id":    destination.EndpointID,
            "url":            destination.URL,
            "weight":         destination.Weight,
            "is_active":      destination.IsActive,
            "current_weight": 0,
        })
        if err != nil {
            return err
        }

        values[destination.ID.String()] = string(data)
    }

    if len(values) == 0 {
        return nil
    }

    return rdb.HSet(context.Background(), destination_key, values).Err()
}


var selectDestinationScript = redis.NewScript(`
    local key = KEYS[1]

    if redis.call("EXISTS", key) == 0 then
        return {0, ""}
    end

    local entries = redis.call("HGETALL", key)

    local destinations = {}
    local totalWeight = 0

    for i = 1, #entries, 2 do
        local id = entries[i]
        local destination = cjson.decode(entries[i + 1])

        local weight = tonumber(destination.weight) or 0
        local currentWeight = tonumber(destination.current_weight) or 0

        if destination.is_active == true and weight > 0 then
            table.insert(destinations, {
                id = id,
                data = destination,
                current_weight = currentWeight + weight
            })

            totalWeight = totalWeight + weight
        end
    end

    if #destinations == 0 then
        return {1, ""}
    end

    table.sort(destinations, function(a, b)
        return a.id < b.id
    end)

    local selected = nil

    for i = 1, #destinations do
        if selected == nil or destinations[i].current_weight > selected.current_weight then
            selected = destinations[i]
        end
    end

    for i = 1, #destinations do
        local destination = destinations[i]

        if destination.id == selected.id then
            destination.current_weight = destination.current_weight - totalWeight
        end

        destination.data.current_weight = destination.current_weight

        redis.call("HSET", key, destination.id, cjson.encode(destination.data))
    end

    selected.data.current_weight = nil

    return {1, cjson.encode(selected.data)}
`)
func (*DestinationRepo) GetProxyReqDestination_RDM(rdb *redis.Client, destination_key string) (Destination, bool, error) {
    var destination Destination

    result, err := selectDestinationScript.Run(context.Background(), rdb, []string{destination_key}).Result()
    if err != nil {
        return destination, false, err
    }

    values := result.([]interface{})

    exists := values[0].(int64) == 1
    if !exists {
        return destination, false, nil
    }

    data := values[1].(string)
    if data == "" {
        return destination, true, errors.New("no active destinations available")
    }

    if err := json.Unmarshal([]byte(data), &destination); err != nil {
        return destination, true, err
    }

    return destination, true, nil
}

func (*DestinationRepo) UpdateDestinationActive_RDB(rdb *redis.Client, destination_key string, destinationID uuid.UUID, isActive bool) error {
    ctx := context.Background()

    data, err := rdb.HGet(ctx, destination_key, destinationID.String()).Result()
    if err == redis.Nil {
        return nil
    }
    if err != nil {
        return err
    }

    var destination map[string]any

    if err := json.Unmarshal([]byte(data), &destination); err != nil {
        return err
    }

    destination["is_active"] = isActive

    updatedData, err := json.Marshal(destination)
    if err != nil {
        return err
    }

    return rdb.HSet(ctx, destination_key, destinationID.String(), string(updatedData)).Err()
}