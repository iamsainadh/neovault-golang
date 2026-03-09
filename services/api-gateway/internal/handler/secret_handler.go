package handler

import (
	secret "github.com/iamsainadh/neovault/secret/proto"
)

type SecretHandler struct {
	client secret.SecretServiceClient
}

func NewSecretHandler(client secret.SecretServiceClient) *SecretHandler {
	return &SecretHandler{
		client: client,
	}
}

//! Create Secret
