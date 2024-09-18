// File: internal/models/models.go

package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Softwares struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Version     string             `bson:"version" json:"version"`
	Description string             `bson:"description" json:"description"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// Add other model structs here, if any

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type User struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Username  string    `json:"username" bson:"username"`
	Email     string    `json:"email" bson:"email"`
	Password  string    `json:"-" bson:"password"` // Password is not returned in JSON
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	Role      string    `json:"role" bson:"role"`
}

type Platform struct {
	ID            string    `json:"id" bson:"_id,omitempty"`
	Name          string    `json:"name" bson:"name"`
	Description   string    `json:"description" bson:"description"`
	Version       string    `json:"version" bson:"version"`
	URL           string    `json:"url" bson:"url"`
	Documentation string    `json:"documentation" bson:"documentation"`
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" bson:"updatedAt"`
}

type APIService struct {
	ID          string     `json:"id" bson:"_id,omitempty"`
	Name        string     `json:"name" bson:"name"`
	Description string     `json:"description" bson:"description"`
	Version     string     `json:"version" bson:"version"`
	Endpoints   []Endpoint `json:"endpoints" bson:"endpoints"`
	SoftwareID  string     `json:"softwareId" bson:"softwareId"`
	PlatformID  string     `json:"platformId" bson:"platformId"`
	CreatedAt   time.Time  `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt" bson:"updatedAt"`
}

type Parameter struct {
	Name        string `json:"name" bson:"name"`
	In          string `json:"in" bson:"in"` // e.g., "query", "path", "body"
	Description string `json:"description" bson:"description"`
	Required    bool   `json:"required" bson:"required"`
	Type        string `json:"type" bson:"type"`
}

type Response struct {
	StatusCode  int    `json:"statusCode" bson:"statusCode"`
	Description string `json:"description" bson:"description"`
	Schema      string `json:"schema" bson:"schema"`
}

type APIRegistration struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	APIServiceID string    `json:"apiServiceId" bson:"apiServiceId"`
	UserID       string    `json:"userId" bson:"userId"`
	Status       string    `json:"status" bson:"status"` // e.g., "pending", "approved", "rejected"
	CreatedAt    time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" bson:"updatedAt"`
}

type Usage struct {
	ID           string    `json:"id" bson:"_id,omitempty"`
	APIServiceID string    `json:"apiServiceId" bson:"apiServiceId"`
	UserID       string    `json:"userId" bson:"userId"`
	Endpoint     string    `json:"endpoint" bson:"endpoint"`
	Method       string    `json:"method" bson:"method"`
	Timestamp    time.Time `json:"timestamp" bson:"timestamp"`
	ResponseTime int64     `json:"responseTime" bson:"responseTime"` // in milliseconds
	StatusCode   int       `json:"statusCode" bson:"statusCode"`
}

type Metrics struct {
	ID                  string    `json:"id" bson:"_id,omitempty"`
	APIServiceID        string    `json:"apiServiceId" bson:"apiServiceId"`
	Date                time.Time `json:"date" bson:"date"`
	TotalCalls          int       `json:"totalCalls" bson:"totalCalls"`
	AverageResponseTime int64     `json:"averageResponseTime" bson:"averageResponseTime"`
	ErrorRate           float64   `json:"errorRate" bson:"errorRate"`
}
