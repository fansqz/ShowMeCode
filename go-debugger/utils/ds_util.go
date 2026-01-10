package utils

import (
	"github.com/emirpasic/gods/sets"
	"github.com/emirpasic/gods/sets/hashset"
)

func List2set[T any](list []T) sets.Set {
	set := hashset.New()
	for _, value := range list {
		set.Add(value)
	}
	return set
}
