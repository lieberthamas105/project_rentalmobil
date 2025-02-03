package model

type UserCredential struct {
	Id       uint32 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}
