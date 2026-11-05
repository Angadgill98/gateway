CREATE TABLE app (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    endpoint_ids UUID[] DEFAULT '{}'
    rate_limit INTEGER
);