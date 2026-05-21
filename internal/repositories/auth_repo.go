package repositories

import (
	"context"
	"go_emmie/internal/database"
	"go_emmie/internal/types"
	"go_emmie/prisma/db"
)

type UserRepository interface {
	Create(ctx context.Context, dto types.RegisterRequestDTO, passwordHash string) (*db.UserModel, error)
	FindByEmail(ctx context.Context, email string) (*db.UserModel, error)
}

 

type userRepo struct {
	db *database.DB
}

// Create implements [UserRepository].
func (u *userRepo) Create(ctx context.Context, dto types.RegisterRequestDTO, passwordHash string) (*db.UserModel, error) {

	user, err := u.db.User.CreateOne(
		db.User.Email.Set(dto.Email),
		db.User.PasswordHash.Set(passwordHash),
		db.User.FirstName.Set(dto.FirstName),
		db.User.LastName.Set(dto.LastName),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindByEmail implements [UserRepository].
func (u *userRepo) FindByEmail(ctx context.Context, email string) (*db.UserModel, error) {
	user, err := u.db.User.FindFirst(
		db.User.Email.Equals(email),
	).Exec(ctx)

	if err != nil {
		return nil, err
	}

	return user, nil

}

func NewUserRepository(db *database.DB) UserRepository {
	return &userRepo{db: db}
}
