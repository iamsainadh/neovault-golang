package model

import "time"

type SecretType int32

const (
	SecretTypeGeneric  SecretType = 0
	SecretTypeAPIKey   SecretType = 1
	SecretTypePassword SecretType = 2
	SecretTypeToken    SecretType = 3
)

type SecretStatus int32

const (
	SecretStatusActive   SecretStatus = 0
	SecretStatusArchived SecretStatus = 1
	SecretStatusDeleted  SecretStatus = 2
)

type Secret struct {
	ID                   string
	OrgID                string
	Name                 string
	SecretType           SecretType
	Description          string
	CurrentVersion       int32
	IsRotationEnabled    bool
	RotationIntervalDays int32
	Status               SecretStatus
	LastRotatedAt        time.Time
	CreatedBy            string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	PayloadFormat        string
}

type SecretValue struct {
	Raw  *string
	JSON *string // validated JSON string
}

type CreateSecretInput struct {
	OrgID                string
	Name                 string
	Description          string
	Value                SecretValue
	SecretType           SecretType
	IsRotationEnabled    bool
	RotationIntervalDays int32
	CreatedBy            string
}

// ------------------------------------- Check Secret Required
type CheckSecretRotationInput struct {
	Limit int32  // Optional
	OrgID string // Optional
}

type CheckSecretRotationResponse struct {
	Secrets []*SecretRotationItem
}

type SecretRotationItem struct {
	Id             string
	OrgID          string
	Name           string
	CurrentVersion int32
}

type UpdateSecretInput struct {
	SecretId int
}