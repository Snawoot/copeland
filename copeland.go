package copeland

import (
	"cmp"
	"errors"
	"fmt"
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

var ErrIncorrectLength = errors.New("incorrect number of items in ballot list")

type UnknownNameError string

func (e UnknownNameError) Error() string {
	return fmt.Sprintf("unknown name \"%s\" in ballot", string(e))
}

func (e UnknownNameError) Name() string {
	return string(e)
}

type MissingNameError string

func (e MissingNameError) Error() string {
	return fmt.Sprintf("missing name \"%s\" in ballot", string(e))
}

func (e MissingNameError) Name() string {
	return string(e)
}

type DuplicateNameError struct {
	name  string
	count int
}

func (e DuplicateNameError) Error() string {
	return fmt.Sprintf("name \"%s\" appears %d times in ballot", e.name, e.count)
}

func (e DuplicateNameError) Name() string {
	return e.name
}

func (e DuplicateNameError) Count() int {
	return e.count
}

func (c *Copeland) ballotToMatrix(ballot []string) (*Matrix, error) {
	size := len(c.names)
	if len(ballot) != size {
		return nil, ErrIncorrectLength
	}
	mapped := make([]int, size)
	counts := make([]int, size)
	for i, name := range ballot {
		mappedName, ok := c.nameToIndex(name)
		if !ok {
			return nil, UnknownNameError(name)
		}
		mapped[i] = mappedName
		counts[mappedName]++
	}
	for i, count := range counts {
		switch count {
		case 0:
			return nil, MissingNameError(c.names[i])
		case 1:
		default:
			return nil, DuplicateNameError{
				name:  c.names[i],
				count: count,
			}
		}
	}
	m := NewMatrix(size)
	for base, winner := range mapped {
		for _, loser := range mapped[base+1:] {
			m.Inc(winner, loser)
		}
	}
	return m, nil
}

func (c *Copeland) Update(ballot []string) error {
	matrix, err := c.ballotToMatrix(ballot)
	if err != nil {
		return fmt.Errorf("state update failed: %w", err)
	}
	c.state.Add(matrix)
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

func CmpScoreEntry(a, b ScoreEntry) int {
	if n := cmp.Compare(b.Score, a.Score); n != 0 {
		return n
	}
	return cmp.Compare(a.Name, b.Name)
}

func RankScore(scores []ScoreEntry) [][]ScoreEntry {
	copied := make([]ScoreEntry, len(scores))
	copy(copied, scores)
	slices.SortFunc(copied, CmpScoreEntry)
	return groupBy(copied, func(a, b ScoreEntry) bool { return a.Score == b.Score })
}

func groupBy[S ~[]E, E any](s S, eq func(E, E) bool) []S {
	var head []S
	for len(s) > 0 {
		key := s[0]
		var idx int
		for idx = 1; idx < len(s); idx++ {
			if !eq(s[idx], key) {
				break
			}
		}
		head = append(head, s[0:idx])
		s = s[idx:]
	}
	return head
}
