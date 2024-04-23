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

type KeyChainItem struct {
	KeyChain string
	Name     string
	Class    string
	Account  string
}

func parseKeyChainItems(output string) ([]KeyChainItem, error) {
	account := os.Getenv("LOGNAME")
	lines := strings.Split(output, "\n")
	i := 0
	var items []KeyChainItem
	for i < len(lines) {
		if lines[i] == "" {
			break
		}
		if !strings.HasPrefix(lines[i], "keychain: \"") {
			return nil, fmt.Errorf("Invalid start of keychain item: %s", lines[i])
		}
		keychain := strings.TrimSuffix(strings.TrimPrefix(lines[i], "keychain: \""), "\"")
		item := KeyChainItem{
			KeyChain: keychain,
		}
		i += 1
		for i < len(lines) && !strings.HasPrefix(lines[i], "attributes:") && !strings.HasPrefix(lines[i], "keychain: ") {
			if strings.HasPrefix(lines[i], "class: ") {
				class := strings.Trim(strings.TrimPrefix(lines[i], "class: "), "\"")
				item.Class = class
			}
			i += 1
		}
		if i >= len(lines) {
			break
		}
		if strings.HasPrefix(lines[i], "keychain: ") {
			continue
		}
		i += 1
		for i < len(lines) && strings.HasPrefix(lines[i], "    ") {
			attr := strings.Trim(lines[i], " ")
			parts := strings.Split(attr, "=")
			if parts[0] == "0x00000007 <blob>" {
				item.Name = strings.Trim(parts[1], "\"")
			} else if parts[0] == "\"acct\"<blob>" {
				item.Account = strings.Trim(parts[1], "\"")
			}
			i += 1
		}
		if item.Class == "genp" && item.Account == account {
			items = append(items, item)
		}
	}
	return items, nil
}

func (a MacSecurityAdapter) List() ([]string, error) {
	cmd := exec.Command("security", "dump-keychain")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	items, err := parseKeyChainItems(string(output))
	if err != nil {
		return nil, err
	}
	names := make([]string, len(items))
	for i, item := range items {
		names[i] = item.Name
	}

	return names, nil
}
