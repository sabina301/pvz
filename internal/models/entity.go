package models

import (
	"pvz/internal/models/auth"
	"time"
)

type User struct {
	Id           string
	Email        string
	PasswordHash string
	Role         auth.Role
}

type Pvz struct {
	Id               string
	RegistrationDate time.Time
	City             string
}

type Reception struct {
	Id       string
	DateTime time.Time
	PvzId    string
	Status   string
}

type Product struct {
	Id          string
	DateTime    time.Time
	Type        string
	ReceptionId string
}

type PvzReception struct {
	PvzId       string
	ReceptionId string
}
