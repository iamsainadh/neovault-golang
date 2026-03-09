PROTO_DIR=proto
SECRET_OUT_DIR=services/secret/proto
VAULT_OUT_DIR=services/vault/proto
AUTH_OUT_DIR=services/auth/proto



proto-auth:
	protoc -I $(PROTO_DIR) \
	--go_out=$(AUTH_OUT_DIR) \
	--go_opt=paths=source_relative \
	--go-grpc_out=$(AUTH_OUT_DIR) \
	--go-grpc_opt=paths=source_relative \
	$(PROTO_DIR)/auth/auth.proto

proto-secret:
	protoc -I $(PROTO_DIR) \
	--go_out=$(SECRET_OUT_DIR) \
	--go_opt=paths=source_relative \
	--go-grpc_out=$(SECRET_OUT_DIR) \
	--go-grpc_opt=paths=source_relative \
	$(PROTO_DIR)/secrets/secrets.proto

proto-vault:
	protoc -I $(PROTO_DIR) \
	--go_out=$(VAULT_OUT_DIR) \
	--go_opt=paths=source_relative \
	--go-grpc_out=$(VAULT_OUT_DIR) \
	--go-grpc_opt=paths=source_relative \
	$(PROTO_DIR)/vault/vault.proto

proto-secret-version:
	protoc -I $(PROTO_DIR) \
	--go_out=$(SECRET_OUT_DIR) \
	--go_opt=paths=source_relative \
	--go-grpc_out=$(SECRET_OUT_DIR) \
	--go-grpc_opt=paths=source_relative \
	$(PROTO_DIR)/secrets/secrets_version.proto