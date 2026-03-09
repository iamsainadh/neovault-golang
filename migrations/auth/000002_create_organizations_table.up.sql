CREATE Type organization_status AS ENUM(
    'active',
    'suspended',
    'deleted'
);

CREATE TABLE IF NOT EXISTS ORGANIZATIONS(
    id uuid NOT NULL PRIMARY KEY DEFAULT uuidv7(),
    name varchar(70) NOT NULL,
    status organization_status NOT NULL DEFAULT 'active',
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_name_status ON organizations (name,status);