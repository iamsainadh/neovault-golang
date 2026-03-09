package mapper

import (
	"strings"

	"github.com/iamsainadh/neovault/secret/internal/model"
	pb "github.com/iamsainadh/neovault/secret/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// proto to model conversion
func CreateRequestToInput(req *pb.CreateSecretRequest) model.CreateSecretInput {
	var val model.SecretValue

	if req.Value != nil {
		switch v := req.Value.Value.(type) {

		case *pb.SecretPayload_RawValue:
			if v != nil {
				val.Raw = &v.RawValue
			}

		case *pb.SecretPayload_JsonValue:
			if v != nil && v.JsonValue != "" {
				val.JSON = &v.JsonValue
			}
		}
	}

	return model.CreateSecretInput{
		OrgID:                req.OrgId,
		Name:                 strings.TrimSpace(req.Name),
		Description:          strings.TrimSpace(req.Description),
		Value:                val,
		SecretType:           protoSecretTypeToModel(req.SecretType),
		IsRotationEnabled:    req.IsRotationEnabled,
		RotationIntervalDays: req.RotationIntervalDays,
		CreatedBy:            req.CreatedBy,
	}
}

// Model to proto
func SecretToProto(s model.Secret) *pb.Secret {
	return &pb.Secret{
		Id:                   s.ID,
		OrgId:                s.OrgID,
		Name:                 s.Name,
		SecretType:           modelSecretTypeToProto(s.SecretType),
		CurrentVersion:       s.CurrentVersion,
		IsRotationEnabled:    s.IsRotationEnabled,
		RotationIntervalDays: s.RotationIntervalDays,
		Status:               modelSecretStatusToProto(s.Status),
		LastRotatedAt:        timestamppb.New(s.LastRotatedAt),
		CreatedAt:            timestamppb.New(s.CreatedAt),
		UpdatedAt:            timestamppb.New(s.UpdatedAt),
		CreatedBy:            s.CreatedBy,
		PayloadFormat:        s.PayloadFormat,
	}
}

func protoSecretTypeToModel(t pb.SecretType) model.SecretType {
	switch t {
	case pb.SecretType_SECRET_TYPE_GENERIC:
		return model.SecretTypeGeneric
	case pb.SecretType_SECRET_TYPE_API_KEY:
		return model.SecretTypeAPIKey
	case pb.SecretType_SECRET_TYPE_PASSWORD:
		return model.SecretTypePassword
	case pb.SecretType_SECRET_TYPE_TOKEN:
		return model.SecretTypeToken
	default:
		return model.SecretTypeGeneric
	}
}

func modelSecretTypeToProto(t model.SecretType) pb.SecretType {
	switch t {
	case model.SecretTypeGeneric:
		return pb.SecretType_SECRET_TYPE_GENERIC
	case model.SecretTypeAPIKey:
		return pb.SecretType_SECRET_TYPE_API_KEY
	case model.SecretTypePassword:
		return pb.SecretType_SECRET_TYPE_PASSWORD
	case model.SecretTypeToken:
		return pb.SecretType_SECRET_TYPE_TOKEN
	default:
		return pb.SecretType_SECRET_TYPE_UNSPECIFIED
	}
}

// Secret Status

func modelSecretStatusToProto(s model.SecretStatus) pb.SecretStatus {
	switch s {
	case model.SecretStatusActive:
		return pb.SecretStatus_SECRET_STATUS_ACTIVE
	case model.SecretStatusArchived:
		return pb.SecretStatus_SECRET_STATUS_ARCHIVED
	case model.SecretStatusDeleted:
		return pb.SecretStatus_SECRET_STATUS_DELETED
	default:
		return pb.SecretStatus_SECRET_STATUS_UNSPECIFIED
	}
}

//---- Check Secret

func CheckSecretRequiredToModel(req *pb.CheckSecretRotationRequest) model.CheckSecretRotationInput {
	return model.CheckSecretRotationInput{
		Limit: req.Limit,
		OrgID: req.OrgId,
	}
}
