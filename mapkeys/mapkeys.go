package mapkeys

import (
	"fmt"
	"strings"

	"golang.org/x/exp/constraints"
)

// Just a map
type Store[T Number] map[string]T

type Number interface {
	constraints.Integer | constraints.Float
}

// Aggregation function interface: take any number of keys and aggregate them somehow (sum, mean, etc.)
// The result will always be a float64
type AggregationFunction[T Number] func(keysAndVals map[string]T) T

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

// AggregateDescendants aggregates values for all keys with a certain prefix.
// It returns a bool indicating whether or not the result is valid (or just a meaningless float64 zero value), and a float64
func (s *Store[T]) AggregateDescendants(prefix string, aggFunc AggregationFunction[T]) (bool, float64) {
	// get descendants with a prefix search
	descendants := s.SearchPrefix(prefix)
	if len(descendants) == 0 {
		return false, 0
	}

	// apply aggregation function
	return true, float64(aggFunc(descendants))
}

func Sum[T Number](keysAndVals map[string]T) T {
	var sum T
	for _, v := range keysAndVals {
		sum += v
	}
	return sum
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
