package reception

import "time"

type CreateRequest struct {
	PvzId string `json:"pvzId" validate:"required,uuid"`
}

type CreateResponse struct {
	Id       string    `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PvzId    string    `json:"pvzId"`
	Status   string    `json:"status"`
}
