package main

import "sync"

var signup = make(chan bool, 1)
var signUpLock sync.Mutex

var days = (*alpha)[2]
var allLibs []library
var allBooks = make(map[int]bookScore)
var alpha *[]int
var seen = make(map[int]bool)
var wait sync.WaitGroup

// describes each library
type library struct {
	ID           int
	SignUpTime   int
	ScansPerDay  int
	ScannedBooks []int
	BookIDs      []int
	IsSignedUp   bool
	Quality      float64
}

// i define this separately to prevent any form of mix up
type bookScore int

func (l *library) calcQuality() {
	x := float64(len(l.BookIDs) / l.ScansPerDay)
	tmp := (x / float64(l.SignUpTime)) * l.avgBookScore()
	l.Quality = tmp
}

func (l *library) avgBookScore() float64 {
	scores := 0.0
	for _, id := range l.BookIDs {
		scores += float64(allBooks[id])
	}
	return scores / float64(len(l.BookIDs))
}
