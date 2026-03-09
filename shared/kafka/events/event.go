package events

import "time"

type Event struct {
	ID                string         `json:"id"`
	OrgID             string         `json:"org_id"`
	ActorUserID       string         `json:"actor_user_id"`
	ActorMembershipID string         `json:"actor_membership_id"`
	Action            string         `json:"action"`
	ResourceType      string         `json:"resource_type"`
	ResourceID        string         `json:"resource_id"`
	ServiceName       string         `json:"service_name"`
	Status            string         `json:"status"`
	Severity          string         `json:"severity"`
	RequestID         string         `json:"request_id"`
	TraceID           string         `json:"trace_id"`
	IPAddress         string         `json:"ip_address"`
	UserAgent         string         `json:"user_agent"`
	Payload           map[string]any `json:"payload,omitempty"`
	EventTime         time.Time      `json:"event_time"`
}

/*
AUTH Service
USER_LOGIN (success/failure)
TOKEN_ISSUED / REVOKED

Authorization Changes

ROLE_CREATED / UPDATED / DELETED
PERMISSION_GRANTED / REVOKED
USER_ADDED_TO_ORG



Secret Access (your most important logs)
SECRET_CREATED
SECRET_READ
SECRET_UPDATED
SECRET_DELETED
SECRET_VERSION_ROLLBACK


Secret Rotation
ROTATION_STARTED
ROTATION_SUCCESS
ROTATION_FAILED


Access Control Violations
ACCESS_DENIED
INVALID_TOKEN_USED
UNAUTHORIZED_SECRET_ACCESS

Critical System Changes
ENCRYPTION_KEY_ROTATED
CONFIG_CHANGED (only security-relevant)
*/
