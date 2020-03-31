package types

import (
	"fmt"

	"github.com/mylukin/easytgbot/configs"
)

type File struct {
	FileId       string
	FileUniqueId string
	FileSize     int
	FilePath     string
}

// Link returns a full path to the download Url for a File.
//
// It requires the Bot Token to create the link.
func (f *File) Link(token string) string {
	return fmt.Sprintf(configs.FileEndpoint, token, f.FilePath)
}
