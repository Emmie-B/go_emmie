package services

import (
	"context"
	"errors"
	"fmt"
	"go_emmie/internal/repositories"
	"go_emmie/internal/types"
	"go_emmie/prisma/db"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	RegisterUser(ctx context.Context, dto types.RegisterRequestDTO) (*types.UserResponseDTO, error)
}

var (
	ErrValidation         = errors.New("validation error")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type authService struct {
	repo repositories.UserRepository
}

func NewAuthService(repo repositories.UserRepository) AuthService {
	return &authService{repo: repo}
}

func (s *authService) RegisterUser(ctx context.Context, dto types.RegisterRequestDTO) (*types.UserResponseDTO, error) {
	// Validate the incoming DTO
	if err := dto.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrValidation, err)
	}

	// check if user already exists
	existingUser, err := s.repo.FindByEmail(ctx, dto.Email)
	if err != nil {
		// Treat Prisma's ErrNotFound as "user does not exist" rather than an error
		if db.IsErrNotFound(err) {
			existingUser = nil
		} else {
			return nil, err
		}
	}

	if existingUser != nil {
		return nil, fmt.Errorf("%w: %s", ErrEmailAlreadyExists, "a user with this email already exists")
	}

	// Hash the password (this is a placeholder, implement proper hashing)

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed hashing password: %w", err)
	}

	// Create the user in the database
	user, err := s.repo.Create(ctx, dto, string(passwordHash))
	if err != nil {
		if info, ok := db.IsErrUniqueConstraint(err); ok {
			if len(info.Fields) == 0 || info.Fields[0] == db.User.Email.Field() {
				return nil, fmt.Errorf("%w: email already exists", ErrEmailAlreadyExists)
			}
		}
		return nil, err
	}

	// Convert to response DTO
	responseDTO := ToUserResponseDTO(user)

	return &responseDTO, nil
}

// ToUserResponseDTO converts a raw Prisma User model into a safe, client-facing DTO
func ToUserResponseDTO(user *db.UserModel) types.UserResponseDTO {
	role := string(user.Role)
	if role == "" {
		role = "CUSTOMER"
	}

	return types.UserResponseDTO{
		ID:         user.ID,
		Email:      user.Email,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Role:       role,
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt,
	}
}
