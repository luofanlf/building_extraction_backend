package dto

type CreateRequestRequest struct {
	RequestCount int    `json:"request_count" binding:"required"`
	Reason       string `json:"reason" binding:"required"`
}

type RequestResponse struct {
	ID           int    `json:"id"`
	UserID       int    `json:"user_id"`
	Username     string `json:"username"`
	RequestCount int    `json:"request_count"`
	Status       int    `json:"status"`
	Reason       string `json:"reason"`
	Reply        string `json:"reply"`
	CreatedAt    string `json:"created_at"`
}

type HandleRequestRequest struct {
	RequestID int    `json:"request_id" binding:"required"`
	Status    int    `json:"status" binding:"required"` // 1: 批准, 2: 拒绝
	Reply     string `json:"reply"`
}

type UserRequestResponse struct {
	ID           int    `json:"id"`
	RequestCount int    `json:"request_count"`
	Status       int    `json:"status"`
	Reason       string `json:"reason"`
	Reply        string `json:"reply"`
	CreatedAt    string `json:"created_at"`
}
