package model

type UserSession struct{
	HashedSessionID string
	HashedCsrfToken string
	UserData UserGeneralInfo
	UserMetaData UserMetaData
}

type Credentials struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserSignUp struct {
	Email string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type Confirmation struct {
	Email string `json:"email"`
	Code string `json:"code"`
}

type PasswordChange struct {
	OldPassword string `json:"oldPassword"`
	NewPassword *string `json:"newPassword"`
	Session UserSession
}