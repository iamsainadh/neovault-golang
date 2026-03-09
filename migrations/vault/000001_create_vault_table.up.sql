CREATE TYPE VAULT_STATUS AS ENUM (
    'active',
    'rotated',
    'retired',
    'compromised'
);

CREATE TABLE IF NOT EXISTS VAULT_KEYS(
    id uuid NOT NULL PRIMARY KEY DEFAULT uuidv7(),
    org_id UUID NOT NULL REFERENCES organizations(id),
    key_version int NOT NULL DEFAULT 1,
    encrypted_key bytea NOT NULL,
    algorithm varchar(20) NOT NULL,
    status VAULT_STATUS NOT NULL DEFAULT 'active',
    rotated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE UNIQUE INDEX IF NOT EXISTS uq_vault_org_id_key_version ON VAULT_KEYS (org_id, key_version) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_vault_org_id_key_version ON VAULT_KEYS (org_id, key_version,status);
