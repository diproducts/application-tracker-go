package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/diproducts/application-tracker-go/internal/domain/models"
	"github.com/diproducts/application-tracker-go/internal/lib/logger/sl"
	"github.com/diproducts/application-tracker-go/internal/repository/storage"
	"log/slog"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type userRepository interface {
	SaveUser(ctx context.Context, user *models.User) (int64, error)
	User(ctx context.Context, email string) (models.User, error)
}

type passwordHasher interface {
	Generate(ctx context.Context, password string) (string, error)
	Compare(ctx context.Context, password string, hashedPassword string) (bool, error)
}

type UserUsecase struct {
	passwordHasher passwordHasher
	userRepository userRepository
	logger         *slog.Logger
}

func NewUserUsecase(
	passwordHasher passwordHasher,
	userRepository userRepository,
	logger *slog.Logger,
) *UserUsecase {
	return &UserUsecase{
		passwordHasher: passwordHasher,
		userRepository: userRepository,
		logger:         logger,
	}
}

// CreateUser creates a new user and stores in into the repository.
// Returns an id of the created user and error.
func (u *UserUsecase) CreateUser(ctx context.Context, email, password, name string) (int64, error) {
	const op = "usecase.CreateUser"

	log := u.logger.With(slog.String("op", op))

	hashedPassword, err := u.passwordHasher.Generate(ctx, password)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	user := models.User{
		Email:          email,
		HashedPassword: hashedPassword,
		Name:           name,
	}

	userId, err := u.userRepository.SaveUser(ctx, &user)
	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			log.Warn("user already exists", sl.Err(err))

			return 0, fmt.Errorf("%s: %w", op, ErrUserAlreadyExists)
		}

		log.Error("failed to save user", sl.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userId, nil
}
