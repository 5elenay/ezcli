package ezcli

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

// Create a command handler with built-in help command.
func NewApp(appName string) *CommandHandler {
	var handler CommandHandler
	handler.SetName(appName)

	handler.SetNotFoundFunction(func() {
		log.Fatal("Command not found! Please run command 'help' for list all commands.")
	})

	handler.AddCommand(&Command{
		Name:        "help",
		Description: "Built-in help command for application.",
		Usages:      []string{"help", "help <command name>"},
		Execute: func(c *Command) {
			data := c.CommandData

			if len(data.Arguments) == 0 {
				fmt.Printf("List of all commands. For more information: %s help <command>\n", handler.Name)
				for _, command := range handler.Commands {
					fmt.Printf("  %s | %s\n", command.Name, command.Description)
				}
			} else {
				commandName := data.Arguments[0]

				err := handler.FindCommand(commandName, func(c *Command) error {
					fmt.Printf("Command %s:\n  Description: %s\n  Usages:\n    %s\n", c.Name, c.Description, strings.Join(c.Usages, "\n    "))

					// Add options
					if len(c.Options) > 0 {
						fmt.Println("  Options:")

						for _, item := range c.Options {
							fmt.Printf("    %s | %s\n", item.Name, item.Description)
						}
					}

					// Sub-commands
					if len(c.SubCommands) > 0 {
						fmt.Println("  Sub-Commands:")

						for _, item := range c.SubCommands {
							fmt.Printf("    Name: %s\n    Description: %s\n    Usages:\n      %s\n\n", item.Name, item.Description, strings.Join(item.Usages, "\n      "))
						}
					}

					return nil
				})

				// Command not found
				if err != nil {
					log.Fatal(err)
				}
			}
		},
	})

	return &handler
}

// Find a command from handler.
func (ch *CommandHandler) FindCommand(name string, fn func(c *Command) error) error {
	for _, item := range ch.Commands {
		if strings.EqualFold(item.Name, name) || stringSliceContains(item.Aliases, name) {
			return fn(item)
		}
	}

	return fmt.Errorf("Command not found! Please check your parameter")
}

// Find an option template from command.
func (c *Command) FindOption(name string, fn func(o *CommandOption)) {
	for _, item := range c.Options {
		if strings.EqualFold(item.Name, name) || stringSliceContains(item.Aliases, name) {
			fn(item)
		}
	}
}

// Find an option from command data.
func (c *CommandData) FindOption(name string, fn func(o *CommandOption)) {
	for _, item := range c.Options {
		if strings.EqualFold(item.Name, name) || stringSliceContains(item.Aliases, name) {
			fn(item)
		}
	}
}

// Ask the question.
func (q *Question) Ask() error {
	// Ask the question
	fmt.Print(q.Input)

	// Wait for input
	reader := bufio.NewReader(os.Stdin)
	result, err := reader.ReadString('\n')

	if err != nil {
		return err
	}

	// Change struct
	q.Answer = strings.TrimSpace(result)
	return nil
}

// Find Sub-command
func (c Command) FindSubcommand(name string, fn func(sc *SubCommand)) error {
	for _, item := range c.SubCommands {
		if strings.EqualFold(item.Name, name) {
			fn(item)
			return nil
		}
	}

	return fmt.Errorf("sub-command not found")
}
