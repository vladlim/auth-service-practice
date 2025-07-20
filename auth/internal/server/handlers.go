package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

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

// Server ...

func (s *Server) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
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
