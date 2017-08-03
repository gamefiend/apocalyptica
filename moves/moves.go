package moves

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

type Roll interface {
	Roll() int
	Display() string
}

func die(dieType int) int {
	now := time.Now()
	//now.Unix() seems to return seconds, so all die rolls in the same second
	// would be seeded the same.  UnixNano() returns a much more fine grained value
	rand.Seed(now.UnixNano())
	//'Roll' the die, +1 since Intn is 0 based
	return (rand.Intn(dieType)+ 1)
}

func (m Move) Roll(bonus int) int {
	return (die(6) + die(6) + bonus)
}

func (m Move) Display(r, bonus int) string {
	choice := ""
	switch {
	case r <= 6:
		choice = m[0].Miss
	case (r >= 7 && r <= 9):
		choice = m[0].Hit
	case (r >= 10):
		choice = m[0].Crit
	}

	return fmt.Sprintf("**%d**[%d *%d*] %s", r, (r - bonus), bonus, choice)
}

type Move []struct {
	Name string `json:"Name"`
	Full string `json:"Full"`
	Miss string `json:"Miss"`
	Hit  string `json:"Hit"`
	Crit string `json:"Crit"`
}

func LoadMoves(filename string) Move {
	var mv Move
	file, e := ioutil.ReadFile(filename)
	if e != nil {
		log.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	r := bytes.NewReader(file)
	if e := json.NewDecoder(r).Decode(&mv); e != nil {
		log.Printf("Problems decoding json: %v\n", e)
		os.Exit(1)
	}
	return mv
}

func FindMove(s string, mv Move) Move {
	var Found Move
	substr := strings.Split(s, " ")
	for _, v := range mv {
		r, err := regexp.Compile(v.Name)
		if err != nil {
			fmt.Printf("Problematic regexp, failing.\n")
			os.Exit(1)
		}
		if r.MatchString(substr[0]) == true {
			Found = append(Found, v)
		}
	}
	return Found
}
