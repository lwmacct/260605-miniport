package handler

import "time"

type HealthDTO struct {
	Status    string    `json:"status" example:"ok"`
	Timestamp time.Time `json:"timestamp" example:"2026-06-15T12:00:00Z"`
	Version   string    `json:"version" example:"0.1.0"`
}

type MetaDTO struct {
	Name     string `json:"name" example:"Miniport"`
	Version  string `json:"version" example:"0.1.0"`
	DocsPath string `json:"docsPath" example:"/api"`
}
