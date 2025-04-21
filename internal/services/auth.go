package services

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
	"pvz/internal/logger"
	"pvz/internal/models"
	"pvz/internal/models/auth"
	"pvz/internal/repositories"
	"pvz/internal/tokens"
	"pvz/pkg/errors"
)

type AuthService interface {
	Login(req auth.LoginRequest) (string, error)
	Register(req auth.RegisterRequest) (auth.RegisterResponse, error)
	DummyLogin(req auth.DummyLoginRequest) (string, error)
}
type authServiceImpl struct {
	userRepo repositories.UserRepository
	conn     *sql.DB
}

func NewAuthService(userRepo repositories.UserRepository, conn *sql.DB) AuthService {
	return &authServiceImpl{
		userRepo: userRepo,
		conn:     conn,
	}
}

func (as *authServiceImpl) Login(req auth.LoginRequest) (string, error) {
	log := logger.Log.With("email", req.Email)
	log.Info("attempting user login")

	user, err := as.userRepo.GetByEmail(as.conn, req.Email)
	if err != nil {
		log.Error("failed to fetch user", "err", err)
		return "", errors.NewInternalError()
	}
	if user == nil {
		log.Warn("user not found")
		return "", errors.NewInvalidCredentials()
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		log.Warn("invalid password")
		return "", errors.NewInvalidCredentials()
	}

	token, err := tokens.GenerateJwt(user.Id, user.Role)
	if err != nil {
		log.Error("failed to generate jwt", "err", err)
		return "", err
	}

	log.Info("user login successful", "user_id", user.Id, "role", user.Role)
	return token, nil
}

func (as *authServiceImpl) Register(req auth.RegisterRequest) (auth.RegisterResponse, error) {
	log := logger.Log.With("email", req.Email)
	log.Info("attempting user registration")

	tx, err := as.conn.Begin()
	if err != nil {
		log.Error("failed to begin transaction", "err", err)
		return auth.RegisterResponse{}, errors.NewInternalError()
	}
	defer tx.Rollback()

	user, err := as.userRepo.GetByEmail(tx, req.Email)
	if err != nil {
		log.Error("failed to check existing user", "err", err)
		return auth.RegisterResponse{}, errors.NewInternalError()
	}
	if user != nil {
		log.Warn("user already exists")
		return auth.RegisterResponse{}, errors.NewObjectAlreadyExists("user", "email", req.Email)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", "err", err)
		return auth.RegisterResponse{}, errors.NewInternalError()
	}

	newUser := models.User{
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		Role:         auth.Role(req.Role),
	}
	userID, err := as.userRepo.Create(tx, newUser)
	if err != nil {
		log.Error("failed to create user", "err", err)
		return auth.RegisterResponse{}, errors.NewInternalError()
	}

	if err := tx.Commit(); err != nil {
		log.Error("failed to commit transaction", "err", err)
		return auth.RegisterResponse{}, errors.NewInternalError()
	}

	log.Info("user registered successfully", "user_id", userID, "role", req.Role)

	return auth.RegisterResponse{
		Id:    userID,
		Email: newUser.Email,
		Role:  req.Role,
	}, nil
}

func (as *authServiceImpl) DummyLogin(req auth.DummyLoginRequest) (string, error) {
	log := logger.Log.With("role", req.Role)
	log.Info("starting dummy login")

	token, err := tokens.GenerateDummyJwt(auth.Role(req.Role))
	if err != nil {
		log.Error("failed to generate dummy jwt", "err", err)
		return "", err
	}

	log.Info("dummy login successful")
	return token, nil
}
