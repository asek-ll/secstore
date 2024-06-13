package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

type SetSecretCommand struct {
	Config
}

func (cmd *SetSecretCommand) Execute(args []string) error {
	if len(args) == 1 {
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer term.Restore(int(os.Stdin.Fd()), oldState)

		rw := struct {
			io.Reader
			io.Writer
		}{
			os.Stdin,
			os.Stdout,
		}

		t := term.NewTerminal(rw, "")
		value, err := t.ReadPassword("Input secret: ")
		if err != nil {
			return err
		}

		return cmd.Config.Store.Save(args[0], strings.TrimSpace(value))
	}

	if len(args) == 2 {
		return cmd.Config.Store.Save(args[0], args[1])
	}

	return fmt.Errorf("Required 1-2 args")
}
