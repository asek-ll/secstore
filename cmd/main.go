package main

import (
	"errors"
	"os"

	"github.com/asek-ll/secstore/internal/cmd"
	"github.com/asek-ll/secstore/pkg/store"

	"github.com/jessevdk/go-flags"
)

type Options struct {
	GetSecretCmd    cmd.GetSecretCommand    `command:"get"`
	SetSecretCmd    cmd.SetSecretCommand    `command:"set"`
	RemoveSecretCmd cmd.RemoveSecretCommand `command:"remove"`
}

func main() {
	var options Options
	p := flags.NewParser(&options, flags.Default)
	p.CommandHandler = func(command flags.Commander, args []string) error {
		store, err := store.GetDefaultStore()
		if err != nil {
			return err
		}
		c := command.(cmd.Command)
		c.SetConfig(&cmd.Config{
			Store: store,
		})
		return c.Execute(args)
	}
	if _, err := p.Parse(); err != nil {
		var flagsErr *flags.Error
		if errors.As(err, &flagsErr) && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}
}
