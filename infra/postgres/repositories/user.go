package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jaeyoung0509/compound-interest/domain/shared"
	"github.com/jaeyoung0509/compound-interest/domain/user"
	"github.com/jaeyoung0509/compound-interest/infra/postgres/sqlc/generated"
)

type UserRepository struct {
	queries *generated.Queries
}

func NewUserRepository(db generated.DBTX) *UserRepository {
	return &UserRepository{queries: generated.New(db)}
}

func (r *UserRepository) Get(ctx context.Context, id user.ID) (*user.User, error) {
	row, err := r.queries.GetUser(ctx, id.Value().String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return toUserDomain(row)
}

func (r *UserRepository) Save(ctx context.Context, u *user.User) error {
	return r.queries.InsertUser(ctx, generated.InsertUserParams{
		ID:        u.ID().Value().String(),
		Name:      u.Name(),
		CreatedAt: pgtype.Timestamptz{Time: u.CreatedAt(), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: u.UpdatedAt(), Valid: true},
	})
}

func toUserDomain(row generated.User) (*user.User, error) {
	id, err := shared.ParseID(row.ID)
	if err != nil {
		return nil, err
	}
	return user.Reconstitute(user.IDFromShared(id), row.Name, row.CreatedAt.Time, row.UpdatedAt.Time)
}

// Ensure UserRepository implements user.Repository interface
var _ user.Repository = (*UserRepository)(nil)
