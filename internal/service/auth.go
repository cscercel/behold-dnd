package service

import (
	"context"
	"fmt"
	"time"

	"github.com/cscercel/behold-dnd/internal/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	queries          *db.Queries
	jwtSecret        []byte
	jwtExpiryHours   int
	registrationCode string
}

func NewAuthService(queries *db.Queries, jwtSecret string, jwtExpiryHours int, registrationCode string) *AuthService {
	return &AuthService{
		queries:          queries,
		jwtSecret:        []byte(jwtSecret),
		jwtExpiryHours:   jwtExpiryHours,
		registrationCode: registrationCode,
	}
}

// Claims is the payload embeded inside the JWT token
// User ID and role are enough to make every auth decision
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}

func (s *AuthService) Register(ctx context.Context, username, email, password, role, code string) (db.User, error) {
	// Validate registration code
	if code != s.registrationCode {
		return db.User{}, fmt.Errorf("invalid registration code")
	}

	// Validate role
	if role != "player" && role != "dm" {
		return db.User{}, fmt.Errorf("role must be 'player' or 'dm', got: '%s'", role)
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return db.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{
		Username:       username,
		Email:          email,
		HashedPassword: string(hash),
		Role:           role,
	})
	if err != nil {
		return db.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, db.User, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return "", db.User{}, fmt.Errorf("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return "", db.User{}, fmt.Errorf("invalid email or password")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return "", db.User{}, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func (s *AuthService) generateToken(user db.User) (string, error) {
	claims := &Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.jwtExpiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
