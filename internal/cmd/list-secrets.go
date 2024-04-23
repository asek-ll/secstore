package cmd

import "fmt"

type ListSecretsCommand struct {
	Config
}

func (cmd *ListSecretsCommand) Execute(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("Required exact 0 arg")
	}

	items, err := cmd.Config.Store.List()
	if err != nil {
		return err
	}

	for _, item := range items {
		fmt.Println(item)
	}

	return nil
}
