package handler

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func utilProblem(status int, problemType string, detail string) error {
	return &huma.ErrorModel{
		Type:   "urn:problem:" + problemType,
		Title:  http.StatusText(status),
		Status: status,
		Detail: detail,
	}
}
