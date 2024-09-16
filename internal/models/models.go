// File: internal/models/models.go

package models

import (
	"time"
)

type SoftwareInfo struct {
	ID                     string                `json:"id" bson:"_id,omitempty"`
	Type                   string                `json:"@type" bson:"type"`
	Context                string                `json:"@context" bson:"context"`
	Name                   string                `json:"name" bson:"name"`
	Description            string                `json:"description" bson:"description"`
	Version                string                `json:"version" bson:"version"`
	Identifier             string                `json:"identifier" bson:"identifier"`
	CodeRepository         string                `json:"codeRepository" bson:"codeRepository"`
	ProgrammingLanguage    ProgrammingLanguage   `json:"programmingLanguage" bson:"programmingLanguage"`
	RuntimePlatform        string                `json:"runtimePlatform" bson:"runtimePlatform"`
	TargetProduct          string                `json:"targetProduct" bson:"targetProduct"`
	ApplicationCategory    string                `json:"applicationCategory" bson:"applicationCategory"`
	ApplicationSubCategory string                `json:"applicationSubCategory" bson:"applicationSubCategory"`
	DownloadURL            string                `json:"downloadUrl" bson:"downloadUrl"`
	InstallURL             string                `json:"installUrl" bson:"installUrl"`
	FileSize               string                `json:"fileSize" bson:"fileSize"`
	MemoryRequirements     string                `json:"memoryRequirements" bson:"memoryRequirements"`
	StorageRequirements    string                `json:"storageRequirements" bson:"storageRequirements"`
	Permissions            string                `json:"permissions" bson:"permissions"`
	ProcessorRequirements  string                `json:"processorRequirements" bson:"processorRequirements"`
	ReleaseNotes           string                `json:"releaseNotes" bson:"releaseNotes"`
	SoftwareHelp           string                `json:"softwareHelp" bson:"softwareHelp"`
	SoftwareRequirements   []SoftwareApplication `json:"softwareRequirements" bson:"softwareRequirements"`
	SoftwareSuggestions    []SoftwareApplication `json:"softwareSuggestions" bson:"softwareSuggestions"`
	OperatingSystems       []string              `json:"operatingSystem" bson:"operatingSystem"`
	SupportingData         string                `json:"supportingData" bson:"supportingData"`
	Author                 []Person              `json:"author" bson:"author"`
	Citation               []Citation            `json:"citation" bson:"citation"`
	Contributor            []Person              `json:"contributor" bson:"contributor"`
	CopyrightHolder        Person                `json:"copyrightHolder" bson:"copyrightHolder"`
	CopyrightYear          int                   `json:"copyrightYear" bson:"copyrightYear"`
	Creator                Person                `json:"creator" bson:"creator"`
	DateCreated            time.Time             `json:"dateCreated" bson:"dateCreated"`
	DateModified           time.Time             `json:"dateModified" bson:"dateModified"`
	DatePublished          time.Time             `json:"datePublished" bson:"datePublished"`
	Editor                 Person                `json:"editor" bson:"editor"`
	Encoding               string                `json:"encoding" bson:"encoding"`
	FileFormat             string                `json:"fileFormat" bson:"fileFormat"`
	Funder                 []Organization        `json:"funder" bson:"funder"`
	Keywords               []string              `json:"keywords" bson:"keywords"`
	License                string                `json:"license" bson:"license"`
	Producer               Organization          `json:"producer" bson:"producer"`
	Provider               Organization          `json:"provider" bson:"provider"`
	Publisher              Organization          `json:"publisher" bson:"publisher"`
	Sponsor                Organization          `json:"sponsor" bson:"sponsor"`
	IsAccessibleForFree    bool                  `json:"isAccessibleForFree" bson:"isAccessibleForFree"`
	IsPartOf               string                `json:"isPartOf" bson:"isPartOf"`
	HasPart                []string              `json:"hasPart" bson:"hasPart"`
	Position               string                `json:"position" bson:"position"`
	URL                    string                `json:"url" bson:"url"`
	SameAs                 string                `json:"sameAs" bson:"sameAs"`
	RelatedLink            []string              `json:"relatedLink" bson:"relatedLink"`
	Review                 []Review              `json:"review" bson:"review"`
	Maintainer             Person                `json:"maintainer" bson:"maintainer"`
	ContinuousIntegration  string                `json:"contIntegration" bson:"contIntegration"`
	BuildInstructions      string                `json:"buildInstructions" bson:"buildInstructions"`
	DevelopmentStatus      string                `json:"developmentStatus" bson:"developmentStatus"`
	EmbargoEndDate         time.Time             `json:"embargoEndDate" bson:"embargoEndDate"`
	Funding                []string              `json:"funding" bson:"funding"`
	IssueTracker           string                `json:"issueTracker" bson:"issueTracker"`
	ReferencePublication   []ScholarlyArticle    `json:"referencePublication" bson:"referencePublication"`
	Readme                 string                `json:"readme" bson:"readme"`
	HasSourceCode          string                `json:"hasSourceCode" bson:"hasSourceCode"`
	IsSourceCodeOf         string                `json:"isSourceCodeOf" bson:"isSourceCodeOf"`
	APIInfo                APIInfo               `json:"apiInfo" bson:"apiInfo"`
}

type Person struct {
	ID          string `json:"@id" bson:"id"`
	Type        string `json:"@type" bson:"type"`
	Name        string `json:"name" bson:"name"`
	GivenName   string `json:"givenName" bson:"givenName"`
	FamilyName  string `json:"familyName" bson:"familyName"`
	Email       string `json:"email" bson:"email"`
	Affiliation string `json:"affiliation" bson:"affiliation"`
}

type Organization struct {
	ID    string `json:"@id" bson:"id"`
	Type  string `json:"@type" bson:"type"`
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	URL   string `json:"url" bson:"url"`
}

type ProgrammingLanguage struct {
	Type    string `json:"@type" bson:"type"`
	Name    string `json:"name" bson:"name"`
	Version string `json:"version" bson:"version"`
	URL     string `json:"url" bson:"url"`
}

type SoftwareApplication struct {
	Type     string       `json:"@type" bson:"type"`
	Name     string       `json:"name" bson:"name"`
	Version  string       `json:"version" bson:"version"`
	Provider Organization `json:"provider" bson:"provider"`
}

type Citation struct {
	Type string `json:"@type" bson:"type"`
	Text string `json:"text" bson:"text"`
	URL  string `json:"url" bson:"url"`
}

type Review struct {
	Type         string `json:"@type" bson:"type"`
	ReviewAspect string `json:"reviewAspect" bson:"reviewAspect"`
	ReviewBody   string `json:"reviewBody" bson:"reviewBody"`
}

type ScholarlyArticle struct {
	Type string `json:"@type" bson:"type"`
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
	DOI  string `json:"doi" bson:"doi"`
}

type APIInfo struct {
	Endpoints []Endpoint `json:"endpoints" bson:"endpoints"`
}

type Endpoint struct {
	Path        string `json:"path" bson:"path"`
	Method      string `json:"method" bson:"method"`
	Description string `json:"description" bson:"description"`
}

// SoftwareInfo and related structs remain the same as in the previous message

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

type Endpoint struct {
	Path        string      `json:"path" bson:"path"`
	Method      string      `json:"method" bson:"method"`
	Description string      `json:"description" bson:"description"`
	Parameters  []Parameter `json:"parameters" bson:"parameters"`
	Responses   []Response  `json:"responses" bson:"responses"`
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
