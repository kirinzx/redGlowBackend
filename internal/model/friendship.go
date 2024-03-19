package model

import (
	"time"
)

type StatusName string

const (
	Ivited StatusName = "Заявка отправлена"
	Friends StatusName = "Друзья"
)

type Friendship struct{
	ID int `json:"id"`
	InviterID int `json:"inviter_id"`
	AccepterID int `json:"accepter_id"`
	Status StatusName `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}