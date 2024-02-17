package presentation

import "github.com/khulnasoft/lazydocker/pkg/commands"

func GetProjectDisplayStrings(project *commands.Project) []string {
	return []string{project.Name}
}
