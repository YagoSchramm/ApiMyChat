package entity

type UpdateUserModel struct {
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Description string `json:"description" binding:"required"`
	UID         string `json:"uid" binding:"required"`
}
