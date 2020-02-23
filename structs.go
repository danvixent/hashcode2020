package main

import (
	"fmt"
	"sort"
	"sync"
)

var signup sync.Mutex
var out sync.Mutex
var see sync.Mutex

var days = 0
var allLibs []library
var booksAndScores = make(map[int]bookScore)
var numOfLibsToShipFrom = 0
var alpha []int
var seen = make(map[int]bool)
var wait sync.WaitGroup

// describes each library
type library struct {
	ID           int
	SignUpTime   int
	ScansPerDay  int
	ScannedBooks *[]int
	BookIDs      *[]int
	IsSignedUp   bool
	Quality      float64
}

// i define this separately to prevent any form of mix up
type bookScore int64

func (l *library) calcQuality() {
	x := float64(len(*l.BookIDs) / l.ScansPerDay)
	tmp := (x / float64(l.SignUpTime)) * l.avgBookScore()
	l.Quality = tmp
}

func (l *library) avgBookScore() float64 {
	scores := 0.0
	for _, id := range *l.BookIDs {
		scores += float64(booksAndScores[id])
	}
	return scores / float64(len(*l.BookIDs))
}

func (l *library) sortBooksByScore() {
	fmt.Printf("Before Sorting %d's Books By Score -> %v\n", l.ID, l.BookIDs)
	sort.SliceStable(l, func(i, j int) bool {
		return booksAndScores[(*l.BookIDs)[i]] > booksAndScores[(*l.BookIDs)[j]]
	})
	fmt.Printf("After Sorting %d's Books By Score -> %v\n", l.ID, l.BookIDs)
}
