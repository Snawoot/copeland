package copeland

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type Matrix struct {
	data []int64
	size int
}

func NewMatrix(size int) *Matrix {
	if size < 0 {
		panic("negative matrix size!")
	}
	return &Matrix{
		data: make([]int64, size*size),
		size: size,
	}
}

func (m *Matrix) Size() int {
	return m.size
}

func (m *Matrix) Add(o *Matrix) *Matrix {
	if m.Size() != o.Size() {
		panic("sizes are not equal!")
	}
	for i := range m.data {
		m.data[i] += o.data[i]
	}
	return m
}

func (m *Matrix) Get(i, j int) int64 {
	return m.data[i*m.size+j]
}

func (m *Matrix) Set(i, j int, value int64) {
	m.data[i*m.size+j] = value
}

func (m *Matrix) Inc(i, j int) {
	m.data[i*m.size+j]++
}

func (m *Matrix) Dump() {
	for i := 0; i < m.size; i++ {
		for j := 0; j < m.size; j++ {
			fmt.Printf("%d\t", m.Get(i, j))
		}
		fmt.Println()
	}
}

type Copeland struct {
	names []string
	state *Matrix
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
		state: NewMatrix(len(copied)),
	}, nil
}

func (c *Copeland) nameToIndex(name string) (int, bool) {
	return slices.BinarySearch(c.names, name)
}

func (c *Copeland) ballotToMatrix(ballot []string) *Matrix {
	mapped := make([]int, len(ballot))
	for i, name := range ballot {
		mappedName, ok := c.nameToIndex(name)
		if !ok {
			panic("unexpected condition: validated ballot has non-mapable name")
		}
		mapped[i] = mappedName
	}
	m := NewMatrix(len(ballot))
	for base, winner := range mapped {
		for _, loser := range mapped[base+1:] {
			m.Inc(winner, loser)
		}
	}
	return m
}

func (c *Copeland) Update(ballot []string) error {
	if len(ballot) != len(c.names) {
		return errors.New("incorrect ballot length")
	}
	copied := make([]string, len(ballot))
	copy(copied, ballot)
	slices.Sort(copied)
	if !slices.Equal(copied, c.names) {
		return errors.New("invalid ballot")
	}
	c.state.Add(c.ballotToMatrix(ballot))
	return nil
}

func (c *Copeland) Dump() {
	fmt.Println(strings.Join(c.names, "\t"))
	c.state.Dump()
}
