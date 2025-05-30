package dto

type AdminStatsResponse struct {
	TotalUsers       int64 `json:"total_users"`
	TotalExtractions int64 `json:"total_extractions"`
	ActiveUsersToday int64 `json:"active_users_today"`
	ExtractionsToday int64 `json:"extractions_today"`
}

type UserManagementResponse struct {
	ID              int    `json:"id"`
	Username        string `json:"username"`
	CreatedAt       string `json:"created_at"`
	ExtractionCount int    `json:"extraction_count"`
	RemainingCount  int    `json:"remaining_count"`
	UserRole        int    `json:"user_role"`
}
