package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vladlim/auth-service-practice/auth/internal/config"
	"github.com/vladlim/auth-service-practice/auth/internal/providers/auth"
	"github.com/vladlim/auth-service-practice/auth/internal/providers/tokens"
)

type AuthProvider interface {
	RegisterUser(ctx context.Context, user auth.RegisterUserData) (string error)
	LoginUser(ctx context.Context, login, password string) (string, error)
	GetUserByID(ctx context.Context, userID string) (User, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetStudentByID(ctx context.Context, userID string) (Student, error)
}

type TokensProvider interface {
	GenerateAccessToken(userID string) (Tokens, error)
	GenerateRefreshToken(userID string) (Tokens, error)
	ValidateRefreshToken(tokenString string) (*tokens.Claims, error)
	GenerateRoleKey(role string, groupID string, universityID string, enrollmentYear int, degree string) (string, error)
	ValidateRoleKey(key string) (jwt.MapClaims, error)
}

type Server struct {
	server         http.Server
	authProvider   auth.AuthProvider
	tokensProvider tokens.TokensProvider
}

func New(conf config.Config, authProvider auth.AuthProvider, tokensProvider tokens.TokensProvider) *Server {
	s := new(Server)
	s.server.Addr = fmt.Sprintf(":%d", conf.Port)
	s.server.Handler = s.setRouter()
	s.authProvider = authProvider
	s.tokensProvider = tokensProvider
	return s
}

func (s *Server) setRouter() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", s.pingHandler)
	mux.HandleFunc("POST /auth/register", s.registerUserHandler)
	mux.HandleFunc("POST /auth/login", s.loginUserHandler)
	mux.HandleFunc("POST /auth/refresh", s.refreshTokenHandler)

	mux.HandleFunc("POST /admin/generate-key", s.generateKeyHandler)
	mux.HandleFunc("POST /auth/activate-key", s.activateKeyHandler)
	mux.HandleFunc("GET /users/{id}", s.getUserByIdHandler)
	mux.HandleFunc("GET /users/email/{email}", s.getUserByEmailHandler)
	mux.HandleFunc("GET /students/{id}", s.getStudentByIdHandler)
	mux.HandleFunc("GET /groups/{id}/students", s.getStudentsByGroupHandler)
	mux.HandleFunc("GET /teachers/{id}", s.getTeacherByIdHandler)
	mux.HandleFunc("GET /universities/{id}/teachers", s.getTeachersByUniHandler)
	mux.HandleFunc("GET /roles/{id}", s.getUserRoleByIdHandler)
	return mux
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}
