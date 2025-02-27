package model

type User struct {
	BaseModel
	Username string `gorm:"not null;column:user_name" json:"username"`
	Password string `gorm:"not null;column:password" json:"password"`
}
