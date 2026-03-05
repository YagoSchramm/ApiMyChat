package entity

type UpdateMessageStatus struct {
	MessageID string `json:"id" binding:"required"`
	Status    string `json:"status" binding:"required,oneof=sent delivered read"`
}
