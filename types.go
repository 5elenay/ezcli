package ezcli

// Command Handler struct.
type CommandHandler struct {
	Commands []*Command
}

// Command option (flag) struct.
type CommandOption struct {
	Name, Description string
}

// Command Data struct (arguments, options...)
type CommandData struct {
	Arguments []string
	Options   []*CommandOption
}

// Command struct.
type Command struct {
	Name, Description string
	Options           []*CommandOption
	Usages            []string
	CommandData       *CommandData
	Execute           func(data *CommandData)
}
