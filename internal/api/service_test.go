package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chenxingqiang/soft-crusher/internal/database"
	"github.com/chenxingqiang/soft-crusher/internal/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDatabase is a mock implementation of the Database interface
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) SaveSoftwareInfo(info *models.SoftwareInfo) error {
	args := m.Called(info)
	return args.Error(0)
}

func (m *MockDatabase) GetSoftwareInfo(id string) (*models.SoftwareInfo, error) {
	args := m.Called(id)
	return args.Get(0).(*models.SoftwareInfo), args.Error(1)
}

// Implement other methods of the Database interface...

func TestCreateSoftwareInfo(t *testing.T) {
	mockDB := new(MockDatabase)
	service := NewService(mockDB)

	info := &models.SoftwareInfo{Name: "Test Software"}
	mockDB.On("SaveSoftwareInfo", info).Return(nil)

	body, _ := json.Marshal(info)
	req, _ := http.NewRequest("POST", "/software", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/software", service.CreateSoftwareInfo).Methods("POST")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockDB.AssertExpectations(t)
}

func TestGetSoftwareInfo(t *testing.T) {
	mockDB := new(MockDatabase)
	service := NewService(mockDB)

	info := &models.SoftwareInfo{ID: "123", Name: "Test Software"}
	mockDB.On("GetSoftwareInfo", "123").Return(info, nil)

	req, _ := http.NewRequest("GET", "/software/123", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/software/{id}", service.GetSoftwareInfo).Methods("GET")
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var response models.SoftwareInfo
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, info.Name, response.Name)
	mockDB.AssertExpectations(t)
}

// Add more tests for other API methods...
