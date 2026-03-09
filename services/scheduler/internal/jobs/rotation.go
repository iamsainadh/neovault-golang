package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	secretPb "github.com/iamsainadh/neovault/secret/proto"
	"github.com/iamsainadh/neovault/utils"
	vaultPb "github.com/iamsainadh/neovault/vault/proto"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type RotationJob struct {
	secretClient        secretPb.SecretServiceClient
	vaultClient         vaultPb.VaultServiceClient
	secretVersionClient secretPb.SecretVersionClient
	Db                  *pgxpool.Pool
}

func NewRotationClient(secretClient secretPb.SecretServiceClient, vault vaultPb.VaultServiceClient, secretVersionClient secretPb.SecretVersionClient, Db *pgxpool.Pool) *RotationJob {
	return &RotationJob{
		secretClient:        secretClient,
		vaultClient:         vault,
		secretVersionClient: secretVersionClient,
		Db:                  Db,
	}
}

func (j *RotationJob) Run(ctx context.Context) {
	//!
	secretResponse, err := j.secretClient.CheckSecretRotationRequired(ctx, &secretPb.CheckSecretRotationRequest{})

	logger.Info(secretResponse)

	if err != nil {
		//! Handle This
		logger.Info(err)
	}

	for _, secret := range secretResponse.Secrets {
		j.rotateSecret(ctx, secret)
	}
}

func (j *RotationJob) rotateSecret(ctx context.Context, secret *secretPb.SecretRotationItem) {
	// Acquire The Single Connection
	conn, lockID, err := j.acquireLock(ctx, secret.Id)
	defer conn.Release()

	if err != nil {
		return
	}

	// Unlock Any error
	defer j.releaseLock(ctx, conn, lockID)

	// Get Secret And Vault
	secretVersion, vault, err := j.fetchSecretData(ctx, secret)

	// Decrypt Secret

	if err != nil {
		return
	}

	// --
	dbConfig, decryptedDEK, err := j.decryptSecret(secretVersion, vault)

	// ------
	if err != nil {
		return
	}
	// Test Connection of Current Db
	if err := j.validateConnection(ctx, dbConfig); err != nil {
		return
	}

	// Generate New Password
	newPassword, err := utils.GeneratePassword(12)

	// Update Connected Db

	lastRotatedAt, err := j.rotatePassword(ctx, dbConfig, newPassword)

	if err != nil {
		return
	}

	if err := j.verifyNewPassword(ctx, dbConfig, newPassword); err != nil {
		return
	}

	/*
		Save New Version With secret_id,version,encrypted_payload, encrypted_dek,vault_key_id,created_by return version and
		Update Old Version Status
	*/

	// Update The Secret Current Version and last updated
	if err := j.updateAndStoreSecret(ctx, dbConfig, newPassword, decryptedDEK, lastRotatedAt, secretVersion, vault, secret); err != nil {
		return
	}
}

type SecretValue struct {
	JsonValue string `json:"json_value"`
}

type DBConfig struct {
	DBUser     string `json:"db_user"`
	DBPassword string `json:"db_password"`
	DBURL      string `json:"db_url"`
}

//-- Acquire Lock

func (j *RotationJob) acquireLock(ctx context.Context, secretId string) (*pgxpool.Conn, int64, error) {
	conn, err := j.Db.Acquire(ctx)

	if err != nil {
		log.Println("failed to acquire connection:", err)
		return nil, 0, err
	}
	lockID := utils.HashUUIDToInt64(secretId)

	var locked bool
	err = conn.QueryRow(ctx, `SELECT pg_try_advisory_lock($1)`, lockID).Scan(&locked)
	if err != nil || !locked {
		conn.Release()
		return nil, 0, err
	}

	return conn, lockID, nil
}

func (j *RotationJob) releaseLock(ctx context.Context, conn *pgxpool.Conn, lockID int64) {
	_, _ = conn.Exec(ctx, `SELECT pg_advisory_unlock($1)`, lockID)
}

// -- Fetch Secret Version and Vault
func (j *RotationJob) fetchSecretData(ctx context.Context, secret *secretPb.SecretRotationItem) (*secretPb.GetSecretVersionBySecretIdResponse, *vaultPb.GetVaultResponse, error) {
	secretVersion, err := j.secretVersionClient.GetSecretVersionById(ctx,
		&secretPb.GetSecretVersionBySecretIdRequest{
			SecretId: secret.Id,
		})

	if err != nil {
		slog.Warn("Secret Version Error", "Secret Version", err)
		return nil, nil, err
	}

	orgVault, err := j.vaultClient.GetVaultByOrgId(ctx, &vaultPb.GetVaultRequest{
		OrgId: secret.OrgId,
	})

	if err != nil {
		slog.Warn("Failed To Get Key:", "Vault Key Error", err)
		return nil, nil, err
	}

	return secretVersion, orgVault, nil
}

