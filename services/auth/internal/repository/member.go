package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MemberRepository interface {
	CreateAdminMember(ctx context.Context, tx pgx.Tx, orgID string, userID string) error
}

type memberRepository struct {
	Db *pgxpool.Pool
}

func NewMemberRepository(pool *pgxpool.Pool) MemberRepository {
	return &memberRepository{
		Db: pool,
	}
}

func (r *memberRepository) CreateAdminMember(
	ctx context.Context,
	tx pgx.Tx,
	orgID string,
	userID string,
) error {
	currCtx, cancel := context.WithTimeout(ctx, time.Second*10)

	defer cancel()

	//TODO: Update This To Specific Role
	query := `
		INSERT INTO organization_members(org_id, user_id, role)
		VALUES ($1, $2, 'admin')
	`

	_, err := tx.Exec(currCtx, query, orgID, userID)

	fmt.Println(err)
	return err
}
