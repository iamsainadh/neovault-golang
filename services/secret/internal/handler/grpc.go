package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/iamsainadh/neovault/secret/internal/app_error"
	"github.com/iamsainadh/neovault/secret/internal/mapper"
	"github.com/iamsainadh/neovault/secret/internal/model"
	"github.com/iamsainadh/neovault/secret/internal/service"
	pb "github.com/iamsainadh/neovault/secret/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO :-- Create Version Handler

type SecretHandler struct {
	pb.UnimplementedSecretServiceServer // embedded by value — standard gRPC-Go pattern
	secretService                       service.SecretService
	logger                              *slog.Logger
}

type SecretVersionHandler struct {
	pb.UnimplementedSecretVersionServer
	secretVersionService service.SecretVersionService
	logger               *slog.Logger
}

func NewSecretHandler(secretService service.SecretService, logger *slog.Logger) *SecretHandler {
	return &SecretHandler{
		secretService: secretService,
		logger:        logger,
	}
}

func NewSecretVersionHandler(secretVersionService service.SecretVersionService, logger *slog.Logger) *SecretVersionHandler {
	return &SecretVersionHandler{
		secretVersionService: secretVersionService,
		logger:               logger,
	}
}

func (h *SecretHandler) CreateSecret(ctx context.Context, req *pb.CreateSecretRequest) (*pb.CreateSecretResponse, error) {

	fmt.Println(req)
	//-- Edge Cases Handling
	if req.OrgId == "" || req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "org_id and name are required")

	}

	secret := mapper.CreateRequestToInput(req)

	created, err := h.secretService.CreateSecretService(ctx, secret)

	if err != nil {
		fmt.Printf("Error Here:%v", err)

		h.logger.ErrorContext(ctx, "CreateSecret failed",
			slog.String("org_id", req.OrgId),
			slog.String("name", req.Name),
			slog.Any("error", err),
		)
		return nil, app_error.ToGRPCError(err)
	}

	h.logger.InfoContext(ctx, "secret created",
		slog.String("secret_id", created.ID),
		slog.String("org_id", created.OrgID),
	)

	return &pb.CreateSecretResponse{
		Secret: mapper.SecretToProto(created),
	}, nil
}

func (h *SecretVersionHandler) CreateSecretVersion(ctx context.Context, req *pb.CreateSecretVersionRequest) (*pb.CreateSecretVersionResponse, error) {
	if req.SecretId == "" {
		return nil, status.Error(codes.InvalidArgument, "secret_id is required")
	}

	secretVersion := model.CreateSecretVersionInput{
		Secret_id:         req.SecretId,
		Version:           req.Version,
		Vault_key_id:      req.VaultKeyId,
		Encrypted_payload: req.EncryptedPayload,
		Encrypted_dek:     req.EncryptedDek,
		Created_by:        req.CreatedBy,
	}

	createdVersion, err := h.secretVersionService.CreateSecretVersion(ctx, secretVersion)

	if err != nil {

		return nil, app_error.ToGRPCError(err)
	}

	return &pb.CreateSecretVersionResponse{
		Id:        createdVersion.Id,
		SecretId:  createdVersion.Secret_id,
		Version:   createdVersion.Version,
		CreatedBy: createdVersion.Created_by,
	}, nil
}

func (h *SecretHandler) CheckSecretRotationRequired(ctx context.Context, req *pb.CheckSecretRotationRequest) (*pb.CheckSecretRotationResponse, error) {
	secret := mapper.CheckSecretRequiredToModel(req)

	secretsRepo, err := h.secretService.CheckSecretRotationScheduler(ctx, secret)

	if err != nil {
		return nil, errors.New("An internal error")
	}

	pbSecrets := make([]*pb.SecretRotationItem, 0, len(secretsRepo.Secrets))

	for i, s := range secretsRepo.Secrets {
		pbSecrets[i] = &pb.SecretRotationItem{
			Id:             s.Id,
			OrgId:          s.OrgID,
			Name:           s.Name,
			CurrentVersion: s.CurrentVersion,
		}
	}

	return &pb.CheckSecretRotationResponse{
		Secrets: pbSecrets,
	}, nil
}

func (h *SecretVersionHandler) GetSecretVersionById(ctx context.Context, req *pb.GetSecretVersionBySecretIdRequest) (*pb.GetSecretVersionBySecretIdResponse, error) {
	if req.SecretId == "" {
		return nil, status.Error(codes.InvalidArgument, "secret_id is required")
	}

	secretVersion, err := h.secretVersionService.GetSecretVersionBySecretId(ctx, req.SecretId)

	if err != nil {
		return nil, errors.New("An internal error")
	}

	return &pb.GetSecretVersionBySecretIdResponse{
		Id:               secretVersion.Id,
		SecretId:         secretVersion.SecretId,
		Version:          secretVersion.Version,
		EncryptedPayload: secretVersion.EncryptedPayload,
		EncryptedDek:     secretVersion.EncryptedDek,
		VaultKeyId:       secretVersion.VaultKeyId,
	}, nil
}

func (h *SecretHandler) UpdateSecretRotated(ctx context.Context, req *pb.UpdateSecretRotatedRequest) (*pb.UpdateSecretRotatedResponse, error) {
	updateSecret, err := h.secretService.UpdateSecretRotated(ctx, req.CurrentVersion, req.LastRotatedAt.AsTime())

	if err != nil {
		return nil, errors.New("An internal error")
	}

	return &pb.UpdateSecretRotatedResponse{
		Updated: updateSecret != "",
	}, nil
}
