package copeland

import (
	"errors"
	"slices"
)

type Scoring struct {
	Win  float64
	Tie  float64
	Loss float64
}

var DefaultScoring = &Scoring{
	Win:  1,
	Tie:  .5,
	Loss: 0,
}

type ScoreEntry struct {
	Name  string
	Score float64
}

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

func (m *Matrix) Row(i int) []int64 {
	return m.data[i*m.size : (i+1)*m.size]
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

func (c *Copeland) Score(scoring *Scoring) []ScoreEntry {
	if scoring == nil {
		scoring = DefaultScoring
	}
	size := c.state.Size()
	res := make([]ScoreEntry, 0, size)
	for i := 0; i < size; i++ {
		score := float64(0)
		for j := 0; j < size; j++ {
			if i == j {
				continue
			}
			runner := c.state.Get(i, j)
			opponent := c.state.Get(j, i)
			switch {
			case runner > opponent:
				score += scoring.Win
			case runner < opponent:
				score += scoring.Loss
			default:
				score += scoring.Tie
			}
		}
		res = append(res, ScoreEntry{
			Name:  c.names[i],
			Score: score,
		})
	}
	return res
}
