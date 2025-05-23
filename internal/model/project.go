package model

type Project struct {
	BaseModel
	ProjectName string `gorm:"not null;column:name" json:"name"`
	InputImage  string `gorm:"not null;column:input_image" json:"input_image"`
	OutputImage string `gorm:"not null;column:output_image" json:"output_image"`
	ModelName   string `gorm:"not null;column:model" json:"model"`
	UserID      int    `gorm:"not null;column:user_id" json:"user_id"`
}

func (Project) TableName() string {
	return "projects"
}
