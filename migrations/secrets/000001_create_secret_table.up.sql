CREATE TYPE secret_type AS ENUM (
    'generic',
    'api_key',
    'password',
    'token'
);

CREATE TYPE secret_status AS ENUM (
    'active',
    'archived',
    'deleted'
);


CREATE TABLE IF NOT EXISTS SECRETS(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    org_id UUID NOT NULL REFERENCES organizations(id),
    name VARCHAR(70) NOT NULL,
    description TEXT,
    payload_format TEXT NOT NULL DEFAULT 'SINGLE_VALUE',
    secret_type secret_type NOT NULL,
    current_version INT NOT NULL DEFAULT 1,
    is_rotation_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    rotation_interval_days INT,
    status secret_status NOT NULL DEFAULT 'active',
    deleted_at TIMESTAMPTZ,
    last_rotated_at TIMESTAMPTZ,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    --  CONSTRAINT fk_secrets_org
    --     FOREIGN KEY (org_id)
    --     REFERENCES organizations(id),

    -- CONSTRAINT fk_secrets_created_by
    --     FOREIGN KEY (created_by)
    --     REFERENCES organization_members(id)
);


-- This enforces uniqueness only for non-deleted records.
CREATE UNIQUE INDEX uq_secrets_org_name_active ON secrets (org_id, name) WHERE deleted_at IS NULL;

CREATE INDEX idx_secrets_org_name ON secrets (org_id, name);
