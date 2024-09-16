// File: internal/models/software_info.go

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
