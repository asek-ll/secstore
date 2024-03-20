package adapter

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
)

type GPGConfig struct {
	KeyId string `json:"keyId"`
	Path  string `json:"path"`
}

type GPGAdapter struct {
	GPGConfig
}

func NewGPGAdapter(config *GPGConfig) (*GPGAdapter, error) {
	if config == nil {
		return nil, fmt.Errorf("No gpg config present")
	}
	store := GPGAdapter{*config}
	err := store.ensureConfig()
	if err != nil {
		return nil, err
	}
	return &store, nil
}

func isGPGKeyExists(key string) (bool, error) {
	cmd := exec.Command("gpg", "--list-keys", key)
	if err := cmd.Start(); err != nil {
		return false, err
	}

	if err := cmd.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (a GPGAdapter) ensureConfig() error {
	if a.KeyId == "" {
		return errors.New("GPG key not specified")
	}
	e, err := isGPGKeyExists(a.KeyId)
	if err != nil {
		return err
	}
	if !e {
		return errors.New("GPG key not exists")
	}

	_, err = os.Stat(a.Path)
	if os.IsNotExist(err) {
		return os.MkdirAll(a.Path, os.ModePerm)
	}

	return err
}

func (a GPGAdapter) Load(key string) (string, error) {
	err := a.ensureConfig()
	if err != nil {
		return "", err
	}

	fullPath := path.Join(a.Path, key)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", &NoExistsKeyError{key: key}
	}

	cmd := exec.Command("gpg", "--decrypt", fullPath)

	result, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func (a GPGAdapter) Save(key string, value string) error {
	err := a.ensureConfig()
	if err != nil {
		return err
	}

	fullPath := path.Join(a.Path, key)

	if _, err = os.Stat(fullPath); err == nil {
		err = os.Remove(fullPath)
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("gpg", "--encrypt", "-o", fullPath, "--recipient", a.KeyId)
	in, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	go func() {
		defer in.Close()
		io.WriteString(in, value)
	}()
	return cmd.Run()
}

func (a GPGAdapter) Delete(key string) error {
	fullPath := path.Join(a.Path, key)

	if _, err := os.Stat(fullPath); err == nil {
		err = os.Remove(fullPath)
		if err != nil {
			return err
		}
	}
	return nil
}
