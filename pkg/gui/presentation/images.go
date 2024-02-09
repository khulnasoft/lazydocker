package presentation

import (
	"github.com/khulnasoft/lazydocker/pkg/commands"
	"github.com/khulnasoft/lazydocker/pkg/utils"
)

func GetImageDisplayStrings(image *commands.Image) []string {
	return []string{
		image.Name,
		image.Tag,
		utils.FormatDecimalBytes(int(image.Image.Size)),
	}
}
