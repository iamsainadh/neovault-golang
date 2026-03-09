CREATE TABLE IF NOT EXISTS audit_events (
    id UUID PRIMARY KEY,

    organization_id UUID NOT NULL,

    actor_type VARCHAR(30) NOT NULL,
    actor_id UUID,
    actor_membership_id UUID,

    action VARCHAR(100) NOT NULL,

    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID,
    resource_name VARCHAR(255),

    service VARCHAR(50) NOT NULL,
    failure_reason TEXT,
    request_id VARCHAR(100),
    trace_id VARCHAR(100),
    session_id VARCHAR(100),
    ip_address INET,
    user_agent TEXT,

    method VARCHAR(10),
    endpoint TEXT,

    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

    event_time TIMESTAMPTZ NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_org_time ON audit_events (organization_id, event_time DESC);

CREATE INDEX idx_audit_actor ON audit_events (organization_id, actor_id, event_time DESC);

CREATE INDEX idx_audit_request
ON audit_events (request_id);

-- Allow Only Insertion , block update or delete 

CREATE OR REPLACE FUNCTION prevent_audit_events_modification()
RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'audit_events is immutable. UPDATE and DELETE are not allowed';
END;
$$ LANGUAGE plpgsql;


DROP TRIGGER IF EXISTS audit_events_no_update ON audit_events;

CREATE TRIGGER audit_events_no_update BEFORE UPDATE ON audit_events FOR EACH ROW EXECUTE FUNCTION prevent_audit_events_modification();


DROP TRIGGER IF EXISTS audit_events_no_delete ON audit_events;

CREATE TRIGGER audit_events_no_delete BEFORE DELETE ON audit_events FOR EACH ROW EXECUTE FUNCTION prevent_audit_events_modification();


CREATE OR REPLACE FUNCTION prevent_audit_events_truncate() RETURNS TRIGGER AS $$
BEGIN
    RAISE EXCEPTION 'audit_events cannot be truncated';
END;
$$ LANGUAGE plpgsql;


DROP TRIGGER IF EXISTS audit_events_no_truncate ON audit_events;

CREATE TRIGGER audit_events_no_truncate BEFORE TRUNCATE ON audit_events FOR EACH STATEMENT EXECUTE FUNCTION prevent_audit_events_truncate();