CREATE TABLE endpoint (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_id UUID NOT NULL REFERENCES app(id) ON DELETE CASCADE,
    route VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    rate_limit INTEGER
); 