// -- Decrypt Secret
func (j *RotationJob) decryptSecret(secretVersion *secretPb.GetSecretVersionBySecretIdResponse, vault *vaultPb.GetVaultResponse) (DBConfig, []byte, error) {
	rmk := utils.ConvertToAESKeyFromRMK(os.Getenv("RMK_KEY"))

	kek, err := utils.DecryptKEK(vault.EncryptedKey, rmk)

	if err != nil {
		return DBConfig{}, nil, err
	}

	decryptedDEK, err := utils.Decrypt(secretVersion.EncryptedDek, kek)

	if err != nil {
		return DBConfig{}, nil, err
	}

	decryptPayload, err := utils.Decrypt(secretVersion.EncryptedPayload, decryptedDEK)

	if err != nil {
		return DBConfig{}, nil, err
	}

	var value SecretValue

	err = json.Unmarshal(decryptPayload, &value)

	if err != nil {
		return DBConfig{}, nil, err
	}

	var dbConfig DBConfig

	err = json.Unmarshal([]byte(value.JsonValue), &dbConfig)

	if err != nil {
		return DBConfig{}, nil, err
	}

	return dbConfig, decryptedDEK, nil
}

// ---- Test Connection

func (j *RotationJob) validateConnection(ctx context.Context, dbConfig DBConfig) error {
	return TestDBConnection(dbConfig.DBURL, dbConfig.DBUser, dbConfig.DBPassword)
}

func (j *RotationJob) rotatePassword(ctx context.Context, dbConfig DBConfig, newPassword string) (time.Time, error) {
	return time.Now(), UpdatePostgresPassword(ctx, dbConfig.DBURL, dbConfig.DBUser, dbConfig.DBPassword, newPassword)
}

func (j *RotationJob) verifyNewPassword(ctx context.Context, dbConfig DBConfig, newPassword string) error {
	dbConfig.DBPassword = newPassword
	return TestDBConnection(dbConfig.DBURL, dbConfig.DBUser, dbConfig.DBPassword)
}

func (j *RotationJob) updateAndStoreSecret(
	ctx context.Context,
	dbConfig DBConfig,
	newPassword string,
	dek []byte,
	LastRotatedAt time.Time,
	secretVersion *secretPb.GetSecretVersionBySecretIdResponse,
	vault *vaultPb.GetVaultResponse,
	secret *secretPb.SecretRotationItem,
) error {

	dbConfig.DBPassword = newPassword

	inner, _ := json.Marshal(dbConfig)

	value := SecretValue{
		JsonValue: string(inner),
	}

	finalPayload, _ := json.Marshal(value)

	encryptedPayload, err := utils.Encrypt(finalPayload, dek)
	if err != nil {
		return err
	}

	// Insert New Version
	newSecretVersion, err := j.secretVersionClient.CreateSecretVersion(ctx, &secretPb.CreateSecretVersionRequest{
		SecretId:         secret.Id,
		Version:          secretVersion.Version + 1,
		VaultKeyId:       vault.Id,
		EncryptedPayload: encryptedPayload,
		EncryptedDek:     secretVersion.EncryptedDek,
		CreatedBy:        "",
	})

	if err != nil {
		return err
	}

	// Update The Secret current version
	_, err = j.secretClient.UpdateSecretRotated(ctx, &secretPb.UpdateSecretRotatedRequest{
		SecretId:       secret.Id,
		CurrentVersion: newSecretVersion.Version,
		LastRotatedAt:  timestamppb.New(LastRotatedAt),
	})

	return nil
}

func TestDBConnection(dbURL, user, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgx.ParseConfig(dbURL)
	if err != nil {
		return err
	}

	config.User = user
	config.Password = password

	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	if err := conn.Ping(ctx); err != nil {
		return err
	}

	return nil
}

// -- Update Postgres
func UpdatePostgresPassword(ctx context.Context, dbURL, user, oldPassword, newPassword string) error {
	config, err := pgx.ParseConfig(dbURL)
	if err != nil {
		return err
	}

	config.User = user
	config.Password = oldPassword

	conn, err := pgx.ConnectConfig(ctx, config)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	query := fmt.Sprintf(`ALTER USER "%s" WITH PASSWORD $1`, user)

	_, err = conn.Exec(ctx, query, newPassword)
	if err != nil {
		return err
	}

	return nil
}
