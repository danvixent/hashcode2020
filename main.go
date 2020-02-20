package main

import (
	"runtime"
	"sort"
	"sync"
	"time"
)

var signup = make(chan bool, 1)
var signUpLock sync.Mutex

var nxtLib = -1
var allLibs []library
var allBooks []Book
var alpha []int
var seen = make(map[int]bool)

type library struct {
	SignUpTime  int
	ScansPerDay int
	Books       []Book
	IsSignedUp  bool
}

type Book struct {
	ID          int
	IsScannedBy []*library
}

func main() {
	sort.SliceStable(allLibs, func(i, j int) bool {
		return len(allLibs[i].Books) > len(allLibs[j].Books) && (allLibs[i].ScansPerDay > allLibs[j].ScansPerDay)
	})
	for _, lib := range allLibs {
		signup <- true
		go procLibs()
	}
	scanBooks()
}

func procLibs() {
	<-signup
	defer runtime.Goexit()

	signUpLock.Lock()
	nxtLib++
	signUpLock.Unlock()

	time.Sleep(1)
	allLibs[nxtLib].IsSignedUp = true
	go scanBooks(allLibs[nxtLib])
	signup <- true
}

func scanBooks() {

}
