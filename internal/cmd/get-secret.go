package cmd

import "fmt"

type SetSecretCommand struct {
	Config
}

func (cmd *SetSecretCommand) Execute(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("Required exact 2 args")
	}

	return cmd.Config.Store.Save(args[0], args[1])
}
