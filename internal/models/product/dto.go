package product

import "time"

type AddInReceptionRequest struct {
	Type  string `json:"type" validate:"required,oneof=электроника одежда обувь"`
	PvzId string `json:"pvzId" validate:"required,uuid"`
}

type AddInReceptionResponse struct {
	Id          string    `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	Type        string    `json:"type"`
	ReceptionId string    `json:"receptionId"`
}
