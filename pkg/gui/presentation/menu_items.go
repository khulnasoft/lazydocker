package presentation

import "github.com/khulnasoft/lazydocker/pkg/gui/types"

func GetMenuItemDisplayStrings(menuItem *types.MenuItem) []string {
	return menuItem.LabelColumns
}
