package main

import (
	"fmt"
	"os"
	"strconv"
)

// extract returns the integer equivalents of numbers in the slice parameter...translated into a slice of ints
func extract(slice []string) *[]int {
	var tmp []int
	for ix := range slice {
		ref, err := strconv.Atoi(slice[ix])
		if err != nil {
			fmt.Println("Conversion failed")
			os.Exit(3) // we make this stringent because this error should never occur
		}
		tmp = append(tmp, ref)
	}
	return &tmp
}

func scanBooks(l *library, shippingDays int) {
	// l.sortBooksByScore()
	maxBooksToShip := shippingDays * l.ScansPerDay
	if maxBooksToShip < len(*l.BookIDs) {
		l.ScannedBooks = shipBooks((*l.BookIDs)[:maxBooksToShip])
		return
	}
	l.ScannedBooks = shipBooks(*l.BookIDs)
}

func shipBooks(IDs []int) *[]int {
	tmp := make([]int, 0)
	for _, id := range IDs {
		see.Lock()
		if !seen[id] {
			seen[id] = true
			tmp = append(tmp, id)
		}
		see.Unlock()
	}
	return &tmp
}
