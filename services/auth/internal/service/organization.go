package service

import (
	"context"
	"fmt"
	"time"

	o "os"

	"github.com/iamsainadh/neovault/auth/internal/model"
	"github.com/iamsainadh/neovault/auth/internal/repository"
	"github.com/iamsainadh/neovault/utils"
	vault "github.com/iamsainadh/neovault/vault/proto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrganizationService interface {
	CreateOrganization(
		ctx context.Context,
		input model.CreateOrganizationInput,
	) (*model.CreateOrganizationResponse, error)
}

type organizationService struct {
	db          *pgxpool.Pool
	orgRepo     repository.OrganizationRepository
	memberRepo  repository.MemberRepository
	vaultClient vault.VaultServiceClient
}

func NewOrganizationService(
	db *pgxpool.Pool,
	orgRepo repository.OrganizationRepository,
	memberRepo repository.MemberRepository,
	vaultClient vault.VaultServiceClient,
) OrganizationService {
	return &organizationService{
		db:          db,
		orgRepo:     orgRepo,
		memberRepo:  memberRepo,
		vaultClient: vaultClient,
	}
}

func (os *organizationService) CreateOrganization(
	ctx context.Context,
	input model.CreateOrganizationInput,
) (*model.CreateOrganizationResponse, error) {
	currCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := os.db.Begin(currCtx)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Rollback(currCtx)
	}()

	// Organization
	createdOrg, err := os.orgRepo.CreateOrganization(currCtx, input, tx)

	fmt.Println("Org" + createdOrg.Id)

	if err != nil {
		return nil, err
	}

	// Member
	//TODO : Get The User Member Id and role

	err = os.memberRepo.CreateAdminMember(currCtx, tx, createdOrg.Id, input.UserId)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(currCtx); err != nil {
		return nil, err
	}

	// Create Vault Key
	encryptedKey, err := utils.GenerateEncryptedKEK(o.Getenv("RMK_KEY"))

	if err != nil {
		return nil, err
	}
	_, err = os.vaultClient.CreateVault(currCtx, &vault.CreateVaultRequest{
		OrgId:        createdOrg.Id,
		EncryptedKey: encryptedKey,
		Algorithm:    "AES-256-GCM",
	})

	fmt.Println(err)

	if err != nil {
		return nil, err
	}
	// TODO : Update The Token After The Role is changed
	token, err := utils.GenerateJwt(o.Getenv("JWT_SECRET"), input.Id, "admin", createdOrg.Id)

	createdOrg.Token = token

	return createdOrg, nil
}
