package repo

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)



type RateLimitingRepo struct {}


func CreateRateLimitingRepo() *RateLimitingRepo {
    return &RateLimitingRepo{}
}

//reoved the clietn chcekign , cahnged 3 to 2 in both loop 
var rateLimitScript = redis.NewScript(`
    local now = tonumber(ARGV[1])
    local buckets = {}

    for i = 1, 2 do
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
            return i + 2
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

        redis.call("HSET", key, "tokens", tokens, "last_refill", now)

        if tokens < 1 then
            return i
        end
    end

    for i = 1, 2 do
        redis.call("HSET", buckets[i].key, "tokens", buckets[i].tokens - 1)
    end

    return 0
`)

func (*RateLimitingRepo)CheckRateLimit(appKey string, endpointKey string, clientKey string, rdb *redis.Client) (int64, error) {
    return rateLimitScript.Run(context.Background(), rdb, []string{appKey, endpointKey, clientKey}, time.Now().UnixMilli()).Int64()
}
// Return value	Meaning
// 0	All buckets allowed
// 1	Application limit exceeded
// 2	Endpoint limit exceeded
// 3	Client limit exceeded //removed for now
// 4	Application bucket missing
// 5	Endpoint bucket missing
// 6	Client bucket missing //removed for now 

// Return	Meaning
// 0	Request allowed
// 1	Application limit exceeded
// 2	Endpoint limit exceeded
// 3	Application bucket missing
// 4	Endpoint bucket missing

func (*RateLimitingRepo) InitializeApplicationBucket(app_name string, token int, limit int, refill_rate float64, rdb *redis.Client) error {
    key := ProjectName + ":" + app_name + ":rate_limit"

    return rdb.HSet(context.Background(), key, map[string]any{
        "tokens":      token,
        "limit":       limit,
        "refill_rate": refill_rate,
        "last_refill": time.Now().UnixMilli(),
    }).Err()
}

func (*RateLimitingRepo) InitializeEndpointBucket(app_name string, endpoint string, token int, limit int, refill_rate float64, rdb *redis.Client) error {
    key := ProjectName + ":" + app_name + ":" + endpoint + ":rate_limit"

    return rdb.HSet(context.Background(), key, map[string]any{
        "tokens":      token,
        "limit":       limit,
        "refill_rate": refill_rate,
        "last_refill": time.Now().UnixMilli(),
    }).Err()
}

func (*RateLimitingRepo) InitializeClientBucket(app_name string, client_ip string, token int, limit int, refill_rate float64, rdb *redis.Client) error {
    key := ProjectName + ":" + app_name + ":" + client_ip + ":rate_limit"

    return rdb.HSet(context.Background(), key, map[string]any{
        "tokens":      token,
        "limit":       limit,
        "refill_rate": refill_rate,
        "last_refill": time.Now().UnixMilli(),
    }).Err()
}