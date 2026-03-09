CREATE TABLE IF NOT EXISTS USERS(
    id uuid NOT NULL PRIMARY KEY DEFAULT uuidv7(),
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    mfa_enabled boolean DEFAULT false,
    is_active boolean DEFAULT true,
    is_verified boolean DEFAULT false,
    deleted_at TIMESTAMPTZ,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_name_email ON USERS (name,email);