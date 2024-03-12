package cmd

import "fmt"

type GetSecretCommand struct {
	Config
}

func (cmd *GetSecretCommand) Execute(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("Required exact 1 arg")
	}

	res, err := cmd.Config.Store.Load(args[0])
	if err != nil {
		return err
	}

	fmt.Print(res)

	return nil
}
