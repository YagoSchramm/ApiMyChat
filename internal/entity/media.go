package entity

type MediaType string

const (
	MediaAudio MediaType = "MP3"
	MediaVideo           = "MP4"
	MediaImage MediaType = "JPG"
)

type Media struct {
	ID       string
	URL      string
	Type     MediaType
	Metadata *Metadata
}
