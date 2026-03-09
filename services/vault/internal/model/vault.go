package model

import "time"

type VaultKeyStatus string

const (
	Active      VaultKeyStatus = "active"
	Rotated     VaultKeyStatus = "rotated"
	Retired     VaultKeyStatus = "retired"
	Compromised VaultKeyStatus = "compromised"
)

type VaultKey struct {
	ID           string         `json:"id"`
	OrgId        string         `json:"org_id"`
	KeyVersion   int32          `json:"key_version"`
	EncryptedKey []byte         `json:"encrypted_key"`
	Algorithm    string         `json:"algorithm"`
	Status       VaultKeyStatus `json:"status"`
	RotatedAt    time.Time      `json:"rotated_at"`
	DeletedAt    time.Time      `json:"deleted_at"`
	CreatedAt    time.Time      `json:"created_at"`
}

type CreateVaultInputs struct {
	OrgId        string
	KeyVersion   int32
	EncryptedKey []byte
	Algorithm    string
}

type GetVaultByOrgInputs struct {
	OrgId string
}

type GetVaultByOrgResponse struct {
	Id            string
	Version       int32
	Encrypted_key []byte
}
