package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
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

			id := -1
			for val := range *(extract(strings.Split(lines[1], " "))) {
				id++
				allBooks[id] = bookScore(val)
			}

			nxtID := -1
			for i := 2; i < len(lines); i = i + 2 {
				tmp := strings.Split(lines[i], " ")
				nxtID++
				struc := &library{}
				struc.ID = nxtID
				struc.SignUpTime, _ = strconv.Atoi(tmp[1])
				struc.ScansPerDay, _ = strconv.Atoi(tmp[2])
				struc.BookIDs = *(extract(strings.Split(lines[i+1], " ")))
				struc.calcQuality()
				allLibs = append(allLibs, *struc)
			}
			sort.SliceStable(allLibs, func(i, j int) bool {
				return allLibs[i].Quality > allLibs[j].Quality
			})
			simulate()
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

func simulate() {
	count := 0
	for _, lib := range allLibs {
		go procLibs(&lib)
		count++
	}
	wait.Add(count)
}

func procLibs(lib *library) {
	<-signup
	defer wait.Done()
	time.Sleep(1)
	lib.IsSignedUp = true
	days = days - lib.SignUpTime
	signup <- true
	scanBooks(lib)
	runtime.Goexit()
}

func scanBooks(lib *library) {
	bksToScan := int(len(lib.BookIDs) / lib.ScansPerDay)
	for {
		if bksToScan > days {
			bksToScan -= lib.ScansPerDay
			continue
		}
		break
	}

	for ix := 0; ix < bksToScan; ix++ {
		if !seen[lib.BookIDs[ix]] {
			lib.ScannedBooks = append(lib.ScannedBooks, ix)
			seen[lib.BookIDs[ix]] = true
		}
	}
}
