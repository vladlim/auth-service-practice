package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/vladlim/auth-service-practice/auth/internal/providers/auth"
	"github.com/vladlim/auth-service-practice/auth/internal/providers/tokens"
)

// Auth...

func (s *Server) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterUserData

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		s.respondWithError(w, http.StatusBadRequest, "Username, email, password are required")
		return
	}

	userID, err := s.authProvider.RegisterUser(r.Context(), ServerRegisterReq2Provider(req))
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrEmailExists):
			s.respondWithError(w, http.StatusConflict, "email exists")
		case errors.Is(err, auth.ErrUsernameExists):
			s.respondWithError(w, http.StatusConflict, "username exists")
		default:
			s.respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	accessToken, err := s.tokensProvider.GenerateAccessToken(userID)
	if err != nil {
		switch {
		case errors.Is(err, tokens.ErrAccessGenerate):
			s.respondWithError(w, http.StatusConflict, "access generate error: "+err.Error())
		default:
			s.respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	refreshToken, err := s.tokensProvider.GenerateRefreshToken(userID)
	if err != nil {
		switch {
		case errors.Is(err, tokens.ErrRefreshGenerate):
			s.respondWithError(w, http.StatusConflict, "refresh generate error: "+err.Error())
		default:
			s.respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	s.respondWithJSON(w, http.StatusCreated, Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (s *Server) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginUserData

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		s.respondWithError(w, http.StatusBadRequest, "login and password are required")
		return
	}

	userID, err := s.authProvider.LoginUser(r.Context(), req.Login, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrIncorrectPassword):
			s.respondWithError(w, http.StatusUnauthorized, "incorrect password")
		case errors.Is(err, auth.ErrUserNotFound):
			s.respondWithError(w, http.StatusNotFound, "user not found")
		default:
			s.respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	accessToken, err := s.tokensProvider.GenerateAccessToken(userID)
	if err != nil {
		switch {
		case errors.Is(err, tokens.ErrAccessGenerate):
			s.respondWithError(w, http.StatusConflict, "access generate error: "+err.Error())
		default:
			s.respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	refreshToken, err := s.tokensProvider.GenerateRefreshToken(userID)
	if err != nil {
		switch {
		case errors.Is(err, tokens.ErrRefreshGenerate):
			s.respondWithError(w, http.StatusConflict, "refresh generate error: "+err.Error())
		default:
			s.respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	s.respondWithJSON(w, http.StatusCreated, Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (s *Server) refreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.RefreshToken == "" {
		s.respondWithError(w, http.StatusBadRequest, "refresh token is required")
		return
	}

	claims, err := s.tokensProvider.ValidateRefreshToken(req.RefreshToken)

	if err != nil {
		switch {
		case errors.Is(err, tokens.ErrInvalidToken):
			s.respondWithError(w, http.StatusUnauthorized, "invalid token")
		default:
			s.respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	accessToken, err := s.tokensProvider.GenerateAccessToken(claims.UserID)
	if err != nil {
		switch {
		case errors.Is(err, tokens.ErrAccessGenerate):
			s.respondWithError(w, http.StatusConflict, "access generate error: "+err.Error())
		default:
			s.respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	refreshToken, err := s.tokensProvider.GenerateRefreshToken(claims.UserID)
	if err != nil {
		switch {
		case errors.Is(err, tokens.ErrRefreshGenerate):
			s.respondWithError(w, http.StatusConflict, "refresh generate error: "+err.Error())
		default:
			s.respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	s.respondWithJSON(w, http.StatusCreated, Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

// Keys

func (s *Server) generateKeyHandler(w http.ResponseWriter, r *http.Request) {
	// TO DO Middleware or check for admin

	var req struct {
		Role           string `json:"role"`
		GroupID        string `json:"group_id,omitempty"`
		UniversityID   string `json:"university_id"`
		EnrollmentYear int    `json:"enrollment_year,omitempty"`
		Degree         string `json:"degree,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "invalid request")
		return
	}

	switch req.Role {
	case "student":
		if req.GroupID == "" || req.UniversityID == "" || req.EnrollmentYear == 0 {
			s.respondWithError(w, http.StatusBadRequest, "group_id, university_id and enrollment_year are required for student")
			return
		}
	case "teacher":
		if req.UniversityID == "" || req.Degree == "" {
			s.respondWithError(w, http.StatusBadRequest, "university_id and degree are required for teacher")
			return
		}
	default:
		s.respondWithError(w, http.StatusBadRequest, "invalid role")
		return
	}

	key, err := s.tokensProvider.GenerateRoleKey(
		req.Role,
		req.GroupID,
		req.UniversityID,
		req.EnrollmentYear,
		req.Degree,
	)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, "failed to generate key")
		return
	}

	if claims, err := s.tokensProvider.ValidateRoleKey(key); err != nil {
		log.Default().Println("[ERR]: ", err.Error())
	} else {
		log.Default().Println("[GENERATED KEY INFO]: \n", claims)
	}

	s.respondWithJSON(w, http.StatusOK, map[string]string{
		"key": key,
	})
}

func (s *Server) activateKeyHandler(w http.ResponseWriter, r *http.Request) {
	claims, err := s.getClaimsFromRequest(r)
	if err != nil {
		s.respondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}
	userID := claims.UserID

	var req struct {
		Key string `json:"key"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondWithError(w, http.StatusBadRequest, "invalid request format")
		return
	}

	keyClaims, err := s.tokensProvider.ValidateRoleKey(req.Key)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var activateErr error
	switch keyClaims["role"] {
	case "student":
		activateErr = s.authProvider.ActivateStudent(r.Context(), userID, keyClaims)
	case "teacher":
		activateErr = s.authProvider.ActivateTeacher(r.Context(), userID, keyClaims)
	default:
		s.respondWithError(w, http.StatusBadRequest, "invalid role in activation key")
		log.Default().Println(keyClaims["role"])
		return
	}

	if activateErr != nil {
		switch {
		case errors.Is(activateErr, auth.ErrUserNotFound):
			s.respondWithError(w, http.StatusNotFound, "user not found")
		default:
			s.respondWithError(w, http.StatusInternalServerError, activateErr.Error())
		}
		return
	}

	s.respondWithJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// User Info...

func (s *Server) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	if userID == "" {
		s.respondWithError(w, http.StatusBadRequest, "user ID is required")
		return
	}

	user, err := s.authProvider.GetUserByID(context.Background(), userID)

	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserNotFound):
			s.respondWithError(w, http.StatusNotFound, "user not found")
		default:
			s.respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := ProviderUser2Server(user)
	s.respondWithJSON(w, http.StatusOK, resp)
}

func (s *Server) getUserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("email")
	if email == "" {
		s.respondWithError(w, http.StatusBadRequest, "email is required")
		return
	}

	log.Default().Println(email)

	user, err := s.authProvider.GetUserByEmail(context.Background(), email)

	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserNotFound):
			s.respondWithError(w, http.StatusNotFound, "user not found")
		default:
			s.respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	resp := ProviderUser2Server(user)
	s.respondWithJSON(w, http.StatusOK, resp)
}

func (s *Server) getStudentByIdHandler(w http.ResponseWriter, r *http.Request) {
	studentID := r.PathValue("id")
	if studentID == "" {
		s.respondWithError(w, http.StatusBadRequest, "student ID is required")
		return
	}

	student, err := s.authProvider.GetStudentByID(r.Context(), studentID)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserNotFound):
			s.respondWithError(w, http.StatusNotFound, "user not found")
		default:
			s.respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	s.respondWithJSON(w, http.StatusOK, student)
}

func (s *Server) getStudentsByGroupHandler(w http.ResponseWriter, r *http.Request) {
	groupIDs := r.PathValue("id")

	log.Default().Println("[GROUP ID] ", groupIDs)

	students, err := s.authProvider.GetStudents(r.Context(), groupIDs)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.respondWithJSON(w, http.StatusOK, students)
}

func (s *Server) getTeacherByIdHandler(w http.ResponseWriter, r *http.Request) {
	teacherID := r.PathValue("id")
	if teacherID == "" {
		s.respondWithError(w, http.StatusBadRequest, "teacher ID is required")
		return
	}

	teacher, err := s.authProvider.GetTeacherByID(r.Context(), teacherID)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrUserNotFound):
			s.respondWithError(w, http.StatusNotFound, "user not found")
		default:
			s.respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	s.respondWithJSON(w, http.StatusOK, teacher)
}

func (s *Server) getTeachersByUniHandler(w http.ResponseWriter, r *http.Request) {
	uniIDs := r.PathValue("id")

	log.Default().Println("[UNI ID] ", uniIDs)

	teachers, err := s.authProvider.GetTeachers(r.Context(), uniIDs)
	if err != nil {
		s.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	s.respondWithJSON(w, http.StatusOK, teachers)
}

func (s *Server) getUserRoleByIdHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")

	roles, err := s.authProvider.GetUserRoles(r.Context(), userID)
	if err != nil {
		log.Printf("GetUserRoles error: %v", err)
		s.respondWithError(w, http.StatusInternalServerError, "failed to get user roles")
		return
	}

	s.respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"user_id": userID,
		"roles":   roles,
	})
}

// Server ...

func (s *Server) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func (s *Server) getClaimsFromRequest(r *http.Request) (*tokens.Claims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("authorization header is required")
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return nil, errors.New("authorization header format must be 'Bearer {token}'")
	}

	tokenString := strings.TrimPrefix(authHeader, bearerPrefix)

	claims, err := s.tokensProvider.ValidateRefreshToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	return claims, nil
}

// Responces:
func setCommonHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func (s *Server) respondWithError(w http.ResponseWriter, code int, message string) {
	s.respondWithJSON(w, code, map[string]string{"error": message})
}

func (s *Server) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	setCommonHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to encode response"))
	}
}
