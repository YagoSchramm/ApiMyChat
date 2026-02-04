package entity

type MessageType string

const (
	TypeText  MessageType = "TEXT"
	TypeImage MessageType = "IMAGE"
	TypeAudio MessageType = "AUDIO"
	TypeVideo MessageType = "VIDEO"
)
