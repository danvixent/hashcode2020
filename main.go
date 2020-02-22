package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	// create output directory
	err := os.Mkdir("outputs", os.ModePerm)
	if err != nil {
		fmt.Printf("Could not create outputs directory: %v\n", err)
	}
	// Read the `input` directory so that we don't have to
	// modify the code whenever we want to test other inputs
	str := time.Now()
	filepath.Walk("inputs", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				// we dont want to stop the whole app if just one file does not open
				return err
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			var lines []string
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}

			alpha = extract(strings.Split(lines[0], " "))
			allLibs = make([]library, (*alpha)[1])
			nxtID := -1
			for i := 2; i < len(lines); i = i + 2 {
				tmp := strings.Split(lines[i], " ")
				nxtID++
				struc := library{
					ID:          nxtID,
					SignUpTime:  strconv.Atoi(tmp[1]),
					ScansPerDay: strconv.Atoi(tmp[2]),
					BookIDs:     *(extract(strings.Split(lines[i+1], " "))),
				}
				allLibs = append(allLibs, struc)
			}
		}
		return nil
	})
	stp := time.Since(str).Seconds()
	fmt.Println("Time:", stp)
}

// extract returns the integer equivalents of numbers in the slice parameter...translated into a slice of ints
func extract(slice []string) *[]int {
	var tmp []int
	for ix := range slice {
		ref, err := strconv.Atoi(slice[ix])
		// Because the error isn't supposed to occur at all, i'll handle it here
		if err != nil {
			fmt.Println("Conversion failed")
			os.Exit(3) // we make this stringent because this error should never occur
		}
		tmp = append(tmp, ref)
	}
	return &tmp
}

func procLibs(lib *library) {
	<-signup
	defer runtime.Goexit()

	signUpLock.Lock()
	nxtLib++
	signUpLock.Unlock()

	time.Sleep(1)
	lib.IsSignedUp = true
	go scanBooks(lib)
	signup <- true
}

func scanBooks(lib *library) {
	bksToScan := int(len(lib.BookIDs) / lib.ScansPerDay)
	var days = (*alpha)[2]

	for {
		if bksToScan > days {
			bksToScan -= lib.ScansPerDay
			continue
		}
		break
	}

	for ix := 0; ix < bksToScan; ix++ {
		if !seen[lib.BookIDs[ix]] {
			findBook(lib.BookIDs[ix]).Scan()
			seen[lib.BookIDs[ix]] = true
		}
	}
}
