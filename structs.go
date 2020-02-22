package main

var signup = make(chan bool, 1)
var signUpLock sync.Mutex

var nxtLib = -1
var allLibs []library
var allBooks []book
var alpha *[]int
var seen = make(map[int]bool)

// describes each library
type library struct {
	ID          int
	SignUpTime  int
	ScansPerDay int
	BookIDs     []int
	IsSignedUp  bool
	IsSelected  bool
}

func (l *library) calcQuality() {

}

// book describes a book
type book struct {
	ID        int
	IsScanned bool
}

// Scan a book
func (b *book) Scan() {
	b.IsScanned = true
}

func findBook(id int) *book {
	for _, book := range allBooks {
		if book.ID == id {
			return &book
		}
	}
	return nil
}
