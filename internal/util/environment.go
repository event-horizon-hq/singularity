package util

func MergeMapValuesWithExtras(env map[string]string, extras []string) []string {
	result := make([]string, 0, len(env)+len(extras))
	for _, v := range env {
		result = append(result, v)
	}
	
	result = append(result, extras...)
	return result
}
