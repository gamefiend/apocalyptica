package moves

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

//grabs the nanosecond as a seed
func seed() int64 {
	now := time.Now()
	return now.Unix()
}

func (m Move) Roll(bonus int) int {
	rand.Seed(seed())
	return ((rand.Intn(10)) + 2) + bonus
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

func LoadMoves() Move {
	var mv Move
	file, e := ioutil.ReadFile("basic.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	r := bytes.NewReader(file)
	if e := json.NewDecoder(r).Decode(&mv); e != nil {
		fmt.Printf("Problems decoding json: %v\n", e)
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
