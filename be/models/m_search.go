package models

type SearchBodyRequest struct {
	Start string `json:"start" validate:"required"`
	Goal string `json:"goal" validate:"required"`
}