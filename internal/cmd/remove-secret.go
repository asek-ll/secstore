package cmd

import "fmt"

type RemoveSecretCommand struct {
	Config
}

func (cmd *RemoveSecretCommand) Execute(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Required exact 1 arg")
	}

	return cmd.Config.Store.Delete(args[0])
}
