package models

type User struct {
	ID           int64  `json:"-"`
	Login        string `json:"login"`    // имя пользователя
	Password     string `json:"password"` // пароль
	PasswordHash string `json:"-"`
}
