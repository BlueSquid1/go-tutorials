package main

import (
	"bytes"
	"fmt"
)

type IntSet struct {
	words []uint64
}

func (s *IntSet) Has(x int) (bool, error) {
	wordIndex := x / 64
	if wordIndex >= len(s.words) {
		return false, fmt.Errorf("x is beyond the word slice")
	}

	bit := x % 64

	mask := uint64(1 << bit)
	word := s.words[wordIndex]
	return word&mask != 0, nil
}

func (s *IntSet) Add(x int) {
	wordIndex := x / 64
	bit := x % 64

	for wordIndex >= len(s.words) {
		s.words = append(s.words, 0)
	}

	mask := uint64(1 << bit)
	s.words[wordIndex] |= mask
}

func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			mask := uint64(1 << j)
			if word&mask != 0 {
				value := (i * 64) + j
				fmt.Fprintf(&buf, " %d", value)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

func (s *IntSet) Len() int {
	counter := 0
	for _, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < 64; j++ {
			mask := uint64(1 << j)
			if word&mask != 0 {
				counter++
			}
		}
	}
	return counter
}

func (s *IntSet) Copy() *IntSet {
	copySet := new(IntSet)
	copySet.words = make([]uint64, len(s.words))
	copy(copySet.words, s.words)
	return copySet
}

func main() {
	x := IntSet{}
	y := IntSet{}
	x.Add(1)
	x.Add(200)
	y.Add(500)
	y.Add(550)
	x.UnionWith(&y)

	z := x.Copy()
	z.Add(2)
	fmt.Printf("z=%v x=%v\n", z, &x)
}
