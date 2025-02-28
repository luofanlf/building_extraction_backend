package model

type User struct {
	BaseModel
	Username string `gorm:"not null;column:username" json:"username"`
	Password string `gorm:"not null;column:password" json:"password"`
}
