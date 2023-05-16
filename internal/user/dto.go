package user

type CreateUserDTO struct {
	Email          string `json:"email"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password" bson:"-"`
}

type UpdateUserDTO struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

type LoginDTO struct {
	username string `json:"username"`
	password string `json:"password"`
}
