// File: internal/database/database.go

package database

import (
	"time"

	"github.com/chenxingqiang/soft-crusher/internal/models"
)

type Database interface {
	// SoftwareInfo CRUD
	SaveSoftwareInfo(info *models.SoftwareInfo) error
	GetSoftwareInfo(id string) (*models.SoftwareInfo, error)
	UpdateSoftwareInfo(info *models.SoftwareInfo) error
	DeleteSoftwareInfo(id string) error
	ListSoftwareInfo(limit, offset int) ([]*models.SoftwareInfo, error)
	SearchSoftwareInfo(query string) ([]*models.SoftwareInfo, error)
	GetSoftwareInfoByCodeRepository(url string) (*models.SoftwareInfo, error)

	// User CRUD
	CreateUser(user *models.User) error
	GetUser(id string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id string) error
	ListUsers(limit, offset int) ([]*models.User, error)

	// Platform CRUD
	CreatePlatform(platform *models.Platform) error
	GetPlatform(id string) (*models.Platform, error)
	UpdatePlatform(platform *models.Platform) error
	DeletePlatform(id string) error
	ListPlatforms(limit, offset int) ([]*models.Platform, error)

	// APIService CRUD
	CreateAPIService(service *models.APIService) error
	GetAPIService(id string) (*models.APIService, error)
	UpdateAPIService(service *models.APIService) error
	DeleteAPIService(id string) error
	ListAPIServices(limit, offset int) ([]*models.APIService, error)

	// APIRegistration CRUD
	CreateAPIRegistration(registration *models.APIRegistration) error
	GetAPIRegistration(id string) (*models.APIRegistration, error)
	UpdateAPIRegistration(registration *models.APIRegistration) error
	DeleteAPIRegistration(id string) error
	ListAPIRegistrations(limit, offset int) ([]*models.APIRegistration, error)

	// Usage logging
	LogUsage(usage *models.Usage) error
	GetUsage(apiServiceID string, startTime, endTime time.Time) ([]*models.Usage, error)

	// Metrics
	GetMetrics(apiServiceID string, startDate, endDate time.Time) ([]*models.Metrics, error)
	UpdateMetrics(metrics *models.Metrics) error
}
