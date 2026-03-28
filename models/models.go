package models

type RegisterRequest struct {
	Email       string  `json:"email"`
	Nickname    string  `json:"nickname"`
	Password    string  `json:"password"`
	Role        string  `json:"role"`
	Photo       *string `json:"photo"`
	City        *string `json:"city"`
	Status      *string `json:"status"`
	AgreementPD *bool   `json:"agreement_pd"`
	AgreementEA *bool   `json:"agreement_ea"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserInfo struct {
	Email    string  `json:"email"`
	Nickname string  `json:"nickname"`
	Photo    *string `json:"photo"`
	City     *string `json:"city"`
	Status   string  `json:"status"`
}

type UpdateProfileRequest struct {
	Nickname *string `json:"nickname"`
	Photo    *string `json:"photo"`
	City     *string `json:"city"`
	Status   *string `json:"status"`
	Password *string `json:"password"`
}
type AuthContext struct {
	UserID int64
	Role   string
}

type VerifyEmailRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type ResendRequest struct {
	Email string `json:"email"`
}
