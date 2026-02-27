CREATE TABLE IF NOT EXISTS tenants (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS usage_logs (
    id UUID PRIMARY KEY,
    api_key TEXT,
    path TEXT,
    ts TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS api_keys (
    id UUID PRIMARY KEY,
    tenant_id UUID REFERENCES tenants(id),
    key TEXT UNIQUE NOT NULL
);