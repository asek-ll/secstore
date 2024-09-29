package adapter

import (
	"os/exec"
	"strings"
)

type SecretToolAdapter struct {
}

func (a SecretToolAdapter) Load(key string) (string, error) {
	cmd := exec.Command("secret-tool", "lookup", "secstore-key", key)
	res, err := cmd.CombinedOutput()
	if err != nil {
		return "", wrapError(key, res, err)
	}

	return strings.TrimSuffix(string(res), "\n"), nil
}

func (a SecretToolAdapter) Save(key string, value string) error {
	cmd := exec.Command("secret-tool", "store", "--label='"+key+"'", "system", "secstore", "secstore-key", key)
	cmd.Stdin = strings.NewReader(value)
	output, err := cmd.CombinedOutput()
	return wrapError(key, output, err)
}

func (a SecretToolAdapter) Delete(key string) error {
	cmd := exec.Command("secret-tool", "clear", "secstore-key", key)
	output, err := cmd.CombinedOutput()
	return wrapError(key, output, err)
}

func (a SecretToolAdapter) List() ([]string, error) {
	cmd := exec.Command("secret-tool", "search", "--all", "system", "secstore")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	items := parseIniRecords(string(output))
	var names []string
	for _, item := range items {
		k, e := item.Props["attribute.secstore-key"]
		if e {
			names = append(names, k)
		}
	}

	return names, nil
}

type SecretToolItem struct {
	Name  string
	Props map[string]string
}

func parseIniRecords(output string) []SecretToolItem {
	lines := strings.Split(output, "\n")
	var result []SecretToolItem
	i := 0
	for i < len(lines) {
		row := lines[i]
		if !strings.HasPrefix(row, "[") {
			break
		}
		item := SecretToolItem{
			Name:  strings.Trim(row, "[]"),
			Props: make(map[string]string),
		}
		i += 1
		for i < len(lines) && !strings.HasPrefix(lines[i], "[") {
			kv := strings.SplitN(lines[i], "=", 2)
			if len(kv) == 2 {
				item.Props[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			}
			i += 1
		}
		result = append(result, item)
	}

	return result
}
