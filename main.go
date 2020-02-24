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
		os.Exit(3)
	}
	// walk through the input directory
	str := time.Now()
	filepath.Walk("inputs", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			var lines []string
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}

			alpha = *extract(strings.Split(lines[0], " "))
			days = alpha[2]                         //package level variable
			allLibs = make([]*library, 0, alpha[1]) //eliminate making calls to append() reallocate

			id := -1
			fmt.Println("line[1]->", lines[1])
			tmp := extract(strings.Split(lines[1], " "))
			for _, val := range *tmp {
				id++
				booksAndScores[id] = bookScore(val)
			}

			nxtID := -1
			for i := 2; i < len(lines); i = i + 2 {
				tmp := strings.Split(lines[i], " ")
				nxtID++
				struc := &library{}
				struc.ID = nxtID
				struc.SignUpTime, _ = strconv.Atoi(tmp[1])
				struc.ScansPerDay, _ = strconv.Atoi(tmp[2])
				struc.BookIDs = extract(strings.Split(lines[i+1], " "))
				struc.calcQuality()
				allLibs = append(allLibs, struc)
			}
			sort.SliceStable(allLibs, func(i, j int) bool {
				return allLibs[i].Quality > allLibs[j].Quality
			})
			simulate()
			printToFile(path)
			clearDataStructures()
		}
		return nil
	})
	stp := time.Since(str).Seconds()
	fmt.Println("Time:", stp)
}

func simulate() {
	count := 0
	for _, lib := range allLibs {
		// time.Sleep(10 * time.Millisecond)
		go procLibs(lib)
		count++
	}
	wait.Add(count)
	wait.Wait()
}

func procLibs(lib *library) {
	defer wait.Done()
	tmp := 0
	signup.Lock()
	lib.signUp()
	tmp = days
	signup.Unlock()

	lib.scanBooks(tmp)
	fmt.Println("Goroutine for ", lib.ID, "finished")
	runtime.Goexit()
}

//Print needed output to file
func printToFile(path string) {
	output := ""
	noLib := 0
	for _, lib := range allLibs {
		if lib.ScannedBooks == nil {
			continue
		}
		output += strconv.Itoa(lib.ID) + " " + strconv.Itoa(len(*lib.ScannedBooks)) + "\n"
		noLib++
		for _, id := range *lib.ScannedBooks {
			output += strconv.Itoa(id) + " "
		}
		output += "\n"
	}
	output = strconv.Itoa(noLib) + "\n" + output

	outFile := strings.Trim(path, "inputs")
	outFile = outFile[1:2] + ".out"
	outFile = "outputs/" + outFile
	f, err := os.Create(outFile)
	defer f.Close()
	if err != nil {
		fmt.Println("Cannot create output file: ", err)
		return // return since there is no file to write to
	}

	_, err = f.Write([]byte(output))
	if err != nil {
		fmt.Println("Cannot write output to file: ", err)
	}
	f.Sync()
}
