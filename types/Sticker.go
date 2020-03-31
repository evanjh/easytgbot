package types

type Sticker struct {
	FileId       string
	FileUniqueId string
	Width        int
	Height       int
	IsAnimated   bool
	Thumb        *PhotoSize
	Emoji        string
	SetName      string
	MaskPosition *MaskPosition
	FileSize     int
}

type StickerSet struct {
	Name          string
	Title         string
	IsAnimated    bool
	ContainsMasks bool
	Stickers      []Sticker
}

type MaskPosition struct {
	Point   string
	x_shift float64
	y_shift float64
	scale   float64
}
