package mapkeys

import (
	"fmt"
	"strings"
)

// Just a map
type Store[T any] map[string]T

func (s *Store[T]) Insert(key string, value T) {
	(*s)[key] = value
}

func (s *Store[T]) Search(key string) (bool, T) {
	var defaultResult T
	val, ok := (*s)[key]
	if !ok {
		return false, defaultResult
	}
	return true, val
}

func (s *Store[T]) SearchPrefix(prefix string) map[string]T {
	results := make(map[string]T)
	// we have to iterate over EVERYTHING
	for key, val := range *s {
		if strings.HasPrefix(key, prefix) {
			results[key] = val
		}
	}

	return results
}

func (s *Store[T]) String() string {
	var resultString string
	if s == nil {
		return ""
	}
	for key, val := range *s {
		resultString = resultString + fmt.Sprintf("\n%s - %v", key, val)
	}
	return resultString
}
