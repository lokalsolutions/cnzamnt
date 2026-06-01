package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"cnzamnt/backend/internal/service"
)

type Server struct {
	db      *sql.DB
	service *service.Service
}

func New(database *sql.DB) *Server {
	return &Server{db: database, service: service.New(database)}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.health)
	mux.HandleFunc("POST /api/dev/users", s.createUser)
	mux.HandleFunc("GET /api/users/me", s.me)
	mux.HandleFunc("POST /api/artworks", s.createArtwork)
	mux.HandleFunc("GET /api/artworks", s.artworks)
	mux.HandleFunc("GET /api/artworks/{id}", s.artwork)
	mux.HandleFunc("POST /api/artworks/{id}/comments", s.createComment)
	mux.HandleFunc("GET /api/artworks/{id}/comments", s.comments)
	return mux
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
		"app":    "CnzAMnt",
	})
}

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var input service.CreateUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	user, err := s.service.CreateUser(input)
	s.respondCreated(w, user, err)
}

func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	userID, err := devUserID(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	user, err := s.service.User(userID)
	s.respond(w, user, err)
}

func (s *Server) createArtwork(w http.ResponseWriter, r *http.Request) {
	userID, err := devUserID(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	var input service.CreateArtworkInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	artwork, err := s.service.CreateArtwork(userID, input)
	s.respondCreated(w, artwork, err)
}

func (s *Server) artworks(w http.ResponseWriter, _ *http.Request) {
	artworks, err := s.service.Artworks()
	s.respond(w, artworks, err)
}

func (s *Server) artwork(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid artwork id")
		return
	}
	artwork, err := s.service.Artwork(id)
	s.respond(w, artwork, err)
}

func (s *Server) createComment(w http.ResponseWriter, r *http.Request) {
	userID, err := devUserID(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	artworkID, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid artwork id")
		return
	}
	var input service.CreateCommentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	comment, err := s.service.CreateComment(r.Context(), userID, artworkID, input)
	s.respondCreated(w, comment, err)
}

func (s *Server) comments(w http.ResponseWriter, r *http.Request) {
	artworkID, err := pathID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid artwork id")
		return
	}
	comments, err := s.service.Comments(artworkID)
	s.respond(w, comments, err)
}

func (s *Server) respond(w http.ResponseWriter, value any, err error) {
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, value)
}

func (s *Server) respondCreated(w http.ResponseWriter, value any, err error) {
	if err != nil {
		writeServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, value)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		writeError(w, http.StatusNotFound, "not found")
	case errors.Is(err, service.ErrInsufficientCNZ):
		writeError(w, http.StatusPaymentRequired, "not enough CNZ")
	case errors.Is(err, service.ErrForbidden):
		writeError(w, http.StatusForbidden, "forbidden")
	default:
		writeError(w, http.StatusBadRequest, err.Error())
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func devUserID(r *http.Request) (int64, error) {
	value := r.Header.Get("X-Dev-User-Id")
	if value == "" {
		return 0, errors.New("X-Dev-User-Id header is required")
	}
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("X-Dev-User-Id header must be a positive integer")
	}
	return id, nil
}

func pathID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}
