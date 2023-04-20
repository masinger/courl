package util

import "strings"

func AllStringsPresent(values ...*string) bool {
	for _, s := range values {
		if s == nil || len(*s) == 0 {
			return false
		}
	}
	return true
}

func PresentOrDefault(value *string, defaultValue string) string {
	if value != nil && len(*value) > 0 {
		return *value
	}
	return defaultValue
}

func CoalesceStrings(values ...*string) *string {
	for _, value := range values {
		if value != nil {
			return value
		}
	}
	return nil
}

func NotEmpty(values ...*string) []*string {
	var result []*string
	for _, value := range values {
		if value != nil && len(strings.TrimSpace(*value)) > 0 {
			result = append(result, value)
		}
	}
	return result
}

func NotNill(values ...*string) []*string {
	var result []*string
	for _, value := range values {
		if value != nil {
			result = append(result, value)
		}
	}
	return result
}

func CopyIfPresent(from *string, target *string) {
	if AllStringsPresent(from) {
		*target = *from
	}
}
