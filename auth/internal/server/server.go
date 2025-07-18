package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/vladlim/auth-service-practice/auth/internal/config"
	"github.com/vladlim/auth-service-practice/auth/internal/providers/auth"
)

type Provider interface {
	RegisterUser(ctx context.Context, user auth.RegisterUserData) (string, error)
	LoginUser(ctx context.Context, login, password string) (string, error)
}

type Server struct {
	server   http.Server
	provider Provider
}

func New(conf config.Config, provider Provider) *Server {
	s := new(Server)
	s.server.Addr = fmt.Sprintf(":%d", conf.Port)
	s.server.Handler = s.setRouter()
	s.provider = provider
	return s
}

// Метод	Роут					Назначение
// POST		/auth/register			Регистрация нового пользователя
// POST		/auth/login				Вход в систему
// POST		/auth/logout			Выход из системы (очистка сессии или куки)
// GET		/users/{id}				Получить пользователя по ID
// GET		/users/email/{email}	Получить пользователя по email
// GET		/students/{id}			Получить студента по ID
// GET		/students				Фильтрация студентов по группам (через query)
// GET		/teachers/{id}			Получить преподавателя по ID
// GET		/teachers				Фильтрация преподавателей по университетам
// GET		/roles/{id}				Получить роль пользователя по ID

func (s *Server) setRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", s.pingHandler)
	mux.HandleFunc("POST /auth/register", s.registerUserHandler)
	mux.HandleFunc("POST /auth/login", s.loginUserHandler)
	// mux.HandleFunc("GET /users/{id}", s.getUserByIdHandler)
	// mux.HandleFunc("GET /users/email/{email}", s.getUserByEmailHandler)
	// mux.HandleFunc("GET /students/{id}", s.getStudentByIdHandler)
	// mux.HandleFunc("GET /students", s.getStudentsFilterByGroupHandler)
	// mux.HandleFunc("GET /teachers/{id}", s.getTeacherByIdHandler)
	// mux.HandleFunc("GET /teachers", s.getTeacherFilterByUniHandler)
	// mux.HandleFunc("GET /roles/{id}", s.getUserRoleByIdHandler)
	return mux
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}
