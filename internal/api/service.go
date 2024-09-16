// File: internal/api/service.go

package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/chenxingqiang/soft-crusher/internal/database"
	"github.com/chenxingqiang/soft-crusher/internal/models"
	"github.com/gorilla/mux"
)

type Service struct {
	DB database.Database
}

func NewService(db database.Database) *Service {
	return &Service{DB: db}
}

func (s *Service) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/software", s.CreateSoftwareInfo).Methods("POST")
	r.HandleFunc("/software/{id}", s.GetSoftwareInfo).Methods("GET")
	r.HandleFunc("/software/{id}", s.UpdateSoftwareInfo).Methods("PUT")
	r.HandleFunc("/software/{id}", s.DeleteSoftwareInfo).Methods("DELETE")
	r.HandleFunc("/software", s.ListSoftwareInfo).Methods("GET")
	r.HandleFunc("/software/search", s.SearchSoftwareInfo).Methods("GET")
	r.HandleFunc("/software/repository", s.GetSoftwareInfoByCodeRepository).Methods("GET")

	// User routes
	r.HandleFunc("/users", s.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", s.GetUser).Methods("GET")
	r.HandleFunc("/users/{id}", s.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", s.DeleteUser).Methods("DELETE")
	r.HandleFunc("/users", s.ListUsers).Methods("GET")

	// Platform routes
	r.HandleFunc("/platforms", s.CreatePlatform).Methods("POST")
	r.HandleFunc("/platforms/{id}", s.GetPlatform).Methods("GET")
	r.HandleFunc("/platforms/{id}", s.UpdatePlatform).Methods("PUT")
	r.HandleFunc("/platforms/{id}", s.DeletePlatform).Methods("DELETE")
	r.HandleFunc("/platforms", s.ListPlatforms).Methods("GET")

	// API Service routes
	r.HandleFunc("/api-services", s.CreateAPIService).Methods("POST")
	r.HandleFunc("/api-services/{id}", s.GetAPIService).Methods("GET")
	r.HandleFunc("/api-services/{id}", s.UpdateAPIService).Methods("PUT")
	r.HandleFunc("/api-services/{id}", s.DeleteAPIService).Methods("DELETE")
	r.HandleFunc("/api-services", s.ListAPIServices).Methods("GET")
}

// Software Info handlers
func (s *Service) CreateSoftwareInfo(w http.ResponseWriter, r *http.Request) {
	var info models.SoftwareInfo
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.DB.SaveSoftwareInfo(&info); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(info)
}

func (s *Service) GetSoftwareInfo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	info, err := s.DB.GetSoftwareInfo(id)
	if err != nil {
		http.Error(w, "Software not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(info)
}

func (s *Service) UpdateSoftwareInfo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var info models.SoftwareInfo
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	info.ID = id
	if err := s.DB.UpdateSoftwareInfo(&info); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(info)
}

func (s *Service) DeleteSoftwareInfo(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := s.DB.DeleteSoftwareInfo(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) ListSoftwareInfo(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	info, err := s.DB.ListSoftwareInfo(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(info)
}

func (s *Service) SearchSoftwareInfo(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	info, err := s.DB.SearchSoftwareInfo(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(info)
}

func (s *Service) GetSoftwareInfoByCodeRepository(w http.ResponseWriter, r *http.Request) {
	repoURL := r.URL.Query().Get("url")
	if repoURL == "" {
		http.Error(w, "Repository URL is required", http.StatusBadRequest)
		return
	}
	info, err := s.DB.GetSoftwareInfoByCodeRepository(repoURL)
	if err != nil {
		http.Error(w, "Software not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(info)
}

// User handlers
// Implement CreateUser, GetUser, UpdateUser, DeleteUser, ListUsers

// Platform handlers
// Implement CreatePlatform, GetPlatform, UpdatePlatform, DeletePlatform, ListPlatforms

// API Service handlers
// Implement CreateAPIService, GetAPIService, UpdateAPIService, DeleteAPIService, ListAPIServices

func (s *Service) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.DB.SaveUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (s *Service) GetUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	user, err := s.DB.GetUser(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (s *Service) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.ID = id
	if err := s.DB.UpdateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (s *Service) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := s.DB.DeleteUser(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Service) ListUsers(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	users, err := s.DB.ListUsers(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}
