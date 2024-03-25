package model

type UserSession struct{
	HashedSessionID string `json:"HashedSessionID"`
	HashedCsrfToken string `json:"HashedCsrfToken"`
	UserData UserGeneralInfo `json:"UserData"`
	UserMetaData UserMetaData `json:"UserMetaData"`
}

type Credentials struct{
	Email string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}

type UserSignUp struct {
	Email string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
	Username string `json:"username" validate:"required"`
}

type EmailConfirmation struct {
	Email string `json:"email" validate:"email,required"`
	Code string `json:"code" validate:"required"`
}

type PasswordChange struct {
	OldPassword *string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}