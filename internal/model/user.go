package model

import (
	"encoding/json"
	"time"
)

type User struct{
	ID int `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Email string `json:"email" db:"email"`
	PhoneNumber *string `json:"phoneNumber" db:"phone_number"`
	IsActive *string `json:"is_active" db:"is_active"`
	SteamID *string `json:"steamId" db:"steam_id"`
	PhotoPath *string `json:"photo" db:"photo_path"`
	BackgroundPath *string `json:"background" db:"background_path"`
}

type UserGeneralInfo struct {
	ID int `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	SteamID bool `json:"steamId" db:"steam_id"`
	PhotoPath *string `json:"photo" db:"photo_path"`
	BackgroundPath *string `json:"background" db:"background_path"`
}

type UserMetaData struct{
	ID int `json:"id" db:"id"`
	Timezone time.Location `json:"timezone" db:"timezone"`
	IPAdress string `json:"ip_adress" db:"ip_adress"`
	Country string `json:"country" db:"country"`
	UserId int `json:"user_id" db:"user_id"`
}

func (u *UserGeneralInfo) MarshalBinary() ([]byte, error) {
    return json.Marshal(u)
}

func (u *UserMetaData) MarshalBinary() ([]byte, error) {
    return json.Marshal(u)
}