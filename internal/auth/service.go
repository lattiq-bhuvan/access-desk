package auth

import (
	"strconv"
	"context"
	"time"

	fjwt "github.com/lattiq/foundry/auth/jwt"
	ferr "github.com/lattiq/foundry/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/lattiq-bhuvan/access-desk/internal/dto"
	"github.com/lattiq-bhuvan/access-desk/internal/model"
	"github.com/lattiq-bhuvan/access-desk/internal/store"
)

type Service struct {
	st           *store.Store
	tokens       fjwt.TokenService
	accessTTL    time.Duration
}

func NewAuthService(st *store.Store, tokens fjwt.TokenService, accessTTL time.Duration) *Service {
	return &Service{st: st, tokens: tokens, accessTTL: accessTTL}
}

func (s *Service) Login(ctx context.Context, in dto.LoginBody) (*dto.TokenResponse, error) {
	u, err := s.st.UserByEmail(ctx, in.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)) != nil {
		// same error whether user missing or password wrong — don't leak which
		return nil, ferr.ErrUnauthorized("invalid email or password")
	}
	return s.issue(ctx, u)
}

func (s *Service) Refresh(ctx context.Context, in dto.RefreshBody) (*dto.TokenResponse, error) {
	pair, err := s.tokens.RefreshToken(ctx, in.RefreshToken, ClaimsFactory)
	if err != nil {
		return nil, ferr.ErrUnauthorized("invalid refresh token")
	}
	return s.pairToResponse(pair), nil
}

func (s *Service) Me(ctx context.Context, userID string) (*model.User, error) {
	u, err := s.st.UserByID(ctx, userID)
	if err != nil {
		return nil, ferr.ErrResourceNotFound("user not found")
	}
	return u, nil
}

func (s *Service) issue(ctx context.Context, u *model.User) (*dto.TokenResponse, error) {
	claims := &AccessClaims{UserID: strconv.Itoa(int(u.ID)), Email: u.Email, Role: u.Role}
	pair, err := s.tokens.GenerateTokenPair(ctx, claims, claims)
	if err != nil {
		return nil, err
	}
	return s.pairToResponse(pair), nil
}

func (s *Service) pairToResponse(pair *fjwt.TokenPair) *dto.TokenResponse {
	now := time.Now()
	return &dto.TokenResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.accessTTL.Seconds()),
		ExpiresAt:    now.Add(s.accessTTL).Unix(),
	}
}