package adapter

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type MacSecurityAdapter struct {
}

func wrapError(key string, output []byte, err error) error {
	if err == nil {
		return nil
	}
	if strings.Contains(string(output), "SecKeychainSearchCopyNext: The specified item could not be found in the keychain.") {
		return &NoExistsKeyError{key: key}
	}
	return fmt.Errorf("Error from security: %v, %s", err, string(output))
}

func (a MacSecurityAdapter) Load(key string) (string, error) {
	cmd := exec.Command("security", "find-generic-password", "-a", os.Getenv("LOGNAME"), "-w", "-s", key)
	res, err := cmd.CombinedOutput()
	if err != nil {
		return "", wrapError(key, res, err)
	}

	return strings.TrimSuffix(string(res), "\n"), nil
}

func (a MacSecurityAdapter) Save(key string, value string) error {
	cmd := exec.Command("security", "add-generic-password", "-a", os.Getenv("LOGNAME"), "-s", key, "-w", value, "-U")
	output, err := cmd.CombinedOutput()
	return wrapError(key, output, err)
}

func (a MacSecurityAdapter) Delete(key string) error {
	cmd := exec.Command("security", "delete-generic-password", "-a", os.Getenv("LOGNAME"), "-s", key)
	output, err := cmd.CombinedOutput()
	return wrapError(key, output, err)
}
