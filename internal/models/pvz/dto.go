package pvz

import (
	"database/sql"
	"time"
)

type CreateRequest struct {
	Id               *string    `json:"id" validate:"omitempty,uuid"`
	RegistrationDate *time.Time `json:"registrationDate" validate:"omitempty"`
	City             string     `json:"city" validate:"required,oneof=Москва Санкт-Петербург Казань"`
}

type CreateResponse struct {
	Id               string    `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             string    `json:"city"`
}

type DeleteLastProductRequest struct {
	PvzId string `json:"pvzId" validate:"required,uuid"`
}

type DeleteLastProductResponse struct {
	Id string `json:"id"`
}

type CloseLastProductResponse struct {
	Id       string    `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PvzId    string    `json:"pvzId"`
	Status   string    `json:"status"`
}

type ListRequest struct {
	StartDate *time.Time `json:"startDate" validate:"required"`
	EndDate   *time.Time `json:"endDate" validate:"required"`
	Page      int        `json:"page" validate:"min=1" default:"1"`
	Limit     int        `json:"limit" validate:"min=1,max=30" default:"10"`
}

type ListResponse struct {
	Pvz        Pvz                 `json:"pvz"`
	Receptions []ReceptionProducts `json:"receptions"`
}

type Pvz struct {
	Id               string    `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             string    `json:"city"`
}

type ReceptionProducts struct {
	Reception Reception `json:"reception"`
	Products  []Product `json:"products"`
}

type Reception struct {
	Id       string    `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PvzId    string    `json:"pvzId"`
	Status   string    `json:"status"`
}

type Product struct {
	Id          string    `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	Type        string    `json:"type"`
	ReceptionId string    `json:"receptionId"`
}

type RawList struct {
	PvzId           string
	PvzRegDate      time.Time
	PvzCity         string
	ReceptionId     string
	ReceptionDate   time.Time
	ReceptionStatus string
	ProductId       sql.NullString
	ProductDate     sql.NullTime
	ProductType     sql.NullString
}
