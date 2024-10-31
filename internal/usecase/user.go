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
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type userRepository interface {
	SaveUser(ctx context.Context, user *models.User) (int64, error)
	User(ctx context.Context, email string) (models.User, error)
}

type passwordHasher interface {
	Generate(password string) (string, error)
	Compare(hashedPassword string, password string) error
}

type tokenManager interface {
	CreateUserAccessToken(user *models.User) (string, error)
	CreateUserRefreshToken(user *models.User) (string, error)
}

type UserUsecase struct {
	passwordHasher passwordHasher
	userRepository userRepository
	tokenManager   tokenManager
	logger         *slog.Logger
}

func NewUserUsecase(
	passwordHasher passwordHasher,
	userRepository userRepository,
	tokenManager tokenManager,
	logger *slog.Logger,
) *UserUsecase {
	return &UserUsecase{
		passwordHasher: passwordHasher,
		userRepository: userRepository,
		tokenManager:   tokenManager,
		logger:         logger,
	}
}

// CreateUser creates a new user and stores in into the repository.
// Returns an id of the created user and error.
func (u *UserUsecase) CreateUser(ctx context.Context, email, password, name string) (int64, error) {
	const op = "usecase.CreateUser"

	log := u.logger.With(slog.String("op", op))

	hashedPassword, err := u.passwordHasher.Generate(password)
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

// Login checks if user exists and checks if the password is correct.
// Returns access/refresh token or error.
func (u *UserUsecase) Login(ctx context.Context, email, password string) (models.Tokens, error) {
	const op = "usecase.Login"

	log := u.logger.With(slog.String("op", op))

	user, err := u.userRepository.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Info("user not found")

			return models.Tokens{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		log.Error("failed to get user from repository", sl.Err(err))

		return models.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	err = u.passwordHasher.Compare(user.HashedPassword, password)
	if err != nil {
		log.Info("incorrect password")

		return models.Tokens{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	tokens, err := u.getTokens(&user)
	if err != nil {
		log.Error("failed to create user tokens", sl.Err(err))

		return models.Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user successfully logged in")

	return tokens, nil
}

// Refresh validates the refresh token and provides a new access token.
func (u *UserUsecase) Refresh(ctx context.Context, refreshToken string) (models.Tokens, error) {
	// TODO: implement
	return models.Tokens{}, nil
}

func (u *UserUsecase) getTokens(user *models.User) (models.Tokens, error) {
	accessToken, err := u.tokenManager.CreateUserAccessToken(user)
	if err != nil {
		return models.Tokens{}, err
	}

	refreshToken, err := u.tokenManager.CreateUserRefreshToken(user)
	if err != nil {
		return models.Tokens{}, err
	}

	return models.Tokens{
		Access:  accessToken,
		Refresh: refreshToken,
	}, nil
}
