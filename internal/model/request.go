package model

type RequestStatus int

const (
	RequestStatusPending  RequestStatus = 0 // 待处理
	RequestStatusApproved RequestStatus = 1 // 已批准
	RequestStatusRejected RequestStatus = 2 // 已拒绝
)

type ExtractionRequest struct {
	BaseModel
	UserID       int           `gorm:"not null;column:user_id" json:"user_id"`
	User         User          `gorm:"foreignKey:UserID" json:"user"` // 添加与User的关联
	RequestCount int           `gorm:"not null;column:request_count" json:"request_count"`
	Status       RequestStatus `gorm:"not null;column:status" json:"status"`
	Reason       string        `gorm:"column:reason" json:"reason"`
	Reply        string        `gorm:"column:reply" json:"reply"`
}

func (ExtractionRequest) TableName() string {
	return "extraction_requests"
}
