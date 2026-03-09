CREATE TYPE SECRET_VERSION_STATUS AS ENUM (
    'active',
    'deprecated',
    'deleted'
);

CREATE TABLE IF NOT EXISTS SECRET_VERSIONS(
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    secret_id UUID NOT NULL,
    version INT NOT NULL,
    encrypted_payload BYTEA NOT NULL,
    encrypted_dek BYTEA NOT NULL,
    vault_key_id UUID NOT NULL,
    status SECRET_VERSION_STATUS NOT NULL DEFAULT 'active',
    created_by UUID,
    deleted_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),


    CONSTRAINT fk_secrets_id FOREIGN KEY (secret_id) REFERENCES SECRETS(id),
    CONSTRAINT uq_secret_version UNIQUE(secret_id, version)

);

