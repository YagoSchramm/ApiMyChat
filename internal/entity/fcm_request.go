package entity

type FcmRequest struct {
	Uid   string `json:"uid" binding:"required"`
	Token string `json:"token" binding:"required"`
}
