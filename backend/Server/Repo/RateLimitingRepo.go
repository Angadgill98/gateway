package repo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)



type RateLimitingRepo struct {}


func CreateRateLimitingRepo() *RateLimitingRepo {
    return &RateLimitingRepo{}
}

var rateLimitScript = redis.NewScript(`
    local now = tonumber(ARGV[1])
    local buckets = {}

    -- Phase 1: Refill and check all three buckets.
    for i = 1, 3 do
        local key = KEYS[i]

        local data = redis.call(
            "HMGET",
            key,
            "tokens",
            "limit",
            "refill_rate",
            "last_refill"
        )

        if not data[1] or not data[2] or not data[3] or not data[4] then
            return redis.error_reply("Rate limit bucket not initialized")
        end

        local tokens = tonumber(data[1])
        local limit = tonumber(data[2])
        local refill_rate = tonumber(data[3])
        local last_refill = tonumber(data[4])

        local elapsed = math.max(0, now - last_refill) / 1000

        tokens = math.min(limit, tokens + (elapsed * refill_rate))

        buckets[i] = {
            key = key,
            tokens = tokens
        }

        -- Save the refilled balance.
        redis.call("HSET", key, "tokens", tokens, "last_refill", now)

        -- Stop at the first bucket that cannot allow the request.
        if tokens < 1 then
            return i
        end
    end

    -- Phase 2: All buckets passed. Consume one token from each.
    for i = 1, 3 do
        redis.call("HSET", buckets[i].key, "tokens", buckets[i].tokens - 1)
    end

    return 0
`)

func CheckRateLimit(appKey string, endpointKey string, clientKey string, rdb *redis.Client) (int64, error) {
    return rateLimitScript.Run(context.Background(), rdb, []string{appKey, endpointKey, clientKey}, time.Now().UnixMilli()).Int64()
}
// Return   Value	Meaning
// 0		All three buckets allowed the request
// 1		Application limit exceeded
// 2		Endpoint limit exceeded
// 3		Client limit exceeded

func InitializeApplicationBucket(application_id uuid.UUID, limit int, refill_rate float64, rdb *redis.Client) error {
    key := "rate_limit:app:" + application_id.String()

    return rdb.HSet(context.Background(), key, map[string]any{
        "tokens":      limit,
        "limit":       limit,
        "refill_rate": refill_rate,
        "last_refill": time.Now().UnixMilli(),
    }).Err()
}