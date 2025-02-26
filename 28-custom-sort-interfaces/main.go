package main

import (
	"fmt"
	"reflect"
	"sort"
	"time"
)

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return d
}

var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

type customSort struct {
	t       []*Track
	less    func(x, y *Track) bool
	history []int
}

func (c *customSort) Len() int           { return len(c.t) }
func (c *customSort) Less(i, j int) bool { return c.less(c.t[i], c.t[j]) }
func (c *customSort) Swap(i, j int)      { c.t[i], c.t[j] = c.t[j], c.t[i] }
func (c *customSort) SortByColumn(i int) {
	c.history = append(c.history, i)

	c.less = func(x, y *Track) bool {
		// Uses reflect to handle different variable types
		p1 := reflect.ValueOf(x)
		p2 := reflect.ValueOf(y)
		v1 := p1.Elem() //Dereference the pointer
		v2 := p2.Elem()
		for h := len(c.history) - 1; h >= 0; h-- {
			index := c.history[h]

			f1 := v1.Field(index)
			f2 := v2.Field(index)
			compare := compareValues(f1, f2)
			if compare == equal {
				continue
			}
			if compare == less {
				return true
			}
			return false
		}
		return false
	}
}

type compareResult int

const (
	less compareResult = iota
	equal
	greater
)

func compareValues(x reflect.Value, y reflect.Value) compareResult {
	switch x.Kind() {
	case reflect.Int, reflect.Int64:
		if x.Int() == y.Int() {
			return equal
		}
		if x.Int() < y.Int() {
			return less
		}
		return greater
	case reflect.String:
		if x.String() == y.String() {
			return equal
		}
		if x.String() < y.String() {
			return less
		}
		return greater
	default:
		panic("didn't handle sorting by type: " + x.Kind().String())
	}
}

func (c *customSort) PrintTracks() {
	for _, t := range c.t {
		fmt.Printf("Title: %v, Artist: %v, Album: %v, Year: %v, Length: %v\n", t.Title, t.Artist, t.Album, t.Year, t.Length)
	}
}

func main() {
	// Program simulates if a user first sorted by column 3 (year) and then sorted by column 0 (title)
	c := &customSort{t: tracks}
	c.SortByColumn(3)
	c.SortByColumn(0)
	sort.Sort(c)
	c.PrintTracks()

}
