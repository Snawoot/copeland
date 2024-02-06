package copeland

import (
	"errors"
	"slices"
)

type Copeland struct {
	names []string
}

func New(names []string) (*Copeland, error) {
	copied := make([]string, len(names))
	copy(copied, names)
	slices.Sort(copied)
	copied = slices.Compact(copied)
	if len(copied) < 2 {
		return nil, errors.New("not enough alternatives")
	}
	return &Copeland{
		names: copied,
	}
}
