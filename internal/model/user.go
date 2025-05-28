package model

type User struct {
	BaseModel
	Username        string `gorm:"not null;column:username" json:"username"`
	Password        string `gorm:"not null;column:password" json:"password"`
	ExtractionCount int    `gorm:"not null;column:extraction_count" json:"extraction_count"`
	RemainingCount  int    `gorm:"not null;column:remaining_count" json:"remaining_count"`
	UserRole        int    `gorm:"not null;column:user_role" json:"user_role"`
}

func (User) TableName() string {
	return "users"
}
