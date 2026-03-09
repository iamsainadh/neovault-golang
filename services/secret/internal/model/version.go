package model

type CreateSecretVersionInput struct {
	Secret_id         string
	Version           int32
	OrgId             string
	Value             string
	Encrypted_payload []byte
	Encrypted_dek     []byte
	Vault_key_id      string
	Created_by        string
}

type CreateSecretVersionResponse struct {
	Id         string
	Secret_id  string
	Version    int32
	Created_by string
}

type GetSecretVersionBySecretIdResponse struct {
	Id               string
	SecretId         string
	Version          int32
	EncryptedPayload []byte
	EncryptedDek     []byte
	VaultKeyId       string
}
