package buildtool

// PrintableCommand ...
type PrintableCommand interface {
	PrintableCommand() string
}

// RunnableCommand ...
type RunnableCommand interface {
	PrintableCommand() string
	SetCustomArgs(args []string)
	Run() error
}

// ModifyAbleCommand ...
type ModifyAbleCommand interface {
	SetCustomArgs(args []string)
}

// BuildCommandContains ...
func BuildCommandContains(commands []PrintableCommand, command PrintableCommand) bool {
	for _, cmd := range commands {
		if cmd.PrintableCommand() == command.PrintableCommand() {
			return true
		}
	}
	return false
}
