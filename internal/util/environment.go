package util

import "fmt"

func MergeMapValuesWithExtras(env map[string]string, extras []string) []string {
	result := make([]string, 0, len(env)+len(extras))
	for k, v := range env {
		result = append(result, fmt.Sprintf("%s=%s", k, v))
	}
	result = append(result, extras...)
	return result
}
